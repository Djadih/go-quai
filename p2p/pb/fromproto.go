package pb

import (
	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
)

func convertProtoToSliceID(protoSliceID *SliceID) (*types.SliceID, error) {
	sliceID := &types.SliceID{
		Region: protoSliceID.Region,
		Zone:   protoSliceID.Zone,
	}
	// TODO: implement
	return sliceID, nil
}

func convertProtoToHash(protoHash *Hash) (*common.Hash, error) {
	hash := &common.Hash{}
	hash.SetBytes(protoHash.Hash)
	// TODO: implement
	return hash, nil
}

func convertProtoToBlock(protoBlock *Block) (*types.Block, error) {

	txs := make([]*types.Transaction, len(protoBlock.Txs))
	for i, protoTransaction := range protoBlock.Txs {
		transaction, err := convertProtoToTransaction(protoTransaction)
		if err != nil {
			return nil, err
		}
		txs[i] = transaction
	}

	block := types.NewBlock(
		nil,
		txs,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
	)
	//! TODO: implement
	return block, nil
}

func convertProtoToHeader(protoHeader *Header) (*types.Header, error) {
	header := new(types.Header)
	header.SetGasLimit(protoHeader.GasLimit)
	header.SetGasUsed(protoHeader.GasUsed)
	// TODO: implement
	return header, nil
}

func convertProtoToTransaction(protoTransaction *Transaction) (*types.Transaction, error) {
	panic("TODO: implement")
	// return nil, nil
}
