package main

import (
	"flag"
	"io"
	"log"

	pb "github.com/kelseyhightower/helloworld/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var (
		server = flag.String("server", "127.0.0.1:10000", "Server address.")
		name   = flag.String("name", "", "Username to use.")
	)
	flag.Parse()

	conn, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	stream, err := c.SayHelloStream(context.Background(), &pb.Request{Name: *name})
	if err != nil {
		log.Fatal(err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(response.Message)
	}
}
