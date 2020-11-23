package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"git.ezbuy.me/ezbuy/evarmor/common/base"
	"git.ezbuy.me/ezbuy/evarmor/common/log"
	mproto "git.ezbuy.me/ezbuy/evarmor/common/proto"
	"git.ezbuy.me/ezbuy/evarmor/common/utils"
)

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:11110")
	if err != nil {
		log.Errorf("net dial failed: %q", err)
		return
	}

	onMessage := base.OnMessageOption(func(msg *mproto.XMessage, c base.WriteCloser) {
		fmt.Printf("client on message: %+v\n", msg)
	})

	options := []base.ServerOption{onMessage}
	conn := base.NewClientConn(0, c, options...)

	defer conn.Close()

	conn.Start()
	for {
		reader := bufio.NewReader(os.Stdin)
		talk, _ := reader.ReadString('\n')
		if talk == "bye\n" {
			break
		} else {
			talk = strings.TrimSpace(talk)
			fmt.Println("client write: ", talk)
			if err := conn.Write(&mproto.XMessage{
				Code: utils.CRC32("SayHello"),
				Data: []byte(talk),
			}); err != nil {
				log.Errorf("client write failed: %q", err)
			} else {
				log.Infof("client write succedd.")
			}
		}
	}
}
