package main

import (
	"log"
	"micro/homepage"
	"micro/server"
	"micro/sscert"
	"micro/util"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/soheilhy/cmux"
)

// This is an example for serving HTTP and HTTPS on the same port.
func main() {
	// ensure server certificate
	hostname, err := os.Hostname()
	log.Printf("Hostname for the machine %s", hostname)
	sscert.EnsureCertificatesCreated(hostname + ",localhost")
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC)

	// create route mux and register handlers
	routeMux := http.NewServeMux()
	homepage.New().Register("/home", routeMux)

	// Create the TCP listener.
	l, err := net.Listen("tcp", ":7000")
	util.LogFatal(err)

	// Create a c-mux.
	rootMux := cmux.New(l)

	// We first match on HTTP 1.1 methods.
	http1L := rootMux.Match(cmux.HTTP1Fast())

	// If not matched, we assume that its TLS.
	//
	// Note that you can take this listener, do TLS handshake and
	// create another mux to multiplex the connections over TLS.
	tlsL := rootMux.Match(cmux.Any())
	tlsL = server.TlsListener(tlsL)

	tlsMux := cmux.New(tlsL)
	grpcL := tlsMux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpsL := tlsMux.Match(cmux.Any())

	go server.ServeHTTP1(http1L, routeMux)
	go server.ServeGRPC(grpcL)
	go server.ServeHTTP1(httpsL, routeMux)

	go func() {
		if err := tlsMux.Serve(); err != cmux.ErrListenerClosed {
			panic(err)
		}
	}()
	if err := rootMux.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
