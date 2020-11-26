package temp

import (
	// "git.ezbuy.me/ezbuy/base/misc/context"
	"context"

	"git.ezbuy.me/ezbuy/evarmor/common/log"
	pchat "git.ezbuy.me/ezbuy/evarmor/rpc/proto/chat"
	// "golang.org/x/net/context"
)

type Chat struct {
}

func (c *Chat) SayHello(ctx context.Context, in *pchat.SayHelloReq) (*pchat.SayHelloResp, error) {
	log.Infof("SayHello I am here")
	return &pchat.SayHelloResp{
		Response: "收到",
	}, nil
}
