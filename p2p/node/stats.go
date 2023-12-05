package node

import (
	"time"

	"github.com/dominant-strategies/go-quai/log"
)

// Returns the number of peers in the routing table, as well as how many active
// connections we currently have.
func (p *P2PNode) connectionStats() (int, int, int) {
	WANroutingTableSize := p.dht.WAN.RoutingTable().Size()
	LANroutingTableSize := p.dht.LAN.RoutingTable().Size()
	peers := p.Host.Network().Peers()
	numConnected := len(peers)

	log.Info("Connected peers: %s", peers)

	return WANroutingTableSize, LANroutingTableSize, numConnected
}

func (p *P2PNode) statsLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			WANsize, LANsize, numConnected := p.connectionStats()
			log.Infof("Routing Table Size: WAN-%d, LAN-%d, Number of Connected Peers: %d", WANsize, LANsize, numConnected)
		case <-p.ctx.Done():
			log.Warnf("Context cancelled. Stopping stats loop...")
			return
		}

	}
}
