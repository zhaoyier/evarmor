package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
)

type ServerOption func(*options)

type Server struct {
	desc             string
	opts             options
	ctx              context.Context
	cancel           context.CancelFunc
	conns            *sync.Map
	plugins          *sync.Map  //其他服务器
	mu               sync.Mutex // guards following
	lis              map[string]net.Listener
	wg               *sync.WaitGroup
	delay            time.Duration
	serviceMap       map[string]*Method
	servicePluginMap map[string][]*ServerConn
	serviceType      int64 //服务器类型
}

// type Codec interface {
// 	Decode(net.Conn) (Message, error)
// 	Encode(Message) ([]byte, error)
// }

func NewServer(desc string, opt ...ServerOption) *Server {
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
		desc:             desc,
		opts:             opts,
		conns:            &sync.Map{},
		plugins:          &sync.Map{},
		wg:               &sync.WaitGroup{},
		lis:              make(map[string]net.Listener),
		serviceMap:       make(map[string]*Method),
		servicePluginMap: make(map[string][]*ServerConn),
		serviceType:      10000, //TODO 默认
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *Server) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
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
		//TODO 是否监听服务关闭
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
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				//TODO 超时处理
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

		netid := getAndIncrement(s.serviceType)
		sc := NewServerConn(netid, s, conn)
		sc.SetName(sc.rawConn.RemoteAddr().String())
		s.conns.Store(netid, sc)

		s.wg.Add(1) // this will be Done() in ServerConn.Close()
		go func() {
			sc.Start()
		}()
	}
}

func (s *Server) SetServiceType(tp int64) {
	s.serviceType = tp
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

// 添加插件服务
func (s *Server) AddPlugin(id int64, sc *ServerConn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.conns.Delete(id)
	// s.plugins.Store(id, sc)
}

func (s *Server) SyncPlugin(data string, sc *ServerConn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.conns.Delete(sc.netid)

	var ping Ping
	if err := proto.Unmarshal([]byte(data), &ping); err != nil {
		// TODO 处理错误
		return
	}

	for _, name := range ping.GetNames() {
		conns, ok := s.servicePluginMap[name]
		if ok {
			for _, conn := range conns {
				if conn.addr == sc.addr {
					continue
				}
			}
			conns = append(conns, sc)
			s.servicePluginMap[name] = conns
			continue
		}
		conns = make([]*ServerConn, 0, 2)
		conns = append(conns, sc)
		s.servicePluginMap[name] = conns
	}
}

func (s *Server) GetPlugin(service string, sc *ServerConn) (*ServerConn, error) {
	conns, ok := s.servicePluginMap[service]
	if !ok {
		return nil, errors.New("")
	}
	index := sc.hashNum % len(conns)
	return conns[index], nil
}
