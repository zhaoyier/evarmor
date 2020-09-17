package gateway

import (
	"git.ezbuy.me/ezbuy/evarmor/common/network"
)

type server struct{}

func StartProxy() {
	s := network.NewServer("proxy")

	// xm := &evarmor.XMessage{}
	// if err := proto.Unmarshal([]byte("hello"), xm); err != nil {
	// 	fmt.Printf("====>>0022:%q\n", err)
	// }
	s.Start("localhost:13100")
}
