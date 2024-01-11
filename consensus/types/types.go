package types

import (
	"encoding/hex"
	"strconv"

	"github.com/dominant-strategies/go-quai/p2p/pb"
	"github.com/pkg/errors"
)

//*Definitions and implementations for a Hash type*//

// Hash represents a 256 bit hash
type Hash [32]byte

// Converts a custom go Hash type (types.Hash) to a protocol buffer Hash type (pb.Hash)
func (h Hash) ToProto() *pb.Hash {

	return &pb.Hash{
		Hash: h[:],
	}
}

// converts a protocol buffer Hash type (pb.Hash) to a custom go Hash type (types.Hash)
func (h *Hash) FromProto(pbHash *pb.Hash) {

	copy(h[:], pbHash.Hash)
}

// returns a new instance of the protocol buffer Hash type (pb.Hash)
func (h Hash) NewProtoInstance() *pb.Hash {
	return &pb.Hash{}
}

func NewHashFromString(s string) (Hash, error) {
	if len(s) != 64 { // 2 hex chars per byte
		return Hash{}, errors.New("invalid string length for hash")
	}
	hashSlice, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}
	var hash Hash
	copy(hash[:], hashSlice)
	return hash, nil
}

//*Definitions and implementations for a Context type*//

type Context struct {
	Location string
	Level    uint32
}

func (c Context) ToProto() *pb.Context {
	return &pb.Context{
		Location: c.Location,
		Level:    c.Level,
	}
}

func (c *Context) FromProto(pbContext *pb.Context) {
	c.Location = pbContext.Location
	c.Level = pbContext.Level
}

func (c Context) NewProtoInstance() *pb.Context {
	return &pb.Context{}
}

var (
	PRIME_CTX  = Context{"prime", 0}
	REGION_CTX = Context{"region", 1}
	ZONE_CTX   = Context{"zone", 2}
)

// *Definitions and implementations for a custom PeerID type*//
type PeerID struct {
	Location    SliceID
	Entropy     uint32
	Zone_height uint32
	User_agent  string
}

//*Definitions and implementations for a SliceID type*//

type SliceID struct {
	Context Context
	Region  uint32
	Zone    uint32
}

func (s SliceID) ToProto() *pb.SliceID {
	return &pb.SliceID{
		Context: s.Context.ToProto(),
		Region:  s.Region,
		Zone:    s.Zone,
	}
}

func (sliceID SliceID) String() string {
	return strconv.Itoa(int(sliceID.Region)) + "." + strconv.Itoa(int(sliceID.Zone))
}

// converts a protocol buffer SliceID type (pb.SliceID) to a custom go SliceID type (types.SliceID)
func (sliceID *SliceID) FromProto(pbSliceID *pb.SliceID) {
	sliceID.Context.FromProto(pbSliceID.Context)
	sliceID.Region = pbSliceID.Region
	sliceID.Zone = pbSliceID.Zone
}

// returns a new instance of the protocol buffer SliceID type (pb.SliceID)
func (sliceID SliceID) NewProtoInstance() *pb.SliceID {
	return &pb.SliceID{}
}
