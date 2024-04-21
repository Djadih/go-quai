package peerdb

import (
	sync "sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/syndtr/goleveldb/leveldb"
)

// contains the information of a peer
// that is stored in the peerDB as Value
type PeerInfo struct {
	AddrInfo  peer.AddrInfo
	PubKey    []byte
	Entropy   uint64
	Protected bool
}

type AddrInfo struct {
	peer.AddrInfo
}

// PeerDB implements the ipfs Datastore interface
// and exposes an API for storing and retrieving peers
// using levelDB as the underlying database
type PeerDB struct {
	db          *leveldb.DB
	peerCounter int
	mu          sync.Mutex
}

// ProtoEncode converts the hash into the ProtoHash type
func (pi PeerInfo) ProtoEncode() *ProtoPeerInfo {
	addrInfo := AddrInfo{pi.AddrInfo}

	return &ProtoPeerInfo{
		AddrInfo:  addrInfo.ProtoEncode(),
		PubKey:    pi.PubKey,
		Entropy:   pi.Entropy,
		Protected: pi.Protected,
	}
}

func (addr AddrInfo) ProtoEncode() *ProtoAddrInfo {
	multiAddrs := make([]string, len(addr.Addrs))
	for i, addr := range addr.Addrs {
		multiAddrs[i] = addr.String()
	}

	return &ProtoAddrInfo{
		ID:    addr.ID.String(),
		Addrs: multiAddrs,
	}
}

func (addr AddrInfo) ProtoDecode() *ProtoAddrInfo {
	multiAddrs := make([]string, len(addr.Addrs))
	for i, addr := range addr.Addrs {
		multiAddrs[i] = addr.String()
	}

	return &ProtoAddrInfo{
		ID:    addr.ID.String(),
		Addrs: multiAddrs,
	}
}
