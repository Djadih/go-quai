package protocol

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/dominant-strategies/go-quai/consensus/types"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// interface required to join the quai protocol network
type QuaiP2PNode interface {
	GetBootPeers() []peer.AddrInfo
	Connect(ctx context.Context, pi peer.AddrInfo) error
	NewStream(ctx context.Context, peerID peer.ID, pids ...protocol.ID) (network.Stream, error)
	Network() network.Network
	// Search for a block in the node's cache, or query the consensus backend if it's not found in cache.
	// Returns nil if the block is not found.
	GetBlock(hash types.Hash, slice types.SliceID) *types.Block
}
