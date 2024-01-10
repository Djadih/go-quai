package types

import (
	"encoding/hex"
	"strconv"

	"github.com/pkg/errors"
)

// Hash represents a 256 bit hash
type Hash [32]byte

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

type Context struct {
	Location string
	Level    uint32
}

var (
	PRIME_CTX  = Context{"prime", 0}
	REGION_CTX = Context{"region", 1}
	ZONE_CTX   = Context{"zone", 2}
)

type SliceID struct {
	Context Context
	Region  uint32
	Zone    uint32
}

type PeerID struct {
	// TODO: Evaluate if entropy and zone_height should stay as uint32 or need to be big.int
	Location    *SliceID
	Entropy     uint32
	Zone_height uint32
	User_agent  string
}

func (sliceID SliceID) String() string {
	return strconv.Itoa(int(sliceID.Region)) + "." + strconv.Itoa(int(sliceID.Zone))
}
