package temp

import (
	"git.ezbuy.me/ezbuy/base/misc/context"
	pchat "git.ezbuy.me/ezbuy/evarmor/rpc/proto/chat"
	// "golang.org/x/net/context"
)

type Chat struct {
}

func (c *Chat) SayHello(ctx context.T, in *pchat.SayHelloReq) (*pchat.SayHelloResp, error) {
	return &pchat.SayHelloResp{
		Response: "收到",
	}, nil
}
