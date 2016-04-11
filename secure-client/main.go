package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/kelseyhightower/helloworld/helloworld"

	"github.com/kelseyhightower/helloworld/credentials/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var (
		caCert = flag.String("ca", "/etc/helloworld/ca.pem", "Trusted CA certificate.")
		server = flag.String("server", "127.0.0.1:10000", "Server address.")
		name   = flag.String("name", "", "Username to use.")
		token  = flag.String("token", "/etc/helloworld/token", "JWT token.")
	)
	flag.Parse()

	rawCACert, err := ioutil.ReadFile(*caCert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rawCACert)

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: caCertPool,
	})

	jwtCreds, err := jwt.NewFromTokenFile(*token)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(*server,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(jwtCreds),
	)
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
