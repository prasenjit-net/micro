package service

import "context"

type PingServerImpl struct {
	UnimplementedPingPongServer
}

func (ps *PingServerImpl) Ping(_ context.Context, ping *Ping) (*Pong, error) {
	return &Pong{Reply: ping.Subject}, nil
}
