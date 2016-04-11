package main

import (
	"crypto/rsa"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/kelseyhightower/helloworld/helloworld"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

type server struct {
	jwtPublicKey *rsa.PublicKey
}

func NewServer(rsaPublicKey string) (*server, error) {
	data, err := ioutil.ReadFile(rsaPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading the jwt public key: %v", err)
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the jwt public key: %s", err)
	}

	return &server{publickey}, nil
}

func (s *server) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	token, err := validateTokenFromContext(ctx, s.jwtPublicKey)
	if err != nil {
		log.Println(err)
		return &pb.Response{}, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	response := &pb.Response{
		Message: fmt.Sprintf("Hello %s (%s)", in.Name, token.Claims["email"]),
	}

	return response, nil
}

func main() {
	var (
		listenAddr   = flag.String("listen-addr", "0.0.0.0:7900", "Listen address.")
		tlsCert      = flag.String("tls-cert", "/etc/helloworld/cert.pem", "TLS server certificate.")
		tlsKey       = flag.String("tls-key", "/etc/helloworld/key.pem", "TLS server key.")
		jwtPublicKey = flag.String("jwt-public-key", "/etc/helloworld/jwt.pem", "The JWT RSA public key.")
	)
	flag.Parse()

	log.Println("Helloworld service starting...")

	cert, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
	if err != nil {
		log.Fatal(err)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	s := grpc.NewServer(grpc.Creds(creds))

	hs, err := NewServer(*jwtPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterGreeterServer(s, hs)

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Serve(ln)
}
