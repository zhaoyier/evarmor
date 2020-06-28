package main

import (
	"echoserver/message"
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
)

//错误检查
func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

//服务器端接收数据线程
//      客户端   conn
//      消息队列 messages
func Handler(conn net.Conn, messages chan string) {
	buf := make([]byte, 1024)
	for {
		// 读取数据
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		reciveStr := string(buf[0:lenght])
		// 数据转发给所有客户端
		messages <- reciveStr
	}
}

//服务器发送数据的线程
//      客户端集合 conns
//      消息队列   messages
func echoHandler(conns *map[string]net.Conn, messages chan string) {
	// 遍历数据
	for {
		msg := <-messages
		// 遍历客户端
		for key, value := range *conns {
			_, err := value.Write([]byte(msg))
			if err != nil {
				fmt.Println(err.Error())
				delete(*conns, key)
			}
		}
	}
}

//启动服务器
//      端口 port
func StartServer(port string) {
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	// 客户端集合
	conns := make(map[string]net.Conn)
	// 发送消息队列
	messages := make(chan string, 10)
	// 启动发送消息线程
	go echoHandler(&conns, messages)

	for {
		// 新客户端
		conn, err := l.Accept()
		checkError(err, "Accept")
		conns[conn.RemoteAddr().String()] = conn
		// 启动接收消息线程
		go Handler(conn, messages)
	}
}

//客户端发送线程
//      发送连接 conn
func chatSend(conn net.Conn) {
	var input string
	username := conn.LocalAddr().String()
	for {
		// 获取输入
		fmt.Scanln(&input)
		if input == "/quit" {
			fmt.Println("ByeBye..")
			conn.Close()
			os.Exit(0)
		}

		msg := &message.Chat{
			User:    proto.String(username),
			Message: proto.String(input),
		}

		// 进行编码
		data, err := proto.Marshal(msg)
		if checkError(err, "ProtoMarshal") == false {
			conn.Close()
			os.Exit(0)
		}

		// 发送给服务端
		lens, err := conn.Write(data)
		if lens > 0 {
		}
		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			break
		}
	}
}

//客户端启动函数
//      远程ip地址和端口 tcpaddr
func StartClient(tcpaddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	checkError(err, "ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")
	// 启动客户端发送线程
	go chatSend(conn)

	// 接收数据
	buf := make([]byte, 1024)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			os.Exit(0)
		}
		msg := &message.Chat{}
		err = proto.Unmarshal(buf[0:lenght], msg)
		if checkError(err, "ProtoUnmarshal") == false {
			conn.Close()
			os.Exit(0)
		}
		fmt.Println(msg.GetUser() + "    " + msg.GetMessage())
	}
}

//启动服务器端：  file server [port]             eg: file server 9090
//启动客户端：    file client [Server Ip Addr]:[Server Port]    eg: file client 192.168.0.74:9090
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong pare")
		os.Exit(0)
	}

	if os.Args[1] == "server" && len(os.Args) == 3 {
		StartServer(os.Args[2])
	}

	if os.Args[1] == "client" && len(os.Args) == 3 {
		StartClient(os.Args[2])
	}
}