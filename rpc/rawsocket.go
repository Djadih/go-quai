package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/INFURA/go-ethlibs/jsonrpc"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/common/hexutil"
	"github.com/dominant-strategies/go-quai/core/types"
)

type MinerSession struct {
	proto string
	ip    string
	port  string
	conn  *net.TCPConn
	enc   *json.Encoder

	// Stratum
	sync.Mutex
	latestId uint64
}

const (
	MAX_REQ_SIZE = 8192
)

func NewMinerConn(endpoint string) (*MinerSession, error) {
	localaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	remoteaddr, err := net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}

	server, err := net.DialTCP("tcp", localaddr, remoteaddr)
	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}
	server.SetDeadline(time.Time{})

	log.Printf("New TCP client made to: %v", server)

	return &MinerSession{proto: "tcp", ip: remoteaddr.AddrPort().Addr().String(), port: "15000", conn: server, latestId: 0, enc: json.NewEncoder(server)}, nil
}

// RPCMarshalHeader converts the given header to the RPC output.
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	result := map[string]interface{}{
		"hash":                head.Hash(),
		"parentHash":          head.ParentHashArray(),
		"nonce":               head.Nonce(),
		"sha3Uncles":          head.UncleHashArray(),
		"logsBloom":           head.BloomArray(),
		"stateRoot":           head.RootArray(),
		"miner":               head.CoinbaseArray(),
		"extraData":           hexutil.Bytes(head.Extra()),
		"size":                hexutil.Uint64(head.Size()),
		"timestamp":           hexutil.Uint64(head.Time()),
		"transactionsRoot":    head.TxHashArray(),
		"receiptsRoot":        head.ReceiptHashArray(),
		"extTransactionsRoot": head.EtxHashArray(),
		"extRollupRoot":       head.EtxRollupHashArray(),
		"manifestHash":        head.ManifestHashArray(),
		"location":            head.Location(),
	}

	number := make([]*hexutil.Big, common.HierarchyDepth)
	difficulty := make([]*hexutil.Big, common.HierarchyDepth)
	gasLimit := make([]hexutil.Uint, common.HierarchyDepth)
	gasUsed := make([]hexutil.Uint, common.HierarchyDepth)
	for i := 0; i < common.HierarchyDepth; i++ {
		number[i] = (*hexutil.Big)(head.Number(i))
		difficulty[i] = (*hexutil.Big)(head.Difficulty(i))
		gasLimit[i] = hexutil.Uint(head.GasLimit(i))
		gasUsed[i] = hexutil.Uint(head.GasUsed(i))
	}
	result["number"] = number
	result["difficulty"] = difficulty
	result["gasLimit"] = gasLimit
	result["gasUsed"] = gasUsed

	if head.BaseFee() != nil {
		results := make([]*hexutil.Big, common.HierarchyDepth)
		for i := 0; i < common.HierarchyDepth; i++ {
			results[i] = (*hexutil.Big)(head.BaseFee(i))
		}
		result["baseFeePerGas"] = results
	}

	return result
}

func (miner *MinerSession) ListenTCP(updateCh chan *types.Header) error {
	connbuff := bufio.NewReaderSize(miner.conn, MAX_REQ_SIZE)

	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			log.Printf("Socket flood detected from %s", miner.ip)
			// miner.policy.BanClient(cs.ip)
			return err
		} else if err == io.EOF {
			log.Printf("Client %s disconnected", miner.ip)
			// miner.removeSession(miner)
			break
		} else if err != nil {
			log.Printf("Error reading from socket: %v", err)
			return err
		}

		if len(data) > 1 {
			var rpcResp *JsonRPCResponse
			err := json.Unmarshal(data, &rpcResp)

			if err != nil {
				log.Printf("Unable to decode RPC Response: %v", err)
				return err
			}
			var header *types.Header
			err = json.Unmarshal(rpcResp.Result, &header)
			if err != nil {
				log.Printf("Unable to decode header: %v", err)
				return err
			}

			updateCh <- header
		}
	}
	return nil
}

func (ms *MinerSession) SendTCPRequest(msg jsonrpc.Request) error {

	ms.Lock()
	defer ms.Unlock()

	// ms.latestId += 1
	message, err := msg.MarshalJSON()
	if err != nil {
		log.Fatalf("Error: %v", err)
		return err
	}

	// return ms.enc.Encode(&message)

	ms.conn.Write(message)
	ms.conn.Write([]byte("\n"))
	// header, _ := ms.conn.Read()
	// resultCh <- header
	return nil
}
