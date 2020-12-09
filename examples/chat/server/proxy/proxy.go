package proxy

import (
	"context"
	"encoding/json"

	tao "git.ezbuy.me/ezbuy/evarmor/common/base"
	"git.ezbuy.me/ezbuy/evarmor/common/log"
	"github.com/leesper/holmes"
)

// ProcessMessage handles the Message logic.
func ProcessMessage(ctx context.Context, conn tao.WriteCloser) {
	switch conn.(type) {
	case *tao.ServerConn:
		log.Infof("_process message start proxy server: %+v", "server")
	case *tao.ClientConn:
		log.Infof("_process message start proxy client: %+v", "client")
	}
	_, ok := tao.ServerFromContext(ctx)
	if ok {
		msg := tao.MessageFromContext(ctx)
		// s.Broadcast(msg)
		data, _ := msg.Serialize()
		xm := &tao.XMessage{}
		json.Unmarshal(data, xm)
		holmes.Infof("ProcessMessage: %+v|%+v", xm, string(xm.Data))
		val, ok := tao.GetServiceHandler(xm.Invoke)
		// val, ok := messageRegistry2[xm.Invoke]
		if ok {
			log.Infof("find message registry 2")
			tao.DealServiceMessage(val, xm)
		}
	}
}
