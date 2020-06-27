package network

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"
)

type ServerOption func(*options)

// type Codec interface {
// 	Decode(net.Conn) (Message, error)
// 	Encode(Message) ([]byte, error)
// }

func NewServer(opt ...ServerOption) *Server {
	var opts options
	for _, o := range opt {
		o(&opts)
	}
	if opts.workerSize <= 0 {
		opts.workerSize = defaultWorkersNum
	}
	if opts.bufferSize <= 0 {
		opts.bufferSize = BufferSize256
	}

	s := &Server{
		opts:       opts,
		conns:      &sync.Map{},
		wg:         &sync.WaitGroup{},
		lis:        make(map[string]net.Listener),
		serviceMap: make(map[string]*Method),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func newListen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func (s *Server) Start(addr string) error {
	l, err := newListen(addr)
	if err != nil {
		fmt.Printf("new listen failed: %q", err)
		return err
	}
	s.lis[addr] = l
	defer func() {
		if _, ok := s.lis[addr]; ok {
			s.lis[addr].Close()
			delete(s.lis, addr)
		}
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if s.delay == 0 {
					s.delay = 5 * time.Millisecond
				} else {
					s.delay *= 2
				}
				if max := 1 * time.Second; s.delay >= max {
					s.delay = max
				}
				select {
				case <-time.After(s.delay):
				case <-s.ctx.Done():
				}
				continue
			}
			return err
		}
		s.delay = 0
		sz := s.ConnsSize()
		if sz > MaxConnections {
			fmt.Printf("max ocnnnections size: %d, refused", MaxConnections)
			conn.Close()
			continue
		}

		netid := getAndIncrement()
		sc := NewServerConn(netid, s, conn)
		sc.SetName(sc.rawConn.RemoteAddr().String())
		s.conns.Store(netid, sc)

		s.wg.Add(1) // this will be Done() in ServerConn.Close()
		go func() {
			sc.Start()
		}()
	}
}

func (s *Server) Stop() {

}

func (s *Server) RegisterServer(srv Service) {
	t := reflect.TypeOf(srv)
	v := reflect.ValueOf(srv)
	if t.NumMethod() == 0 {
		// TODO 临时注销
		// panic("no method found for serivce: " + t.Name())
	}

	for i := 0; i < t.NumMethod(); i++ {
		name := t.Method(i).Name
		_, ok := s.serviceMap[name]
		if ok {
			panic("duplicate register service:" + name)
		}
		s.serviceMap[name] = &Method{
			Method:    v.Method(i),
			ParamType: t.Method(i).Type.In(2),
		}
	}

}

func (s *Server) ConnsSize() int {
	var sz int
	s.conns.Range(func(k, v interface{}) bool {
		sz++
		return true
	})
	return sz
}

func (s *Server) GetHandlerFunc(name string) *Method {
	entry, ok := s.serviceMap[name]
	if !ok {
		return nil
	}
	return entry
}
