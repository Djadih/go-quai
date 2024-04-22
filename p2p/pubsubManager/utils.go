package pubsubManager

import (
	"strings"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
)

const (
	// Data types for gossipsub topics
	C_workObjectType  = "blocks"
	C_transactionType = "transactions"
	C_headerType      = "headers"
	C_hashType        = "hash"
)

// gets the name of the topic for the given type of data
func TopicName(genesis common.Hash, location common.Location, data interface{}) (string, error) {
	baseTopic := strings.Join([]string{genesis.String(), location.Name()}, "/")
	switch data.(type) {
	case *types.WorkObject:
		return strings.Join([]string{baseTopic, C_workObjectType}, "/"), nil
	case common.Hash:
		return strings.Join([]string{baseTopic, C_hashType}, "/"), nil
	case *types.Transaction:
		return strings.Join([]string{baseTopic, C_transactionType}, "/"), nil
	default:
		return "", ErrUnsupportedType
	}
}

func getTopicType(topic string) string {
	return topic[strings.LastIndex(topic, "/")+1:]
}

// lists our peers which provide the associated topic
func (g *PubsubManager) PeersForTopic(location common.Location, data interface{}) ([]peer.ID, error) {
	topicName, err := TopicName(g.genesis, location, data)
	if err != nil {
		return nil, err
	}
	return g.topics[topicName].ListPeers(), nil
}

// Creates a Cid from a location to be used as DHT key
func LocationToCid(location common.Location) cid.Cid {
	sliceBytes := []byte(location.Name())

	// create a multihash from the slice ID
	mhash, _ := multihash.Encode(sliceBytes, multihash.SHA2_256)

	// create a Cid from the multihash
	return cid.NewCidV1(cid.Raw, mhash)
}
