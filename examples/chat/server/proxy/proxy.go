package proxy

import (
	"context"
	"encoding/json"

	tao "git.ezbuy.me/ezbuy/evarmor/common/base"
	"git.ezbuy.me/ezbuy/evarmor/common/log"
	"github.com/leesper/holmes"
	// "github.com/golang/protobuf/proto"
)

// ProcessMessage handles the Message logic.
func ProcessMessage(ctx context.Context, conn tao.WriteCloser) {
	switch conn.(type) {
	case *tao.ServerConn:
		log.Infof("_process message start proxy server: %+v", "server")
		netId := tao.NetIDFromContext(ctx)
		s, ok := tao.ServerFromContext(ctx)
		if ok {
			msg := tao.MessageFromContext(ctx)
			data, _ := msg.Serialize()
			xm := &tao.XMessage{}
			json.Unmarshal(data, xm)
			holmes.Infof("ProcessMessage: %+v|%+v", xm, string(xm.Data))
			val, ok := tao.GetServiceHandler(xm.Invoke)
			if ok {
				log.Infof("find message registry 2")
				resp, err := tao.DealServiceMessage(val, xm)
				if err != nil {
					log.Errorf("deal server message failed: %+v", resp)
				}
				log.Infof("process message response: %+v", resp)
				rd, _ := json.Marshal(&tao.XMessage{
					Invoke: "SayHello",
					Data:   resp,
				})
				s.Unicast(netId, tao.DMessage{Content: rd})
			}
		}
	case *tao.ClientConn:
		log.Infof("_process message start proxy client: %+v", "client")
		msg := tao.MessageFromContext(ctx)
		data, _ := msg.Serialize()
		xm := &tao.XMessage{}
		json.Unmarshal(data, xm)
		holmes.Infof("ProcessMessage client: %+v|%+v", xm, string(xm.Data))
		// proto.Unmarshal()
	}

}
