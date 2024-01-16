package protocol

import (
	"context"
	"testing"
	"time"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/p2p/pb"
	"github.com/dominant-strategies/go-quai/p2p/protocol/mocks"
	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p/core/network"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuaiProtocolHandler(t *testing.T) {

	sliceID := types.SliceID{
		Region: 0,
		Zone:   0,
	}

	hash := common.Hash{}
	// set hash with all 0x11 bytes
	for i := 0; i < 32; i++ {
		hash[i] = 0x11
	}
	header := new(types.Header)
	header.SetGasLimit(1000)
	header.SetGasUsed(100)

	testCases := []struct {
		name      string
		mockStub  func(*mocks.MockQuaiP2PNode)
		action    pb.QuaiRequestMessage_ActionType
		slice     *types.SliceID
		hash      *common.Hash
		ExpectErr bool
	}{
		{
			name:   "request header",
			action: pb.QuaiRequestMessage_REQUEST_HEADER,
			slice:  &sliceID,
			hash:   &hash,
			mockStub: func(mockedQuaiNode *mocks.MockQuaiP2PNode) {
				mockedQuaiNode.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(header).Times(1)
			},

			ExpectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a mock network
			mnet := mocknet.New()

			// Setup two nodes on the mock network
			host1, err := mnet.GenPeer()
			require.NoError(t, err)
			defer host1.Close()

			host2, err := mnet.GenPeer()
			require.NoError(t, err)
			defer host2.Close()

			// Connect the two nodes
			err = mnet.LinkAll()
			require.NoError(t, err)
			err = mnet.ConnectAllButSelf()
			require.NoError(t, err)

			// Set up the protocol handler on the second node with the mockedQuaiNode
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedQuaiNode := mocks.NewMockQuaiP2PNode(ctrl)
			tc.mockStub(mockedQuaiNode)
			host2.SetStreamHandler(ProtocolVersion, func(s network.Stream) {
				QuaiProtocolHandler(s, mockedQuaiNode)
			})

			// Set up a stream between the two hosts
			stream, err := host1.NewStream(ctx, host2.ID(), ProtocolVersion)
			assert.NoError(t, err)
			defer stream.Close()

			// Encode the message to send
			data, err := pb.EncodeQuaiRequest(tc.action, tc.slice, tc.hash)
			require.NoError(t, err)

			// Send the message
			_, err = stream.Write(data)
			// sleep for a bit to allow the stream to be read
			time.Sleep(1 * time.Second)

			if tc.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// read the response

			readData := make([]byte, 256)
			_, err = stream.Read(readData)
			// sleep for a bit to allow the stream to be read
			time.Sleep(1 * time.Second)

			if tc.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
