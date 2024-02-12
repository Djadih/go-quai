// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"io"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/common/hexutil"
	"github.com/dominant-strategies/go-quai/rlp"
)

//go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

// Logs is a list of log objects.
type Logs []*Log

type logMarshaling struct {
	Data        hexutil.Bytes
	BlockNumber hexutil.Uint64
	TxIndex     hexutil.Uint
	Index       hexutil.Uint
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

// rlpStorageLog is the storage encoding of a log.
type rlpStorageLog rlpLog

// legacyRlpStorageLog is the previous storage encoding of a log including some redundant fields.
type legacyRlpStorageLog struct {
	Address     common.Address
	Topics      []common.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      common.Hash
	TxIndex     uint
	BlockHash   common.Hash
	Index       uint
}

// EncodeRLP implements rlp.Encoder.
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

// DecodeRLP implements rlp.Decoder.
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage Log

// ProtoEncode converts the log to a protobuf representation.
func (l LogForStorage) ProtoEncode() *ProtoLogForStorage {
	address := l.Address.ProtoEncode()
	topics := make([]*common.ProtoHash, len(l.Topics))
	for i, t := range l.Topics {
		topics[i] = t.ProtoEncode()
	}
	return &ProtoLogForStorage{
		Address: address,
		Topics:  topics,
		Data:    l.Data,
	}
}

// ProtoDecode converts the protobuf to a log representation.
func (l *LogForStorage) ProtoDecode(protoLog *ProtoLogForStorage, location common.Location) error {
	address := new(common.Address)
	err := address.ProtoDecode(protoLog.GetAddress(), location)
	if err != nil {
		return err
	}
	topics := make([]common.Hash, len(protoLog.Topics))
	for i, t := range protoLog.Topics {
		topic := new(common.Hash)
		topic.ProtoDecode(t)
		topics[i] = *topic
	}
	l.Address = *address
	l.Topics = topics
	l.Data = protoLog.GetData()
	return nil
}

// EncodeRLP implements rlp.Encoder.
func (l *LogForStorage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpStorageLog{
		Address: l.Address,
		Topics:  l.Topics,
		Data:    l.Data,
	})
}

// DecodeRLP implements rlp.Decoder.
//
// Note some redundant fields(e.g. block number, tx hash etc) will be assembled later.
func (l *LogForStorage) DecodeRLP(s *rlp.Stream) error {
	blob, err := s.Raw()
	if err != nil {
		return err
	}
	var dec rlpStorageLog
	err = rlp.DecodeBytes(blob, &dec)
	if err == nil {
		*l = LogForStorage{
			Address: dec.Address,
			Topics:  dec.Topics,
			Data:    dec.Data,
		}
	} else {
		// Try to decode log with previous definition.
		var dec legacyRlpStorageLog
		err = rlp.DecodeBytes(blob, &dec)
		if err == nil {
			*l = LogForStorage{
				Address: dec.Address,
				Topics:  dec.Topics,
				Data:    dec.Data,
			}
		}
	}
	return err
}
