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

	"github.com/dominant-strategies/go-quai/cmd/options"
	"github.com/dominant-strategies/go-quai/consensus"
	"github.com/dominant-strategies/go-quai/log"
	"github.com/dominant-strategies/go-quai/p2p/protocol"
)

// P2PNode represents a libp2p node
type P2PNode struct {
	// Host interface
	host.Host

	// Backend for handling consensus data
	consensus consensus.ConsensusBackend

	// List of peers to introduce us to the network
	bootpeers []peer.AddrInfo

	// DHT instance
	dht *dual.DHT

	// runtime context
	ctx context.Context
}

// Returns a new libp2p node.
// The node is created with the given context and options passed as arguments.
func NewNode(ctx context.Context) (*P2PNode, error) {
	ipAddr := viper.GetString(options.IP_ADDR)
	port := viper.GetString(options.PORT)

	// Load bootpeers
	bootpeers, err := loadBootPeers()
	if err != nil {
		log.Errorf("error loading bootpeers: %s", err)
		return nil, err
	}

	// Define a connection manager
	connectionManager, err := connmgr.NewConnManager(
		viper.GetInt(options.MAX_PEERS),   // LowWater
		2*viper.GetInt(options.MAX_PEERS), // HighWater
	)
	if err != nil {
		log.Fatalf("error creating libp2p connection manager: %s", err)
		return nil, err
	}

	// Create the libp2p host
	var dht *dual.DHT
	host, err := libp2p.New(
		// use a private key for persistent identity
		libp2p.Identity(GetNodeKey()),

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
			if viper.GetBool(options.PORTMAP) {
				return libp2p.NATPortMap()
			} else {
				return nil
			}
		}(),

		// Enable NAT detection service
		libp2p.EnableNATService(),

		// If publicly reachable, provide a relay service for other peers
		libp2p.EnableRelayService(),
		libp2p.EnableRelay(),

		// If behind NAT, automatically advertise relay address through relay peers
		// TODO: today the bootnodes act as static relays. In the future we should dynamically select relays from publicly reachable peers.
		libp2p.EnableAutoRelayWithStaticRelays(bootpeers),

		// Attempt to open a direct connection with relayed peers, using relay
		// nodes to coordinate the holepunch.
		libp2p.EnableHolePunching(),

		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			dht, err = dual.New(ctx, h,
				dual.WanDHTOption(
					kaddht.Mode(kaddht.ModeServer),
					kaddht.BootstrapPeers(bootpeers...),
					kaddht.RoutingTableRefreshPeriod(time.Minute),
					kaddht.ProtocolPrefix("quai"),
				),
			)
			return dht, err
		}),
	)
	if err != nil {
		log.Fatalf("error creating libp2p host: %s", err)
		return nil, err
	}

	if err := dht.Bootstrap(ctx); err != nil {
		log.Fatalf("error bootstrapping libp2p host: %s", err)
	}

	idOpts := []identify.Option{
		identify.UserAgent("go-quai-libp2p"),
		identify.ProtocolVersion(protocol.ProtocolVersion),
	}
	// Create the identity service
	idServ, err := identify.NewIDService(host, idOpts...)
	if err != nil {
		log.Fatalf("error creating libp2p identity service: %s", err)
		return nil, err
	}
	// Register the identity service with the host
	idServ.Start()

	log.Debugf("host created")

	// log the p2p node's ID
	nodeID := host.ID()
	log.Infof("node created: %s", nodeID)

	node := &P2PNode{
		ctx:       ctx,
		Host:      host,
		bootpeers: bootpeers,
		dht:       dht,
	}

	return node, nil
}

// Get the full multi-address to reach our node
func (p *P2PNode) p2pAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", p.ID()))
}

// Dial bootpeers and bootstrap the DHT
func (p *P2PNode) bootstrap() error {
	log.Warn("We are bootstrapping")
	// Bootstrap the dht
	if err := p.dht.Bootstrap(p.ctx); err != nil {
		log.Warnf("error bootstrapping DHT: %s", err)
		return err
	}

	// Load bootpeers
	bootpeers, err := loadBootPeers()
	if err != nil {
		log.Errorf("error loading bootpeers: %s", err)
		return err
	}

	closestPeers, err := p.dht.WAN.GetClosestPeers(p.ctx, bootpeers[0].ID.String())
	log.Warnf("closest peers: %v", closestPeers)
	if err != nil {
		log.Warnf("ERROR GETTING PEERS: %s", err)
		return err
	}
	for _, peer := range closestPeers {
		routing.PublishQueryEvent(p.ctx, &routing.QueryEvent{
			ID:   peer,
			Type: routing.FinalPeer,
		})
	}
	return nil
}
