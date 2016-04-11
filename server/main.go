package main

import (
	"log"
	"net"

	pb "github.com/kelseyhightower/helloworld/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return &pb.Response{Message: "Hello " + request.Name}, nil
}

func main() {
	log.Println("Helloworld service starting...")
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
    log.Println("Helloworld service started successfully.")
	s.Serve(ln)
}
