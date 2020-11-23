package service

import (
	"git.ezbuy.me/ezbuy/evarmor/common/base"
	mproto "git.ezbuy.me/ezbuy/evarmor/common/proto"
	"git.ezbuy.me/ezbuy/evarmor/common/utils"

	// lproto "git.ezbuy.me/ezbuy/evarmor/common/proto"
	"git.ezbuy.me/ezbuy/evarmor/common/log"
)

// ProxyServer is the proxy server.
type ProxyServer struct {
	*base.Server
}

// NewProxyServer returns a ProxyServer.
func NewProxyServer() *ProxyServer {
	onConnectOption := base.OnConnectOption(func(conn base.WriteCloser) bool {
		// holmes.Infoln("on connect")
		log.Infof("option on connect")
		return true
	})
	onErrorOption := base.OnErrorOption(func(conn base.WriteCloser) {
		log.Infof("option on error")

		// holmes.Infoln("on error")
	})
	onCloseOption := base.OnCloseOption(func(conn base.WriteCloser) {
		// holmes.Infoln("close chat client")
		log.Infof("option on close")

	})

	onMessageOption := base.OnMessageOption(func(msg *mproto.XMessage, conn base.WriteCloser) {
		log.Infof("option on message: %+v", msg)
		// holmes.Infoln("close chat client")
		// switch msg.GetCode() {
		// case utils.CRC32(base.InternalRegisterHandler): //注册逻辑服务接口
		// 	internalRegisterHandler(msg.GetData(), conn)
		// default:
		// 	proxyMessage(msg, conn)
		// }
	})

	// etcd 通知其他服务有注册
	return &ProxyServer{
		base.NewServer(onConnectOption, onErrorOption, onCloseOption, onMessageOption),
	}
}

func internalRegisterHandler(data []byte, conn base.WriteCloser) error {
	name := string(data)
	sc := conn.(*base.ServerConn)

	// 注册代理服务接口
	sc.Belong().RegistProxy(utils.CRC32(name), sc.NetID())
	return nil
}

func proxyMessage(msg *mproto.XMessage, conn base.WriteCloser) {
	// 根据名称获取逻辑服务连接句柄
	sc := conn.(*base.ServerConn)

	if msg.GetNetid() == 0 { //来自用户端
		conn, err := sc.Belong().GetProxyConn(msg.GetCode(), sc.NetID())
		if err != nil {
			//TODO
			log.Errorf("get proxy conn failed: %q", err)
			return
		}

		msg.Netid = sc.NetID()
		if err := conn.Write(msg); err != nil {
			log.Errorf("proxy write to server failed: %q", err)
			return
		}
	}
	// 来自服务端
	if cn, ok := sc.Belong().Conn(msg.GetNetid()); ok {
		msg.Netid = 0
		if err := cn.Write(msg); err != nil {
			log.Errorf("proxy write client failed: %q", err)
			return
		}
	}
}
