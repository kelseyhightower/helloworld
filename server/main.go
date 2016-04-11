package main

import (
	"log"
	"net"

	pb "github.com/kelseyhightower/helloworld/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	return &pb.Response{Message: "Hello " + in.Name}, nil
}

func main() {
	ln, err := net.Listen("tcp", ":7800")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(ln)
}
