package service

import (
	"git.ezbuy.me/ezbuy/evarmor/common/base"
)

// ProxyServer is the proxy server.
type ProxyServer struct {
	*base.Server
}

// NewProxyServer returns a ProxyServer.
func NewProxyServer() *ProxyServer {
	onConnectOption := base.OnConnectOption(func(conn base.WriteCloser) bool {
		// holmes.Infoln("on connect")
		return true
	})
	onErrorOption := base.OnErrorOption(func(conn base.WriteCloser) {
		// holmes.Infoln("on error")
	})
	onCloseOption := base.OnCloseOption(func(conn base.WriteCloser) {
		// holmes.Infoln("close chat client")
	})
	return &ProxyServer{
		base.NewServer(onConnectOption, onErrorOption, onCloseOption),
	}
}
