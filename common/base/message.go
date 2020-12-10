package base

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"path"
	"reflect"

	"git.ezbuy.me/ezbuy/evarmor/common/log"
	// "git.ezbuy.me/ezbuy/base/misc/log"
	"github.com/golang/protobuf/proto"
	"github.com/leesper/holmes"
)

const (
	// HeartBeat is the default heart beat message number.
	HeartBeat        = 0
	ProxyMessageType = -1
)

// Handler takes the responsibility to handle incoming messages.
type Handler interface {
	Handle(context.Context, interface{})
}

// HandlerFunc serves as an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(context.Context, WriteCloser)
type HandlerFunc2 func(context.Context, proto.Message)

// Handle calls f(ctx, c)
func (f HandlerFunc) Handle(ctx context.Context, c WriteCloser) {
	f(ctx, c)
}

// Handle calls f(ctx, c)
func (f HandlerFunc2) Handle(ctx context.Context, c proto.Message) {
	f(ctx, c)
}

// UnmarshalFunc unmarshals bytes into Message.
type UnmarshalFunc func([]byte) (Message, error)
type UnmarshalFunc2 func([]byte) (Message, error)

// handlerUnmarshaler is a combination of unmarshal and handle functions for message.
type handlerUnmarshaler struct {
	handler     HandlerFunc
	unmarshaler UnmarshalFunc
}

//handlerUnmarshaler2 service的方法签名满足 func(context.Context, XXXXRequest) Response
//XXXXRequest 实现 Request interface
type handlerUnmarshaler2 struct {
	Method    reflect.Value
	ParamType reflect.Type //XXXXRequest的实际类型
	OutType   reflect.Type //XXXXRequest的实际类型
}

var (
	buf *bytes.Buffer
	// messageRegistry is the registry of all
	// message-related unmarshal and handle functions.
	messageRegistry  map[int32]handlerUnmarshaler
	messageRegistry2 map[string]*handlerUnmarshaler2
)

func init() {
	// messageRegistry = map[int32]handlerUnmarshaler{-1: {
	// 	unmarshaler: _deserializeMessage,
	// 	handler:     _processMessage,
	// }}
	messageRegistry = make(map[int32]handlerUnmarshaler)
	fmt.Printf("====>>msg 01：%+v\n", len(messageRegistry))
	messageRegistry2 = make(map[string]*handlerUnmarshaler2)
	buf = new(bytes.Buffer)
}

// Register registers the unmarshal and handle functions for msgType.
// If no unmarshal function provided, the message will not be parsed.
// If no handler function provided, the message will not be handled unless you
// set a default one by calling SetOnMessageCallback.
// If Register being called twice on one msgType, it will panics.
func Register(msgType int32, unmarshaler func([]byte) (Message, error), handler func(context.Context, WriteCloser)) {
	if _, ok := messageRegistry[msgType]; ok {
		panic(fmt.Sprintf("trying to register message %d twice", msgType))
	}

	messageRegistry[msgType] = handlerUnmarshaler{
		unmarshaler: unmarshaler,
		handler:     HandlerFunc(handler),
	}
}

func Dispatch(msgType int32, handler func(context.Context, WriteCloser)) {
	entry, ok := messageRegistry[msgType]
	if !ok {
		entry = handlerUnmarshaler{
			unmarshaler: _deserializeMessage,
			handler:     HandlerFunc(handler),
		}
	}
	entry.handler = HandlerFunc(handler)
	messageRegistry[msgType] = entry
}

func RegisterService(srvs ...interface{}) { //
	// fmt.Printf("====>>500: %+v\n", GetServiceName(srv))
	// name := GetServiceName(srv)
	for _, srv := range srvs {
		log.Infof("%+v\n", GetServiceName(srv))

		refTyp := reflect.TypeOf(srv)
		refVal := reflect.ValueOf(srv)
		if refTyp.NumMethod() == 0 {
			panic("no method found for serivce: " + refTyp.Name())
		}
		for m := 0; m < refTyp.NumMethod(); m++ {
			method := refTyp.Method(m)
			fmt.Printf("=====>>1000:%+v|%+v|%+v\n", method.Name, method.PkgPath, refTyp.Method(m).Type.Out(0))
			if _, ok := messageRegistry2[method.Name]; ok {
				panic("duplicate register service:" + method.Name)
			}
			messageRegistry2[method.Name] = &handlerUnmarshaler2{
				Method:    refVal.Method(m),
				ParamType: refTyp.Method(m).Type.In(2),
				OutType:   refTyp.Method(m).Type.Out(0),
			}
		}
	}

	// for m := 0; i < tf.NumMethod(); i++ {
	// 	hash := utils.GetNameHash(tf.Method(i).Name)
	// 	fmt.Printf("=====>>505:%+v|%+v|%+v\n", tf.Method(i).Name, vf.Method(i), hash)
	// 	method := tf.Method(i)
	// 	// name := tf.Method(i).Name
	// 	// msgType := utils.CRC32(name)
	// 	// if _, ok := messageRegistry[msgType]; ok {
	// 	// 	panic("duplicate register service:" + tf.Method(i).Name)
	// 	// }
	// 	// messageRegistry[msgType] = &handlerUnmarshaler{
	// 	// 	Method:    vf.Method(i),
	// 	// 	ParamType: tf.Method(i).Type.In(2),
	// 	// }
	// 	// // TODO 利用etcd注册服务，注册地址信息
	// 	// etcdSer.PutService(fmt.Sprintf(RegisterServiceHandler, name), "heiheihei")
	// }
}

