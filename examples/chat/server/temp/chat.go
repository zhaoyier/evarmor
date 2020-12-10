package temp

import (
	// "git.ezbuy.me/ezbuy/base/misc/context"
	"context"
	"fmt"

	"git.ezbuy.me/ezbuy/evarmor/common/log"
	pchat "git.ezbuy.me/ezbuy/evarmor/rpc/proto/chat"
	// "golang.org/x/net/context"
)

type Chat struct {
}

func (c *Chat) SayHello(ctx context.Context, in *pchat.SayHelloReq) (*pchat.SayHelloResp, error) {
	log.Infof("SayHello I am here: %+v", in.GetRequest())
	return &pchat.SayHelloResp{
		Response: fmt.Sprintf("收到: %+v", in.GetRequest()),
	}, nil
}
