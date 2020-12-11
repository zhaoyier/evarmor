package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	tao "git.ezbuy.me/ezbuy/evarmor/common/base"
	"git.ezbuy.me/ezbuy/evarmor/examples/chat"
	pchat "git.ezbuy.me/ezbuy/evarmor/rpc/proto/chat"
	"github.com/golang/protobuf/proto"
	"github.com/leesper/holmes"
)

func main() {
	defer holmes.Start().Stop()

	tao.Register(tao.ProxyMessageType, tao.DeserializeMessage, chat.ProcessMessage)

	c, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		holmes.Fatalln(err)
	}

	onConnect := tao.OnConnectOption(func(c tao.WriteCloser) bool {
		holmes.Infoln("on connect")
		return true
	})

	onError := tao.OnErrorOption(func(c tao.WriteCloser) {
		holmes.Infoln("on error")
	})

	onClose := tao.OnCloseOption(func(c tao.WriteCloser) {
		holmes.Infoln("on close")
	})

	onMessage := tao.OnMessageOption(func(msg tao.Message, c tao.WriteCloser) {
		// fmt.Print(msg.(chat.Message).Content)
	})

	options := []tao.ServerOption{
		onConnect,
		onError,
		onClose,
		onMessage,
		tao.ReconnectOption(),
	}

	conn := tao.NewClientConn(0, c, options...)
	defer conn.Close()

	conn.Start()
	for {
		reader := bufio.NewReader(os.Stdin)
		talk, _ := reader.ReadString('\n')
		if talk == "bye\n" {
			break
		} else {
			data, _ := proto.Marshal(&pchat.SayHelloReq{
				Request: strings.TrimSpace(talk),
			})

			xm := &tao.XMessage{
				Id:     100,
				Client: "abc",
				Invoke: "SayHello",
				Data:   data,
			}
			// rawBytes, _ := json.Marshal(xm)
			// msg := tao.DMessage{
			// 	Content: rawBytes,
			// }

			if err := conn.Write(xm); err != nil {
				holmes.Infoln("error", err)
			}
		}
	}
	fmt.Println("goodbye")
}
