package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type ClientConn struct {
	addr      string
	opts      options
	netid     int64
	rawConn   net.Conn
	once      *sync.Once
	wg        *sync.WaitGroup
	sendCh    chan []byte
	handlerCh chan MessageHandler
	// timing    *TimingWheel
	mu          sync.Mutex // guards following
	name        string
	heart       int64
	pending     []int64
	ctx         context.Context
	cancel      context.CancelFunc
	serviceType int64 //服务器类型
}

func NewClient(addr string, opt ...ServerOption) *ClientConn {
	c, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		//TODO
		return nil
	}

	var opts options
	for _, o := range opt {
		o(&opts)
	}
	//TODO
	// if opts.codec == nil {
	// 	opts.codec = TypeLengthValueCodec{}
	// }
	if opts.bufferSize <= 0 {
		opts.bufferSize = BufferSize256
	}
	return newClientConnWithOptions(c, opts)
}

func newClientConnWithOptions(c net.Conn, opts options) *ClientConn {
	cc := &ClientConn{
		addr: c.RemoteAddr().String(),
		opts: opts,
		// netid:     netid,
		rawConn:   c,
		once:      &sync.Once{},
		wg:        &sync.WaitGroup{},
		sendCh:    make(chan []byte, opts.bufferSize),
		handlerCh: make(chan MessageHandler, opts.bufferSize),
		heart:     time.Now().UnixNano(),
	}
	cc.ctx, cc.cancel = context.WithCancel(context.Background())
	// cc.timing = NewTimingWheel(cc.ctx)
	cc.name = c.RemoteAddr().String()
	cc.pending = []int64{}
	return cc
}

func (cc *ClientConn) NetID() int64 {
	return cc.netid
}

func (cc *ClientConn) SetServiceType(tp int64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.serviceType = tp
}

func (cc *ClientConn) Start() {
	onConnect := cc.opts.onConnect
	if onConnect != nil {
		onConnect(cc)
	}

}

func (cc *ClientConn) Close() {
	//TODO
}

// Write writes a message to the client.
func (cc *ClientConn) Write(message proto.Message) error {
	return cc.asyncWrite(message)
	// TODO
	// return nil
}

func (cc *ClientConn) asyncWrite(message proto.Message) (err error) {
	pkt, _ := proto.Marshal(message)

	select {
	case cc.sendCh <- pkt:
		err = nil
	default:
		err = ErrWouldBlock
	}
	return
}

func (cc *ClientConn) readLoop(c WriteCloser, wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// holmes.Errorf("panics: %v\n", p)
		}
		wg.Done()
		// holmes.Debugln("readLoop go-routine exited")
		cc.Close()
	}()

	for {
		select {
		case <-cc.ctx.Done(): // connection closed
			fmt.Println("receiving cancel signal from conn")
			return
		case <-cc.belong.ctx.Done(): // server closed
			fmt.Println("receiving cancel signal from server")
			return
		default:
			msg, err = codec.Decode(rawConn)
			if err != nil {
				fmt.Errorf("error decoding message %v\n", err)
				if _, ok := err.(ErrUndefined); ok {
					// update heart beats
					setHeartBeatFunc(time.Now().UnixNano())
					continue
				}
				return
			}
			setHeartBeatFunc(time.Now().UnixNano())
			handler := GetHandlerFunc(msg.MessageNumber())
			if handler == nil {
				if onMessage != nil {
					fmt.Printf("message %d call onMessage()\n", msg.MessageNumber())
					onMessage(msg, c.(WriteCloser))
				} else {
					fmt.Printf("no handler or onMessage() found for message %d\n", msg.MessageNumber())
				}
				continue
			}
			handlerCh <- MessageHandler{msg, handler}
		}
	}
}
