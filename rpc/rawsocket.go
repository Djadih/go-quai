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
	MAX_REQ_SIZE = 1024
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

func (miner *MinerSession) ListenTCP() {
	var accept = make(chan int, 1)
	n := 0

	accept <- n
	go func(ms *MinerSession) {
		err := ms.handleTCPClient(ms)
		if err != nil {
			ms.conn.Close()
		}
		<-accept
	}(miner)

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

func (ms *MinerSession) SendTCPRequest(msg jsonrpcMessage, resultCh chan *types.Header) error {

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
