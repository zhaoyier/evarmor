package proxy

import (
	"context"

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
			xm := tao.MessageFromContext(ctx).(*tao.XMessage)
			log.Infof("======>>998:%+v\n", xm)
			// data, _ := msg.Serialize()
			// xm := &tao.XMessage{}
			// json.Unmarshal(data, xm)
			// holmes.Infof("ProcessMessage 02: %+v|%+v", xm, string(xm.Data))
			val, ok := tao.GetServiceHandler(string(xm.Invoke))
			if ok {
				log.Infof("find message registry 2")
				var err error
				xm.Data, err = tao.DealServiceMessage(val, xm)
				if err != nil {
					log.Errorf("deal server message failed: %+v", xm.Data)
				}
				log.Infof("process message response: %+v", xm)
				// rd, _ := json.Marshal(&tao.XMessage{
				// 	I,
				// 	Invoke: xm.Invoke,
				// 	Data:   resp,
				// })
				// xm.Data = resp
				// rd, _ := json.Marshal(xm)
				// log.Infof("=====>>775:%+v|%+v", xm, rd)

				s.Unicast(netId, xm)
			}
		}
	case *tao.ClientConn:
		log.Infof("_process message start proxy client: %+v", "client")
		msg := tao.MessageFromContext(ctx)
		xm := tao.MessageFromContext(ctx).(*tao.XMessage)
		// data, _ := msg.Serialize()
		// xm := &tao.XMessage{}
		// json.Unmarshal(data, xm)
		holmes.Infof("ProcessMessage client 002: %+v|%+v|%+v", xm, string(xm.Data), msg)
		// proto.Unmarshal()
	}

}
