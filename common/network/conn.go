package network

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	// "google.golang.org/protobuf/proto"
)

type MessageHandler struct {
	code   string
	data   string
	method *Method
}

type WriteCloser interface {
	Write(proto.Message) error
	Close()
}

type ServerConn struct {
	netid     int64
	name      string
	rawConn   net.Conn
	pending   []int64
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex // guards following
	belong    *Server
	wg        *sync.WaitGroup
	once      *sync.Once
	handlerCh chan MessageHandler
	reader    *bufio.Reader
	writer    *bufio.Writer
	sendCh    chan []byte
}

func NewServerConn(id int64, s *Server, conn net.Conn) *ServerConn {
	sc := &ServerConn{
		netid:     id,
		rawConn:   conn,
		belong:    s,
		wg:        &sync.WaitGroup{},
		once:      &sync.Once{},
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
		sendCh:    make(chan []byte, s.opts.bufferSize),
		handlerCh: make(chan MessageHandler, s.opts.bufferSize),
	}
	sc.ctx, sc.cancel = context.WithCancel(context.WithValue(s.ctx, serverCtx, s))
	sc.name = conn.RemoteAddr().String()
	sc.pending = []int64{}
	return sc
}

func (sc *ServerConn) SetName(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.name = name
}

func (sc *ServerConn) Start() {
	if sc.belong == nil {
		fmt.Printf("====>>1021:\n")
	}
	// if sc.belong.opts == nil {
	// 	fmt.Printf("====>>003:\n")
	// }
	onConnect := sc.belong.opts.onConnect
	if onConnect != nil {
		onConnect(sc)
	}

	go sc.readLoop()
	go sc.writeLoop()
	go sc.handleLoop()

	// loopers := []func(WriteCloser, *sync.WaitGroup){readLoop, writeLoop, handleLoop}
	// for _, l := range loopers {
	// 	looper := l
	// 	sc.wg.Add(1)
	// 	go looper(sc, sc.wg)
	// }
}

func (sc *ServerConn) Close() {

}

func (sc *ServerConn) Write(message proto.Message) error {
	return nil //TODO
	// return asyncWrite(sc, message)
}

func (sc *ServerConn) readLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("read panic: %q", p)
		}
		// sc.wg.Done()
		sc.Close()
	}()
	// var cDone <-chan struct{}
	// cDone = sc.ctx.Done()
	// var cDone <-chan struct{}
	// cDone = sc.ctx.Done()
	for {
		select {
		case <-sc.ctx.Done(): // connection closed
			fmt.Printf("receiving cancel signal from conn")
			return
		case <-sc.belong.ctx.Done():
			fmt.Printf("receiving cancel signal from server")
			return
		default:
			//读消息并回调接口
			// b, _, err := sc.reader.ReadLine()
			// if err != nil {
			// 	log.Fatalf("reader read line failed: %q", err)
			// 	return
			// }
			buf := make([]byte, 1024)
			reqLen, _ := sc.rawConn.Read(buf)
			fmt.Printf("====>>>data: %+v|%+v\n", string(buf[:reqLen]), reqLen)

			xm := &XMessage{}
			if err := proto.Unmarshal(buf[:reqLen], xm); err != nil {
				fmt.Printf("====>>0021:%q\n", err)
				return
			}

			fmt.Printf("====>>0022:%+v\n", xm)

			onMessage := sc.belong.opts.onMessage
			handler := sc.belong.GetHandlerFunc(xm.GetCode())
			if handler == nil {
				if onMessage != nil {
					onMessage(xm, sc)
				} else {
					fmt.Printf("no handler or onMessage() found for message %d\n", xm.GetCode())
				}
			}
			sc.handlerCh <- MessageHandler{xm.GetCode(), xm.GetData(), handler}
		}
	}
}

func (sc *ServerConn) writeLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("write loop panics: %v\n", p)
		}
	}()
	for {
		select {
		case <-sc.ctx.Done(): // connection closed
			fmt.Printf("receiving cancel signal from conn")
			return
		case <-sc.belong.ctx.Done():
			fmt.Printf("receiving cancel signal from server")
			return

		case pkt := <-sc.sendCh:
			if pkt != nil {
				// _, err := sc.writer.Write(pkt)
				// if err != nil {

				// }
				// sc.writer.Flush() //TODO 是否断开连接

				if _, err := sc.rawConn.Write(pkt); err != nil {
					fmt.Printf("write loop data: %q\n", err)
					return
				}
			}
		}
	}
}

func (sc *ServerConn) handleLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("handle panic: %q", p)
		}
		// sc.wg.Done()
		sc.Close()
	}()
	for {
		select {
		case <-sc.ctx.Done():
			fmt.Printf("receiving cancel signal from conn")
		case <-sc.belong.ctx.Done():
			fmt.Printf("receiving cancel signal from server")
		case hc := <-sc.handlerCh:
			fmt.Printf("handle do :%+v\n", hc)
			ctx := sc.ctx
			msg, method := hc.data, hc.method
			fmt.Printf("====>>0023:%+v|%+v\n", msg, method)

			// var req method.ParamType
			// req := &evarmor.HelloRequest{}
			req := reflect.New(method.ParamType.Elem()).Interface().(proto.Message)
			fmt.Printf("====>>0024:%+v|%+v\n", []byte(msg), method.ParamType)
			if err := proto.Unmarshal([]byte(msg), req); err != nil {
				fmt.Printf("proto unmarshal failed :%+v", hc)
				return
			}
			fmt.Printf("====>>0025:%+v|%+v\n", req, nil)

			// handler.Method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
			method.Method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})

		}
	}
}
