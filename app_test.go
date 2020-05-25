package main

import (
	"context"
	"google.golang.org/grpc"
	"micro/service"
	"testing"
)

func TestGrpc(t *testing.T) {
	dial, err := grpc.Dial("localhost:7000", grpc.WithInsecure())
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
