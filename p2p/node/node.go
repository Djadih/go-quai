package node

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"

	"github.com/dominant-strategies/go-quai/cmd/options"
	"github.com/dominant-strategies/go-quai/consensus"
	"github.com/dominant-strategies/go-quai/log"
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
	dht *kaddht.IpfsDHT

	// runtime context
	ctx context.Context
}

// Returns a new libp2p node.
// The node is created with the given context and options passed as arguments.
func NewNode(ctx context.Context) (*P2PNode, error) {
	ipAddr := viper.GetString(options.IP_ADDR)
	port := viper.GetString(options.PORT)

	// Load bootpeers
	var bootpeers []peer.AddrInfo
	for _, p := range viper.GetStringSlice(options.BOOTPEERS) {
		addr, err := multiaddr.NewMultiaddr(p)
		if err != nil {
			log.Fatalf("error loading multiaddress for %s in bootpeers: %s", p, err)
		}
		info, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Fatalf("error getting address info for %s in bootpeers: %s", addr, err)
		}
		bootpeers = append(bootpeers, *info)
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
	var dht *kaddht.IpfsDHT
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

		// If behind NAT, automatically advertise relay address through relay peers
		// TODO: today the bootnodes act as static relays. In the future we should dynamically select relays from publicly reachable peers.
		libp2p.EnableAutoRelayWithStaticRelays(bootpeers),

		// Attempt to open a direct connection with relayed peers, using relay
		// nodes to coordinate the holepunch.
		libp2p.EnableHolePunching(),

		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			dht, err = kaddht.New(ctx, h,
				kaddht.Mode(kaddht.ModeServer),
				kaddht.BootstrapPeers(bootpeers...),
				kaddht.BootstrapPeersFunc(func() []peer.AddrInfo { return bootpeers }),
			)
			return dht, err
		}),
	)
	if err != nil {
		log.Fatalf("error creating libp2p host: %s", err)
		return nil, err
	}

	log.Debugf("host created")

	// log the p2p node's ID
	nodeID := host.ID().Pretty()
	log.Infof("node created: %s", nodeID)

	return &P2PNode{
		ctx:       ctx,
		Host:      host,
		bootpeers: bootpeers,
		dht:       dht,
	}, nil
}

// Get the full multi-address to reach our node
func (p *P2PNode) p2pAddress() (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", p.ID()))
}

// Dial bootpeers and bootstrap the DHT
func (p *P2PNode) bootstrap() error {
	// Connect to boot peers
	for _, addr := range p.bootpeers {
		log.Debugf("dialing boot peer: %s", addr)
		if err := p.Host.Connect(p.ctx, addr); err != nil {
			log.Warnf("error dialing boot peer %s: %s", addr, err)
		}

		ch := ping.Ping(context.Background(), p.Host, addr.ID)
		for i := 0; i < 3; i++ {
			res := <-ch
			if res.Error != nil {
				log.Warnf("error pinging boot peer %s: %s", addr, res.Error)
			} else {
				log.Warnf("pinged boot peer %s in %s", addr, res.RTT)
			}
		}
	}

	// Bootstrap the dht
	if err := p.dht.Bootstrap(p.ctx); err != nil {
		log.Warnf("error bootstrapping DHT: %s", err)
		return err
	}
	return nil
}
