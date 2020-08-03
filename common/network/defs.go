package network

import (
	"context"
	"errors"
	"reflect"
	"time"

	// "google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

const (
	MaxConnections    = 1000
	BufferSize128     = 128
	BufferSize256     = 256
	BufferSize512     = 512
	BufferSize1024    = 1024
	defaultWorkersNum = 20
)

var (
	ErrWouldBlock = errors.New("would block")
)

const (
	serverCtx contextKey = "server"
	netIDCtx  contextKey = "netid"
)

var (
	serviceMap = make(map[string]Method)
)

type contextKey string

type onScheduleFunc func(time.Time, WriteCloser)

// 连接回调
type onConnectFunc func(WriteCloser) bool

//消息通知
type onMessageFunc func(proto.Message, WriteCloser)

//关闭通知
type onCloseFunc func(WriteCloser)

// 工作池
type workerFunc func()

//错误通知
type onErrorFunc func(WriteCloser)

var (
	globalWorkerPool *WorkerPool
)

type Method struct {
	Method    reflect.Value
	ParamType reflect.Type //XXXXRequest的实际类型
}

// HandlerFunc serves as an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(context.Context, WriteCloser)

type Service interface {
}
