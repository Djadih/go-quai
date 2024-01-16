package protocol

import (
	"fmt"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/log"
	"github.com/dominant-strategies/go-quai/p2p/pb"
	"github.com/libp2p/go-libp2p/core/network"
)

func QuaiProtocolHandler(stream network.Stream, node QuaiP2PNode) {
	defer stream.Close()

	log.Debugf("Received a new stream from %s", stream.Conn().RemotePeer())

	// if there is a protocol mismatch, close the stream
	if stream.Protocol() != ProtocolVersion {
		log.Warnf("Invalid protocol: %s", stream.Protocol())
		// TODO: add logic to drop the peer
		return
	}

	// Enter the read loop for the stream and handle messages
	for {
		data, err := common.ReadMessageFromStream(stream)
		if err != nil {
			log.Errorf("error reading message from stream: %s", err)
			return
		}

		action, slice, hash, err := pb.DecodeQuaiRequest(data)
		if err != nil {
			log.Errorf("error decoding quai request: %s", err)
			return
		}

		switch action {
		case pb.QuaiRequestMessage_REQUEST_BLOCK:
			log.Tracef("Received block request for slice %+v and hash %s", slice, hash)
			handleBlockRequest(slice, hash, stream, node)
		case pb.QuaiRequestMessage_REQUEST_HEADER:
			log.Tracef("Received header request for slice %+v and hash %s", slice, hash)
			handleHeaderRequest(slice, hash, stream, node)
		case pb.QuaiRequestMessage_REQUEST_TRANSACTION:
			handleTransactionRequest(slice, hash, stream, node)
		default:
			log.Errorf("invalid action type: %s", action)
			return
		}
	}
}

// Seeks the block in the cache or database and sends it to the peer in a Quai Response Message
func handleBlockRequest(slice *types.SliceID, hash *common.Hash, stream network.Stream, node QuaiP2PNode) {
	// check if we have the block in our cache or database
	block := node.GetBlock(*hash, *slice)
	if block == nil {
		// TODO: handle block not found
		log.Warnf("block not found")
		return
	}
	// create a Quai Message Response with the block
	action := pb.QuaiResponseMessage_RESPONSE_BLOCK
	data, err := pb.EncodeQuaiResponse(action, block)
	if err != nil {
		log.Errorf("error encoding quai response: %s", err)
		return
	}
	err = common.WriteMessageToStream(stream, data)
	if err != nil {
		log.Errorf("error writing message to stream: %s", err)
		// TODO: handle error
		return
	}
	log.Debugf("Sent block %s to peer %s", block.Hash, stream.Conn().RemotePeer())
}

func handleHeaderRequest(slice *types.SliceID, hash *common.Hash, stream network.Stream, node QuaiP2PNode) {
	header := node.GetHeader(*hash, *slice)
	if header == nil {
		// TODO: handle header not found
		log.Warnf("header not found")
		return
	}
	log.Tracef("header found: %+v", header)
	// create a Quai Message Response with the block
	action := pb.QuaiResponseMessage_RESPONSE_HEADER
	data, err := pb.EncodeQuaiResponse(action, header)
	if err != nil {
		log.Errorf("error encoding quai response: %s", err)
		return
	}
	fmt.Printf("EncodeQuaiResponse - data length: %d\n", len(data))
	err = common.WriteMessageToStream(stream, data)
	if err != nil {
		log.Errorf("error writing message to stream: %s", err)
		// TODO: handle error
		return
	}
}

func handleTransactionRequest(slice *types.SliceID, hash *common.Hash, stream network.Stream, node QuaiP2PNode) {
	// TODO: implement
}
