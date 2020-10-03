package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	// Flag
	port := flag.String("p", "1553", "Port")
	flag.Parse()

	// API list
	apiList := []string{
            "/api1",
            "/api2",
            "/api3",    
	}
	startHttpServer(*port, apiList)
}

// handles common request
func commonHandler(w http.ResponseWriter, r *http.Request) {
	debugHttpRequest(r)
	w.WriteHeader(http.StatusOK)
}

// Debug HTTP request
func debugHttpRequest(r *http.Request) {
	fmt.Printf("%s\n", strings.Repeat("=", 100))
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}

// starts http server
func startHttpServer(port string, apiList []string) {
	srv := createHttpServer(port, apiList)
	ch := make(chan struct{})
	go func() {
		defer close(ch)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-ch
}

// creates HTTP server
func createHttpServer(port string, apiList []string) *http.Server {
	//r := mux.NewRouter()
	r := http.NewServeMux()
	for _, api := range apiList {
		r.HandleFunc(api, commonHandler)
		fmt.Printf("[api] %s\n", api)
	}
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("server listening on %s\n", srv.Addr)
	return srv

}
