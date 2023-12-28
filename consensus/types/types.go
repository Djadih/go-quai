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
	var h Hash
	_, err := hex.Decode(h[:], []byte(s))
	if err != nil {
		return Hash{}, err
	}
	return h, nil
}

type Context struct {
	location string
	level    int
}

var (
	PRIME_CTX  = Context{"prime", 0}
	REGION_CTX = Context{"region", 1}
	ZONE_CTX   = Context{"zone", 2}
)

type SliceID struct {
	Context Context
	Region  int
	Zone    int
}

func (sliceID SliceID) String() string {
	return strconv.Itoa(sliceID.Region) + "." + strconv.Itoa(sliceID.Zone)
}
