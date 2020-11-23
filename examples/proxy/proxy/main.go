package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"git.ezbuy.me/ezbuy/evarmor/common/log"

	"git.ezbuy.me/ezbuy/evarmor/examples/proxy/proxy/service"
)

func main() {
	// 启动服务
	// 监听注册服务
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 11110))
	if err != nil {
		log.Errorf("proxy listen failed: %q", err)
		return
	}

	proxyServer := service.NewProxyServer()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		proxyServer.Stop()
	}()

	proxyServer.Start(l)
}
