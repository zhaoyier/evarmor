package network

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
)

var (
	connId          int64 = 10000000
	userMapOnce     sync.Once
	userSessionMap  map[string]*UserConn
	userSessionOnce sync.Once
	serverMapOnce   sync.Once
	serverMap       map[string][]*UserConn
)

func init() {
	userMapOnce.Do(func() {
		userSessionMap = make(map[string]*UserConn)
	})
	serverMapOnce.Do(func() {
		serverMap = make(map[string][]*UserConn)
	})

}

func NewConn(conn net.Conn) *UserConn {
	uc := &UserConn{
		connId: atomic.AddInt64(&connId, 1),
		addr:   conn.RemoteAddr().String(),
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
	userSessionMap[uc.addr] = uc
	return uc
}

func (uc *UserConn) Dispatch(msg []byte) {
	switch uc.serverType {
	case 0: //gateway
		//TODO 寻址转发
	case 1: //logic
		//TODO寻找具体接口
	}
}

func CreateClient(addr string) {

}
