package client

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	pb2 "git.ezbuy.me/ezbuy/evarmor/common/network"
	pb "git.ezbuy.me/ezbuy/evarmor/rpc/evarmor"
	"github.com/golang/protobuf/proto"
)

type Msg struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}
type Resp struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
}

func Client() {
	flag.Parse()
	conn, err := net.Dial("tcp", "localhost:13101")
	if err != nil {
		fmt.Println("connect error", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connecting to localhost:13101")
	var wg sync.WaitGroup
	wg.Add(2)
	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)
	wg.Wait()
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	req := &pb.HelloRequest{
		Name: "zhao",
	}

	data, _ := proto.Marshal(req)

	xm := &pb2.XMessage{
		Code: "SayHello",
		Data: string(data),
	}
	data, _ = proto.Marshal(xm)
	fmt.Printf("====>>101:%+v|%+v\n", string(data), xm.String())
	writer := bufio.NewWriter(conn)
	writer.Write(data)
	// writer.Write([]byte("\n"))
	writer.Flush()

	// for i := 10; i > 0; i-- {
	// 	d := "hello" + strconv.Itoa(i)
	// 	msg := Msg{
	// 		Data: d,
	// 		Type: 1,
	// 	}
	// 	b, _ := json.Marshal(msg)
	// 	writer := bufio.NewWriter(conn)
	// 	writer.Write(b)
	// 	writer.Write([]byte("\n"))
	// 	writer.Flush()
	// }
	fmt.Println("write done")
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(conn)
	for i := 1; i <= 10; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		var resp Resp
		json.Unmarshal(line, &resp)
		fmt.Println("status", resp.Status, " content:", resp.Data)
	}
	fmt.Println("read done")
}
