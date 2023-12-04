package node

import (
	"time"

	"github.com/dominant-strategies/go-quai/log"
)

// Returns the number of peers in the routing table, as well as how many active
// connections we currently have.
func (p *P2PNode) connectionStats() (int, int) {
	routingTableSize := p.dht.RoutingTable().Size()
	numConnected := len(p.Host.Network().Peers())
	return routingTableSize, numConnected
}

func (p *P2PNode) statsLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.bootstrap("12D3KooWCCueXNT8qnrUVq78KEVg9xhKiFkrmH2nk4K49741kqUx")
			routingTableSize, numConnected := p.connectionStats()
			log.Infof("Routing Table Size: %d, Number of Connected Peers: %d", routingTableSize, numConnected)
		case <-p.ctx.Done():
			log.Warnf("Context cancelled. Stopping stats loop...")
			return
		}

	}
}
