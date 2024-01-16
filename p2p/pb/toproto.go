package pb

import (
	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/pkg/errors"
)

func convertDataToProtoResponse(action QuaiResponseMessage_ActionType, data interface{}) (*Response, error) {
	switch action {
	case QuaiResponseMessage_RESPONSE_BLOCK:
		if block, ok := data.(*types.Block); ok {
			protoBlock := convertBlockToProto(block)
			return &Response{
				Response: &Response_Block{
					Block: protoBlock,
				},
			}, nil
		}
	case QuaiResponseMessage_RESPONSE_HEADER:
		if header, ok := data.(*types.Header); ok {
			protoHeader := convertHeaderToProto(header)
			return &Response{
				Response: &Response_Header{
					Header: protoHeader,
				},
			}, nil
		}
	case QuaiResponseMessage_RESPONSE_TRANSACTION:
		if transaction, ok := data.(*types.Transaction); ok {
			protoTransaction := convertTransactionToProto(transaction)
			return &Response{
				Response: &Response_Transaction{
					Transaction: protoTransaction,
				},
			}, nil
		}
	}
	return nil, errors.Errorf("invalid data type or action")
}

func convertBlockToProto(block *types.Block) *Block {
	protoBlock := new(Block)
	//! TODO: implement
	return protoBlock
}

func convertHeaderToProto(header *types.Header) *Header {
	protoHeader := new(Header)
	protoHeader.GasLimit = header.GasLimit()
	protoHeader.GasUsed = header.GasUsed()
	// TODO: implement
	return protoHeader
}

func convertTransactionToProto(transaction *types.Transaction) *Transaction {
	panic("TODO: implement")
	// return nil
}

func convertHashToProto(hash *common.Hash) (*Hash, error) {
	hashBytes := hash.Bytes()
	protoHash := &Hash{
		Hash: hashBytes[:],
	}
	// TODO: implement
	return protoHash, nil
}

func convertSliceIDToProto(sliceID *types.SliceID) (*SliceID, error) {
	protoSliceID := &SliceID{
		Region: sliceID.Region,
		Zone:   sliceID.Zone,
	}
	// TODO: implement
	return protoSliceID, nil
}
