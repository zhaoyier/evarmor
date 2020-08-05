package network

import (
	"bufio"
	"crypto/tls"
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
