package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

var (
	addr    = flag.String("listen", ":8080", "websocket address to listen on")
	id      = flag.String("id", "vipstats", "id of the ethstats server")
	autotls = flag.Bool("autotls", true, "setup TLS on port :443 when listen is on port :80")
)

func main() {
	ethstats := &Server{
		Name: *id,
	}

	_, port, err := net.SplitHostPort(*addr)
	if err != nil {
		exit(1, "failed to parse address", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", ethstats)

	if port == "80" && *autotls {
		log.Print("starting autocert process")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("certs"),
		}

		https := &http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		log.Fatal(https.ListenAndServeTLS("", ""))
	} else {
		log.Printf("listening on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, mux))
	}

}

// exit prints an error and exits with the given code
func exit(code int, msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(code)
}
