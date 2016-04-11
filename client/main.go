package main

import (
	"flag"
	"fmt"
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

	response, err := c.SayHello(context.Background(), &pb.Request{Name: *name})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message)
}
