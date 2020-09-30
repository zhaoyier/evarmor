package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"git.ezbuy.me/ezbuy/evarmor/examples/proxy/service"
)

func main() {
	// 启动服务
	// 监听注册服务
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 12345))
	if err != nil {
		// holmes.Fatalln("listen error", err)
	}

	proxyServer := service.NewProxyServer()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		proxyServer.Stop()
	}()
	// service.Beacon([]string{})

	proxyServer.Start(l)
}
