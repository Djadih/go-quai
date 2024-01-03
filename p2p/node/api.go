package node

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/dominant-strategies/go-quai/cmd/utils"
	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/consensus/types"
	"github.com/dominant-strategies/go-quai/log"
	"github.com/dominant-strategies/go-quai/p2p"
	quaiprotocol "github.com/dominant-strategies/go-quai/p2p/protocol"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Starts the node and all of its services
func (p *P2PNode) Start() error {
	log.Infof("starting P2P node...")

	// Start any async processes belonging to this node
	log.Debugf("starting node processes...")
	go p.eventLoop()
	go p.statsLoop()

	// Is this node expected to have bootstrap peers to dial?
	if !viper.GetBool(utils.BootNodeFlag.Name) && !viper.GetBool(utils.SoloFlag.Name) && len(p.bootpeers) == 0 {
		err := errors.New("no bootpeers provided. Unable to join network")
		log.Errorf("%s", err)
		return err
	}

	// Register the Quai protocol handler
	p.SetStreamHandler(quaiprotocol.ProtocolVersion, func(s network.Stream) {
		quaiprotocol.QuaiProtocolHandler(s, p)
	})

	// If the node is a bootnode, start the bootnode service
	if viper.GetBool(utils.BootNodeFlag.Name) {
		log.Infof("starting node as a bootnode...")
		return nil
	}

	// Open data streams with connected Quai peers
	go quaiprotocol.OpenPeerStreams(p)

	return nil
}

type stopFunc func() error

// Function to gracefully shtudown all running services
func (p *P2PNode) Stop() error {
	// define a list of functions to stop the services the node is running
	stopFuncs := []stopFunc{
		p.Host.Close,
		p.dht.Close,
	}
	// create a channel to collect errors
	errs := make(chan error, len(stopFuncs))
	// run each stop function in a goroutine
	for _, fn := range stopFuncs {
		go func(fn stopFunc) {
			errs <- fn()
		}(fn)
	}

	var allErrors []error
	for i := 0; i < len(stopFuncs); i++ {
		select {
		case err := <-errs:
			if err != nil {
				log.Errorf("error during shutdown: %s", err)
				allErrors = append(allErrors, err)
			}
		case <-time.After(5 * time.Second):
			err := errors.New("timeout during shutdown")
			log.Warnf("error: %s", err)
			allErrors = append(allErrors, err)
		}
	}
	close(errs)
	if len(allErrors) > 0 {
		return errors.Errorf("errors during shutdown: %v", allErrors)
	} else {
		return nil
	}
}

func (p *P2PNode) SetConsensusBackend(be common.ConsensusAPI) {
	p.consensus = be
}

func (p *P2PNode) BroadcastBlock(slice types.SliceID, block types.Block) error {
	return p.pubsub.BroadcastBlock(p.topics[slice][p2p.C_blockTopicName], block)
}

func (p *P2PNode) BroadcastTransaction(tx types.Transaction) error {
	panic("todo")
}

func (p *P2PNode) BroadcastHeader(slice types.SliceID, header types.Header) error {
	panic("todo")
}

func (p *P2PNode) requestDHT(slice types.SliceID, hash types.Hash, reqFunc func(types.Hash, peer.ID)(*types.Block, error)) (chan interface{}, error) {
	resultChan := make(chan interface{}, 1)
	const (
		maxDHTQueryRetries    = 3  // Maximum number of retries for DHT queries
		peersPerDHTQuery      = 10 // Number of peers to query per DHT attempt
		dhtQueryRetryInterval = 5  // Time to wait between DHT query retries
	)
	// create a Cid from the slice ID
	shardCid := shardToCid(slice)
	for retries := 0; retries < maxDHTQueryRetries; retries++ {
		log.Debugf("Querying DHT for slice Cid %s (retry %d)", shardCid, retries)
		// query the DHT for peers in the slice
		peerChan := p.dht.FindProvidersAsync(p.ctx, shardCid, peersPerDHTQuery)
		for peerInfo := range peerChan {
			block, err := p.requestBlockFromPeer(hash, peerInfo.ID)
			if err == nil {
				log.Debugf("Received block %s from peer %s", block.Hash, peerInfo.ID)
				p.blockCache.Add(hash, block)
				resultChan <- block
				return resultChan, nil
			}
		}
		// if the block is not found, wait for a bit and try again
		log.Debugf("Block %s not found in slice %s. Retrying...", hash, slice)
		time.Sleep(dhtQueryRetryInterval * time.Second)
	}
	return nil, errors.New("item not found")
}

// Request a block from the network for the specified slice
func (p *P2PNode) RequestBlock(hash types.Hash, slice types.SliceID) chan *types.Block {
	resultChan, err := p.requestDHT(slice, hash, p.requestBlockFromPeer)
	if err != nil {
		log.Errorf("error requesting block from DHT: %s", err)
		return nil
	}
	return resultChan.(chan *types.Block)

	value := <-resultChan
	block := value.(*types.Block)

	return block
}

// Request a header from the network for the specified slice
func (p *P2PNode) RequestHeader(hash types.Hash, slice types.SliceID) chan *types.Header {
	panic("todo")
}

func (p *P2PNode) RequestTransaction(hash types.Hash, loc types.SliceID) chan *types.Transaction {
	panic("todo")
}

func (p *P2PNode) ReportBadPeer(peer p2p.PeerID) {
	panic("todo")
}

// Returns the list of bootpeers
func (p *P2PNode) GetBootPeers() []peer.AddrInfo {
	return p.bootpeers
}

// Opens a new stream to the given peer using the given protocol ID
func (p *P2PNode) NewStream(peerID peer.ID, protocolID protocol.ID) (network.Stream, error) {
	return p.Host.NewStream(p.ctx, peerID, protocolID)
}

// Connects to the given peer
func (p *P2PNode) Connect(pi peer.AddrInfo) error {
	return p.Host.Connect(p.ctx, pi)
}

// Start gossipsub protocol
func (p *P2PNode) StartGossipSub(ctx context.Context) error {
	for _, slice := range p.consensus.GetRunningSlices() {
		blockTopic, err := p.pubsub.Join(slice.SliceID.String() + "/" + p2p.C_blockTopicName)
		if err != nil {
			return err
		}
		sub, err := blockTopic.Subscribe()
		if err != nil {
			return err
		}

		go p.handleBlocksSubscription(sub)

		sliceTopics, exists := p.topics[slice.SliceID]
		if !exists {
			sliceTopics = make(map[string]*pubsub.Topic)
		}
		sliceTopics[p2p.C_blockTopicName] = blockTopic
		p.topics[slice.SliceID] = sliceTopics
	}
	return nil
}

// Checks if the cache has a block with the given hash. If the block is not found, returns nil.
func (p *P2PNode) GetBlock(hash types.Hash) *types.Block {
	block, ok := p.blockCache.Get(hash)
	if !ok {
		return nil
	}
	return block
}

func (p *P2PNode) GetHeader(hash types.Hash) *types.Header {
	tx, ok := p.headerCache.Get(hash)
	if !ok {
		return nil
	}
	return tx
}

func (p *P2PNode) LookupTransaction(hash types.Hash) *types.Transaction {
	tx, ok := p.txCache.Get(hash)
	if !ok {
		return nil
	}
	return tx
}
