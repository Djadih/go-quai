package node

import (
	"time"

	"github.com/dominant-strategies/go-quai/log"
)

// Returns the number of peers in the routing table, as well as how many active
// connections we currently have.
func (p *P2PNode) connectionStats() {
	WANroutingTable := p.dht.WAN.RoutingTable().ListPeers()
	LANroutingTable := p.dht.LAN.RoutingTable().ListPeers()
	peers := p.Host.Network().Peers()
	numConnected := len(peers)

	log.Warnf("Routing Table Size: WAN-%d, LAN-%d, Number of Connected Peers: %d", len(WANroutingTable), len(LANroutingTable), numConnected)
	log.Infof("Entries in WAN table: %v", WANroutingTable)
	log.Infof("Entries in LAN table: %v", LANroutingTable)
}

func (p *P2PNode) statsLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.connectionStats()
			peers, err := p.dht.WAN.GetClosestPeers(p.ctx, p.Host.ID().String())
			if err != nil {
				log.Warnf("error getting closest peers: %s", err)
			}

			for _, peer := range peers {
				log.Infof("Peer: %s", peer)
			}
		case <-p.ctx.Done():
			log.Warnf("Context cancelled. Stopping stats loop...")
			return
		}

	}
}
