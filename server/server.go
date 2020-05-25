package server

import (
	"crypto/tls"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"micro/service"
	"micro/util"
	"net"
	"net/http"
)

func ServeHTTP1(l net.Listener, mux *http.ServeMux) {
	s := &http.Server{
		Handler: mux,
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func ServeHTTPS(l net.Listener, mux *http.ServeMux) {
	// Load certificates.
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Panic(err)
	}

	config := getTLSConfig(certificate)

	// Create TLS listener.
	http1 := tls.NewListener(l, config)

	// Serve HTTP over TLS.
	ServeHTTP1(http1, mux)
}

func ServeGRPC(l net.Listener) {
	server := grpc.NewServer()
	service.RegisterPingPongServer(server, &service.PingServerImpl{})
	if err := server.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func TlsListener(l net.Listener) net.Listener {
	// Load certificates.
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	util.LogFatalWithMessage("Failed to load certificates", err)

	config := getTLSConfig(certificate)

	// Create TLS listener.
	return tls.NewListener(l, config)
}
