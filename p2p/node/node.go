package node

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	dual "github.com/libp2p/go-libp2p-kad-dht/dual"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"

	"github.com/dominant-strategies/go-quai/cmd/utils"
	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/log"
	"github.com/dominant-strategies/go-quai/p2p/protocol"
	"github.com/dominant-strategies/go-quai/p2p/pubsubManager"
	"github.com/dominant-strategies/go-quai/quai"
	lru "github.com/hashicorp/golang-lru/v2"
)

// P2PNode represents a libp2p node
type P2PNode struct {
	// Host interface
	host.Host

	// Backend for handling consensus data
	consensus quai.ConsensusAPI

	// List of peers to introduce us to the network
	bootpeers []peer.AddrInfo

	// Set of nodes to block connections to/from
	peerBlackList map[peer.ID]struct{}

	// TODO: Consolidate into network interface, and consensus interface
	// DHT instance
	dht *dual.DHT

	// Gossipsub instance
	pubsub *pubsubManager.PubsubManager

	// Caches for each type of data we may receive
	cache map[string]*lru.Cache[common.Hash, interface{}]

	// runtime context
	ctx context.Context
}

// Returns a new libp2p node.
// The node is created with the given context and options passed as arguments.
func NewNode(ctx context.Context) (*P2PNode, error) {
	ipAddr := viper.GetString(utils.IPAddrFlag.Name)
	port := viper.GetString(utils.P2PPortFlag.Name)

	// Load bootpeers
	bootpeers, err := loadBootPeers()
	if err != nil {
		log.Errorf("error loading bootpeers: %s", err)
		return nil, err
	}

	// Define a connection manager
	connectionManager, err := connmgr.NewConnManager(
		viper.GetInt(utils.MaxPeersFlag.Name),   // LowWater
		2*viper.GetInt(utils.MaxPeersFlag.Name), // HighWater
	)
	if err != nil {
		log.Fatalf("error creating libp2p connection manager: %s", err)
		return nil, err
	}

	peerBlackList := make(map[peer.ID]struct{})

	// Create the libp2p host
	var dht *dual.DHT
	host, err := libp2p.New(
		// use a private key for persistent identity
		libp2p.Identity(getNodeKey()),

		// pass the ip address and port to listen on
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/%s/tcp/%s", ipAddr, port),
		),

		// support all transports
		libp2p.DefaultTransports,

		// support Noise connections
		libp2p.Security(noise.ID, noise.New),

		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connectionManager),

		// Optionally attempt to configure network port mapping with UPnP
		func() libp2p.Option {
			if viper.GetBool(utils.PortMapFlag.Name) {
				return libp2p.NATPortMap()
			} else {
				return nil
			}
		}(),

		// Enable NAT detection service
		libp2p.EnableNATService(),

		// If publicly reachable, provide a relay service for other peers
		libp2p.EnableRelayService(),

		// If behind NAT, automatically advertise relay address through relay peers
		// TODO: today the bootnodes act as static relays. In the future we should dynamically select relays from publicly reachable peers.
		libp2p.EnableAutoRelayWithStaticRelays(bootpeers),

		// Attempt to open a direct connection with relayed peers, using relay
		// nodes to coordinate the holepunch.
		libp2p.EnableHolePunching(),

		// Create a connection gater that will reject connections from peers in the blocklist
		libp2p.ConnectionGater(pubsubManager.NewConnGater(&peerBlackList)),

		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			dht, err = dual.New(ctx, h,
				dual.WanDHTOption(
					kaddht.Mode(kaddht.ModeServer),
					kaddht.BootstrapPeersFunc(func() []peer.AddrInfo {
						log.Debugf("Bootstrapping to the following peers: %v", bootpeers)
						return bootpeers
					}),
					kaddht.ProtocolPrefix("/quai"),
					kaddht.RoutingTableRefreshPeriod(1*time.Minute),
				),
			)
			return dht, err
		}),
	)
	if err != nil {
		log.Fatalf("error creating libp2p host: %s", err)
		return nil, err
	}

	idOpts := []identify.Option{
		// TODO: Add version number + commit hash
		identify.UserAgent("go-quai"),
		identify.ProtocolVersion(string(protocol.ProtocolVersion)),
	}

	// Create the identity service
	idServ, err := identify.NewIDService(host, idOpts...)
	if err != nil {
		log.Fatalf("error creating libp2p identity service: %s", err)
		return nil, err
	}
	// Register the identity service with the host
	idServ.Start()

	// log the p2p node's ID
	nodeID := host.ID()
	log.Infof("node created: %s", nodeID)

	// Create a gossipsub instance with helper functions
	ps, err := pubsubManager.NewGossipSubManager(ctx, host)

	if err != nil {
		return nil, err
	}

	// Create a new LRU cache for each data-type we support caching
	cache := map[string]*lru.Cache[common.Hash, interface{}]{
		"blocks": func() *lru.Cache[common.Hash, interface{}] {
			cache, err := lru.New[common.Hash, interface{}](10)
			if err != nil {
				log.Fatal("error initializing cache;", err)
			}
			return cache
		}(),
	}

	return &P2PNode{
		ctx:       ctx,
		Host:      host,
		bootpeers: bootpeers,
		dht:       dht,
		pubsub:    ps,
		cache:     cache,
		peerBlackList: peerBlackList,
	}, nil
}

// Get the full multi-address to reach our node
func (p *P2PNode) p2pAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", p.ID()))
}

// Helper to access the corresponding data cache
func (p *P2PNode) pickCache(datatype interface{}) *lru.Cache[common.Hash, interface{}] {
	switch datatype.(type) {
	case *types.Block:
		return p.cache["blocks"]
	default:
		log.Fatalf("unsupported type")
		return nil
	}
}

// Add a datagram into the corresponding cache
func (p *P2PNode) cacheAdd(hash common.Hash, data interface{}) {
	cache := p.pickCache(data)
	cache.Add(hash, data)
}

// Get a datagram from the corresponding cache
func (p *P2PNode) cacheGet(hash common.Hash, datatype interface{}) (interface{}, bool) {
	cache := p.pickCache(datatype)
	return cache.Get(hash)
}
