package main

import (
	"log"
	"net"
	"time"

	pb "github.com/kelseyhightower/helloworld/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	return &pb.Response{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloStream(in *pb.Request, stream pb.Greeter_SayHelloStreamServer) error {
	for {
		err := stream.Send(&pb.Response{Message: "Hello " + in.Name})
		if err != nil {
			log.Println(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(ln)
}
