package network

import (
	"bufio"
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"
	// "google.golang.org/protobuf/proto"
)

type options struct {
	tlsCfg *tls.Config
	// codec      Codec
	onConnect  onConnectFunc
	onMessage  onMessageFunc
	onClose    onCloseFunc
	onError    onErrorFunc
	workerSize int  // numbers of worker go-routines
	bufferSize int  // size of buffered channel
	reconnect  bool // for ClientConn use only
}

type Server struct {
	opts       options
	ctx        context.Context
	cancel     context.CancelFunc
	conns      *sync.Map
	mu         sync.Mutex // guards following
	lis        map[string]net.Listener
	wg         *sync.WaitGroup
	delay      time.Duration
	serviceMap map[string]*Method
}

// 分割线

type UserConn struct {
	addr       string
	connId     int64
	serverType int32
	reader     *bufio.Reader
	writer     *bufio.Writer
	readerC    int64 //读取消息数
	writerC    int64 //发送消息数
}
