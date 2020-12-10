package chat

import (
	"context"
	"encoding/json"

	tao "git.ezbuy.me/ezbuy/evarmor/common/base"
	"git.ezbuy.me/ezbuy/evarmor/common/log"
	pchat "git.ezbuy.me/ezbuy/evarmor/rpc/proto/chat"
	"github.com/golang/protobuf/proto"
	"github.com/leesper/holmes"
)

const (
	// ChatMessage is the message number of chat message.
	ChatMessage int32 = 0
)

// Message defines the chat message.
type Message struct {
	Content []byte
}

// MessageNumber returns the message number.
func (cm Message) MessageNumber() int32 {
	return ChatMessage
}

// Serialize Serializes Message into bytes.
func (cm Message) Serialize() ([]byte, error) {
	return cm.Content, nil
}

// DeserializeMessage deserializes bytes into Message.
func DeserializeMessage(data []byte) (message tao.Message, err error) {
	if data == nil {
		return nil, tao.ErrNilData
	}
	// content := string(data)
	msg := Message{
		Content: data,
	}
	return msg, nil
}

// ProcessMessage handles the Message logic.
func ProcessMessage(ctx context.Context, conn tao.WriteCloser) {
	holmes.Infof("===>>> ProcessMessage")
	msg := tao.MessageFromContext(ctx)
	data, _ := msg.Serialize()
	xm := &tao.XMessage{}
	json.Unmarshal(data, xm)
	log.Infof("ProcessMessage client: %+v|%+v", xm, xm.Data)
	// var resp *pchat.SayHelloResp
	in := &pchat.SayHelloResp{}
	if err := proto.Unmarshal(xm.Data, in); err != nil {
		log.Errorf("process message unmarshal failed: %q", err)
		return
	}
	log.Infof("process message result: %+v", in.GetResponse())
	// s, ok := tao.ServerFromContext(ctx)
	// if ok {
	// 	msg := tao.MessageFromContext(ctx)
	// 	s.Broadcast(msg)
	// }
}
