package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/common/hexutil"
	"github.com/dominant-strategies/go-quai/common"
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
	timeout  time.Duration
}

const (
	MAX_REQ_SIZE = 8192
)

func NewMinerConn(endpoint string) (*MinerSession, error) {
	localaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:15000")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	remoteaddr, err := net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// server, err := net.Dial("tcp", endpoint)
	server, err := net.DialTCP("tcp", localaddr, remoteaddr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	server.SetDeadline(time.Time{})
	// defer server.Close()
	

	log.Printf("New TCP client made to: %v", server)

	// conn, err := server.AcceptTCP()
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// 	return nil, err
	// }
	// conn.SetKeepAlive(true)

	// ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())

	return &MinerSession{proto: "tcp", ip: remoteaddr.AddrPort().Addr().String(), port: "15000", conn: server, latestId: 0}, nil
	// return &MinerSession{proto: "tcp", ip: "", port: "", conn: nil, latestId: 0}, nil
}

// RPCMarshalHeader converts the given header to the RPC output .
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

type JSONRpcResp struct {
	Id      json.RawMessage `json:"id"`
	Version string          `json:"jsonrpc"`
	Result  types.Header    `json:"result"`
	Error   interface{}     `json:"error,omitempty"`
}


func (miner *MinerSession) ListenTCP(updateCh chan *types.Header) error {
	connbuff := bufio.NewReaderSize(miner.conn, MAX_REQ_SIZE)

	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			log.Printf("Socket flood detected from %s", miner.ip)
			// mienr.policy.BanClient(cs.ip)
			return err
		} else if err == io.EOF {
			log.Printf("Client %s disconnected", miner.ip)
			// s.removeSession(miner)
			break
		} else if err != nil {
			log.Printf("Error reading from socket: %v", err)
			return err
		}


		if len(data) > 1 {
			var rpcResp JSONRpcResp
			err = json.Unmarshal(data, &rpcResp)
			if err != nil {
			log.Printf("Unable to decode header: %v", err)
			return err
			}
			// var header *types.Header
			// json.Unmarshal(rpcResp.Result, &header)
			// header := *types.Header(rpcResp.Result)
			 updateCh <- &rpcResp.Result
		}
	}
	return nil
}

func (miner *MinerSession) handleTCPClient(ms *MinerSession) error {
	ms.enc = json.NewEncoder(ms.conn)
	connbuff := bufio.NewReaderSize(ms.conn, MAX_REQ_SIZE)
	// s.setDeadline(cs.conn)
	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			log.Printf("Socket flood detected from %s", miner.ip)
			// mienr.policy.BanClient(cs.ip)
			return err
		} else if err == io.EOF {
			log.Printf("Client %s disconnected", miner.ip)
			// s.removeSession(miner)
			break
		} else if err != nil {
			log.Printf("Error reading from socket: %v", err)
			return err
		}

		if len(data) > 1 {
			var req StratumReq
			err = json.Unmarshal(data, &req)
			if err != nil {
				// s.policy.ApplyMalformedPolicy(cs.ip)
				log.Printf("Malformed stratum request from %s: %v", ms.ip, err)
				return err
			}
			// s.setDeadline(cs.conn)
			err = ms.handleTCPMessage(&req)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ms *MinerSession) sendTCPResult(result json.RawMessage) error {

	ms.Lock()
	defer ms.Unlock()

	// ms.latestId += 1

	message, err := json.Marshal(jsonrpcMessage{ID: json.RawMessage(strconv.FormatUint(ms.latestId, 10)), Version: "2.0", Error: nil, Result: result})
	if err != nil {
		return err
	}

	ms.conn.Write(message)
	return nil
}

func (ms *MinerSession) SendTCPRequest(msg jsonrpcMessage) error {

	ms.Lock()
	defer ms.Unlock()

	ms.latestId += 1
	// msg.ID = json.RawMessage(ms.latestId)
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ms.conn.Write(message)
	ms.conn.Write([]byte("\n"))
	// header, _ := ms.conn.Read()
	// resultCh <- header
	return nil
}

func (ms *MinerSession) handleTCPMessage(req *StratumReq) error {
	// Handle RPC methods
	// switch req.Message.Method {
	// case "quai_getPendingHeader":
	// 	// reply, errReply := s.handleGetWorkRPC(cs)
	// 	reply, err := api.
	// 	if errReply != nil {
	// 		return cs.sendTCPError(req.Id, errReply)
	// 	}
	// 	return cs.sendTCPResult(req.Id, &reply)
	// }
	// Println(req.Message)

	return nil
}

// type tcpConn struct {
// 	in	io.Reader
// 	out	io.Writer
// }

// func DialTCP(ctx context.Context, endpoint string) (*Client, error) {
// 	return DialTCPIO(ctx, endpoint, tcpConn)
// 	// base_rpc.Dial("tcp", endpoint)

// 	// new_client, err :=
// 	// new_client, err := newClient(ctx, func(_ context.Context) (ServerCodec, error) {
// 	// 	return NewCodec(stdioConn{
// 	// 		in: in,
// 	// 		out: out,
// 	// 	}), nil
// 	// })

// 	// return new_client, err
// }

// func DialTCPIO(ctx context.Context, in io.Reader, out io.Writer) (*Client, error) {
// 	return nil, nil
// }
