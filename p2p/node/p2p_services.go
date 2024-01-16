package node

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/p2p/pb"
	"github.com/dominant-strategies/go-quai/p2p/protocol"
	"github.com/pkg/errors"
)

// Opens a stream to the given peer and requests a block for the given hash and slice.
//
// If a block is not found, an error is returned
func (p *P2PNode) requestBlockFromPeer(hash common.Hash, slice types.SliceID, peerID peer.ID) (*types.Block, error) {
	// Open a stream to the peer using a specific protocol for block requests
	stream, err := p.NewStream(peerID, protocol.ProtocolVersion)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	// create a block request protobuf message
	blockReq, err := pb.EncodeQuaiRequest(pb.QuaiRequestMessage_REQUEST_BLOCK, &slice, &hash)
	if err != nil {
		return nil, err
	}

	// Send the block request to the peer
	err = common.WriteMessageToStream(stream, blockReq)
	if err != nil {
		return nil, err
	}

	// Read the response from the peer
	blockResponse, err := common.ReadMessageFromStream(stream)
	if err != nil {
		return nil, err
	}

	// Decode the response
	action, data, err := pb.DecodeQuaiResponse(blockResponse)
	if err != nil {
		return nil, err
	}

	if action != pb.QuaiResponseMessage_RESPONSE_BLOCK {
		return nil, errors.New("invalid response type")
	}

	block, ok := data.(*types.Block)
	if !ok {
		return nil, errors.New("invalid block type")
	}

	return block, nil

}
