package node

import (
	"time"

	"github.com/dominant-strategies/go-quai/log"
)

// Returns the number of peers in the routing table, as well as how many active
// connections we currently have.
func (p *P2PNode) connectionStats() (int) {
	peers := p.Host.Network().Peers()
	numConnected := len(peers)

	return numConnected
}

func (p *P2PNode) statsLoop() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			peersConnected := p.connectionStats()

			log.Global.Debugf("Number of peers connected: %d", peersConnected)
		case <-p.ctx.Done():
			log.Global.Warnf("Context cancelled. Stopping stats loop...")
			return
		}
	}
}
