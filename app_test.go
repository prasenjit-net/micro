package main

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"micro/service"
	"testing"
)

func TestGrpc(t *testing.T) {
	go main()
	dial, err := grpc.Dial("localhost:7000",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		t.Error(err)
	}
	client := service.NewPingPongClient(dial)

	pong, err := client.Ping(context.Background(), &service.Ping{Subject: "World"})
	if err != nil {
		t.Error(err)
	}
	println(pong.Reply)
}
