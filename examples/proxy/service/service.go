package service

import (
	"git.ezbuy.me/ezbuy/evarmor/common/base"
	mproto "git.ezbuy.me/ezbuy/evarmor/common/proto"
	// lproto "git.ezbuy.me/ezbuy/evarmor/common/proto"
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

	onMessageOption := base.OnMessageOption(func(msg *mproto.XMessage, conn base.WriteCloser) {
		// holmes.Infoln("close chat client")
		switch msg.GetCode() {
		case base.InternalRegisterHandler: //注册逻辑服务接口
			internalRegisterHandler(msg.GetData(), conn)
		}
	})
	return &ProxyServer{
		base.NewServer(onConnectOption, onErrorOption, onCloseOption, onMessageOption),
	}
}

func internalRegisterHandler(data []byte, conn base.WriteCloser) error {
	sc := conn.(*base.ServerConn)
	// var handler *mproto.RegisterHandler
	// if err := proto.Unmarshal(data, handler); err != nil {
	// 	log.Errorf("proto unmarshal failed: %+v", err)
	// 	return err
	// }
	// 注册代理服务接口
	sc.Belong().RegistProxy("handler.GetName()", sc.NetID())
	// sc.belong.RegistProxy(handler.GetName(), sc.NetID())
	return nil
}
