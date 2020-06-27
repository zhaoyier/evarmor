package network

import (
	"net"
	"os"

	"git.ezbuy.me/ezbuy/evarmor/common/log"
)

// func NewServer(addr string) {
// 	//初始化服务
// 	// initServer(addr)
// }

// func initServer(addr string) {
// 	server := &Server{
// 		addr:    addr,
// 		connMap: make(map[string]*UserConn),
// 	}
// 	server.NewListen(addr)
// }

func (s *Server) NewListen(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed: %q", err)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("accept error: %q", err)
			os.Exit(1)
		}
		log.Infof("message %s->%s", conn.RemoteAddr(), conn.LocalAddr())
		go s.HandleRequest(conn)
	}
}

func (s *Server) HandleRequest(conn net.Conn) {
	//TODO 记录session
	cn := NewConn(conn)
	defer func() { //TODO 监听关闭信号
		// log.Warnf("disconnect ip:%s", ip)
		conn.Close()
	}()

	//
	// s.connMap[cn.addr] = cn

	for {
		b, _, err := cn.reader.ReadLine()
		if err != nil {
			log.Fatalf("reader read line failed: %q", err)
			return
		}

		log.Infof("reader read line: %s", string(b))
		//TODO 回调接口
		// resp := Resp{
		// 	Data:   time.Now().String(),
		// 	Status: 200,
		// }
		// r, _ := json.Marshal(resp)
		// cn.writer.Write(r)
		// cn.writer.Write([]byte("\n"))
		// cn.writer.Flush()
	}
}

// // 创建连接
// func CreateClient() {
// 	//
// 	// 创建ticket

// 	t := time.NewTicker(time.Duration(time.Minute * 30))
// 	defer t.Stop()

// 	for { //保活
// 		switch {
// 		case <-t.C:

// 		}
// 		runtime.Gosched()
// 	}
// }

func Heartbeat() {
	for _, servers := range serverMap {
		for _, server := range servers {
			server.writer.Write([]byte("ping"))
			server.writer.Flush() //TODO 是否断开连接
			//TODO 断开重连
		}
	}
}