// GetUnmarshalFunc returns the corresponding unmarshal function for msgType.
func GetUnmarshalFunc(msgType int32) UnmarshalFunc {
	entry, ok := messageRegistry[msgType]
	if !ok {
		entry, ok := messageRegistry[-1]
		if !ok {
			return nil
		}
		return entry.unmarshaler
	}

	return entry.unmarshaler
}

// GetDefaultUnmarshalFunc returns the corresponding unmarshal function for msgType.
func GetDefaultUnmarshalFunc() UnmarshalFunc {
	entry, ok := messageRegistry[-1]
	if !ok {
		return nil
	}
	return entry.unmarshaler
}

// GetHandlerFunc returns the corresponding handler function for msgType.
func GetHandlerFunc(msgType int32) HandlerFunc {
	entry, ok := messageRegistry[msgType]
	if !ok {
		entry, ok := messageRegistry[-1]
		if !ok {
			return nil
		}
		return entry.handler
	}
	return entry.handler
}

// GetDefaultHandlerFunc returns the 0 handler function for msgType.
func GetDefaultHandlerFunc() HandlerFunc {
	entry, ok := messageRegistry[-1]
	if !ok {
		return nil
	}
	return entry.handler
}

func GetServiceHandler(invoke string) (*handlerUnmarshaler2, bool) {
	entry, ok := messageRegistry2[invoke]
	return entry, ok
}

// Message represents the structured data that can be handled.
type Message interface {
	MessageNumber() int32
	Serialize() ([]byte, error)
}

type XMessage struct {
	Id     int64  `json: "id"`    //代理服务生成的连接ID
	Client string `json:"client"` //用户端唯一标识
	Invoke string `json:"invoke"` //接口hash值
	Data   []byte `json:"data"`   //消息体
}

type DMessage struct {
	Content []byte
}

// HeartBeatMessage for application-level keeping alive.
type HeartBeatMessage struct {
	Timestamp int64
}

