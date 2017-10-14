package main

import (
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func validateTokenFromContext(ctx context.Context, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	var (
		token *jwt.Token
		err   error
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	token, err = jwt.Parse(jwtToken[0], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && token.Valid {
		return token, nil
	}
	return nil, err
}
