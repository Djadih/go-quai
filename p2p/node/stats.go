package node

import (
	"runtime/debug"
	"time"

	"github.com/dominant-strategies/go-quai/log"
	"github.com/prometheus/client_golang/prometheus"
)

// Returns the number of peers in the routing table, as well as how many active
// connections we currently have.
func (p *P2PNode) connectionStats() int {
	peers := p.peerManager.GetHost().Network().Peers()
	numConnected := len(peers)

	return numConnected
}

var (
	peerMetrics *prometheus.GaugeVec
)

func (p *P2PNode) statsLoop() {
	defer func() {
		if r := recover(); r != nil {
			p.quitCh <- struct{}{}
			log.Global.WithFields(log.Fields{
				"error":      r,
				"stacktrace": string(debug.Stack()),
			}).Error("Go-Quai Panicked")
		}
	}()
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			peersConnected := p.connectionStats()
			peerMetrics.WithLabelValues("numPeers").Set(float64(peersConnected))
			log.Global.Debugf("Number of peers connected: %d", peersConnected)
		case <-p.ctx.Done():
			log.Global.Warnf("Context cancelled. Stopping stats loop...")
			return
		}
	}
}
