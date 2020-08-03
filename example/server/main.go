package server

import (
	"context"
	"log"

	"git.ezbuy.me/ezbuy/evarmor/common/network"
	pb "git.ezbuy.me/ezbuy/evarmor/rpc/evarmor"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("request: ", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func Server() {
	s := network.NewServer("desc")

	s.RegisterServer(&server{})

	// xm := &evarmor.XMessage{}
	// if err := proto.Unmarshal([]byte("hello"), xm); err != nil {
	// 	fmt.Printf("====>>0022:%q\n", err)
	// }
	s.Start("localhost:13101")
}
