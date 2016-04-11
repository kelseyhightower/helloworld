package main

import (
	"crypto/rsa"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/kelseyhightower/helloworld/helloworld"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
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

func (s *server) SayHello(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	token, err := validateTokenFromContext(ctx, s.jwtPublicKey)
	if err != nil {
		log.Println(err)
		return &pb.Response{}, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	return &pb.Response{
		Message: fmt.Sprintf("Hello %s (%s)", request.Name, token.Claims["email"]),
	}, nil
}

func (s *server) SayHelloStream(request *pb.Request, stream pb.Greeter_SayHelloStreamServer) error {
	_, err := validateTokenFromContext(stream.Context(), s.jwtPublicKey)
	if err != nil {
		log.Println(err)
		return grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	for {
		err := stream.Send(&pb.Response{Message: "Hello " + request.Name})
		if err != nil {
			log.Println(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var (
		debugAddr    = flag.String("debug-addr", "0.0.0.0:10001", "Debug listen address.")
		listenAddr   = flag.String("listen-addr", "0.0.0.0:10000", "HTTP listen address.")
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

	go s.Serve(ln)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/statusmanager", statusHandler)

	log.Println("Helloworld service started successfully.")
	log.Fatal(http.ListenAndServe(*debugAddr, nil))
}