// Serialize serializes HeartBeatMessage into bytes.
func (hbm HeartBeatMessage) Serialize() ([]byte, error) {
	buf.Reset()
	err := binary.Write(buf, binary.LittleEndian, hbm.Timestamp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MessageNumber returns message number.
func (hbm HeartBeatMessage) MessageNumber() int32 {
	return HeartBeat
}

// DeserializeHeartBeat deserializes bytes into Message.
func DeserializeHeartBeat(data []byte) (message Message, err error) {
	var timestamp int64
	if data == nil {
		return nil, ErrNilData
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &timestamp)
	if err != nil {
		return nil, err
	}
	return HeartBeatMessage{
		Timestamp: timestamp,
	}, nil
}

// HandleHeartBeat updates connection heart beat timestamp.
func HandleHeartBeat(ctx context.Context, c WriteCloser) {
	msg := MessageFromContext(ctx)
	switch c := c.(type) {
	case *ServerConn:
		c.SetHeartBeat(msg.(HeartBeatMessage).Timestamp)
	case *ClientConn:
		c.SetHeartBeat(msg.(HeartBeatMessage).Timestamp)
	}
}

// Codec is the interface for message coder and decoder.
// Application programmer can define a custom codec themselves.
type Codec interface {
	Decode(net.Conn) (Message, error)
	Encode(Message) ([]byte, error)
}

// TypeLengthValueCodec defines a special codec.
// Format: type-length-value |4 bytes|4 bytes|n bytes <= 8M|
type TypeLengthValueCodec struct{}

// Decode decodes the bytes data into Message
func (codec TypeLengthValueCodec) Decode(raw net.Conn) (Message, error) {
	byteChan := make(chan []byte)
	errorChan := make(chan error)

	go func(bc chan []byte, ec chan error) {
		typeData := make([]byte, MessageTypeBytes)
		_, err := io.ReadFull(raw, typeData)
		if err != nil {
			ec <- err
			close(bc)
			close(ec)
			holmes.Debugln("go-routine read message type exited")
			return
		}
		bc <- typeData
	}(byteChan, errorChan)

	var typeBytes []byte

	select {
	case err := <-errorChan:
		return nil, err

	case typeBytes = <-byteChan:
		if typeBytes == nil {
			holmes.Warnln("read type bytes nil")
			return nil, ErrBadData
		}
		fmt.Println("====>>001")
		typeBuf := bytes.NewReader(typeBytes)
		var msgType int32
		if err := binary.Read(typeBuf, binary.LittleEndian, &msgType); err != nil {
			return nil, err
		}
		fmt.Println("====>>002")

		lengthBytes := make([]byte, MessageLenBytes)
		_, err := io.ReadFull(raw, lengthBytes)
		if err != nil {
			return nil, err
		}
		fmt.Println("====>>003")

		lengthBuf := bytes.NewReader(lengthBytes)
		var msgLen uint32
		if err = binary.Read(lengthBuf, binary.LittleEndian, &msgLen); err != nil {
			return nil, err
		}
		fmt.Println("====>>004")

		if msgLen > MessageMaxBytes {
			holmes.Errorf("message(type %d) has bytes(%d) beyond max %d\n", msgType, msgLen, MessageMaxBytes)
			return nil, ErrBadData
		}
		fmt.Println("====>>005")

		// read application data
		msgBytes := make([]byte, msgLen)
		_, err = io.ReadFull(raw, msgBytes)
		if err != nil {
			return nil, err
		}
		fmt.Println("====>>006")

		// deserialize message from bytes
		unmarshaler := GetUnmarshalFunc(msgType) //TODO msgType==0
		// unmarshaler := GetDefaultUnmarshalFunc()
		fmt.Printf("====>>007: %+v\n", unmarshaler == nil)

		if unmarshaler == nil {
			return nil, ErrUndefined(msgType)
		}
		return unmarshaler(msgBytes)
	}
}

// Encode encodes the message into bytes data.
func (codec TypeLengthValueCodec) Encode(msg Message) ([]byte, error) {
	data, err := msg.Serialize()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg.MessageNumber())
	binary.Write(buf, binary.LittleEndian, int32(len(data)))
	buf.Write(data)
	packet := buf.Bytes()
	return packet, nil
}

// ContextKey is the key type for putting context-related data.
type contextKey string

// Context keys for messge, server and net ID.
const (
	messageCtx contextKey = "message"
	serverCtx  contextKey = "server"
	netIDCtx   contextKey = "netid"
)

// NewContextWithMessage returns a new Context that carries message.
func NewContextWithMessage(ctx context.Context, msg Message) context.Context {
	return context.WithValue(ctx, messageCtx, msg)
}

// MessageFromContext extracts a message from a Context.
func MessageFromContext(ctx context.Context) Message {
	return ctx.Value(messageCtx).(Message)
}

// NewContextWithNetID returns a new Context that carries net ID.
func NewContextWithNetID(ctx context.Context, netID int64) context.Context {
	return context.WithValue(ctx, netIDCtx, netID)
}

// NetIDFromContext returns a net ID from a Context.
func NetIDFromContext(ctx context.Context) int64 {
	return ctx.Value(netIDCtx).(int64)
}

// MessageNumber returns the message number.
func (dm DMessage) MessageNumber() int32 {
	return -1
}

// Serialize Serializes Message into bytes.
func (dm DMessage) Serialize() ([]byte, error) {
	return []byte(dm.Content), nil
}

func (dm DMessage) Data() []byte {
	return dm.Content
}

func DeserializeMessage(data []byte) (message Message, err error) {
	if data == nil {
		return nil, ErrNilData
	}
	msg := DMessage{
		Content: data,
	}
	return msg, nil
}

func _deserializeMessage(data []byte) (message Message, err error) {
	if data == nil {
		return nil, ErrNilData
	}
	// content := string(data)
	msg := DMessage{
		Content: data,
	}
	return msg, nil
}

func _processMessage(ctx context.Context, conn WriteCloser) {
	switch conn.(type) {
	case *ServerConn:
		log.Infof("_process message start server: %+v", "server")
	case *ClientConn:
		log.Infof("_process message start client: %+v", "client")
	}
	_, ok := ServerFromContext(ctx)
	if ok {
		msg := MessageFromContext(ctx)
		// s.Broadcast(msg)
		data, _ := msg.Serialize()
		xm := &XMessage{}
		json.Unmarshal(data, xm)
		holmes.Infof("ProcessMessage: %+v|%+v", xm, string(xm.Data))
		val, ok := messageRegistry2[xm.Invoke]
		if ok {
			log.Infof("find message registry 2")
			DealServiceMessage(val, xm)
		}
	}
}

func GetServiceName(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return path.Base(t.PkgPath())
}

func _invokeFunc(method handlerUnmarshaler2, xm *XMessage) {
	req := reflect.New(method.ParamType.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(xm.Data, req); err != nil {
		log.Errorf("proto unmarshal failed: %+v", err)
		return
	}

	ctx := context.Background()
	results := method.Method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	log.Infof("method call result: %+v", len(results))
}

func DealServiceMessage(method *handlerUnmarshaler2, xm *XMessage) ([]byte, error) {
	req := reflect.New(method.ParamType.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(xm.Data, req); err != nil {
		log.Errorf("proto unmarshal failed: %+v", err)
		return nil, err
	}

	ctx := context.Background()
	results := method.Method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	serviceResp := results[0].Interface().(proto.Message)
	log.Infof("method call result: %+v|%+v", len(results), serviceResp)
	return proto.Marshal(serviceResp)
}

func _dealClientMessage() {

}
