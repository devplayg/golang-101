package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var cmdFlags *flag.FlagSet

func init() {
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
}

func main() {
	var (
		host = cmdFlags.String("h", "127.0.0.1", "Host")
		port = cmdFlags.String("p", "8000", "Port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}

	cmdFlags.Parse(os.Args[1:])

	// Start server
	ln, err := net.Listen("tcp", *host+":"+*port)
	if nil != err {
		panic(err)
	}
	defer ln.Close()
	log.Printf("Server started listening on %s:%s", *host, *port)

	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println("failed to accept:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("new connection %v", conn.RemoteAddr().String())
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}
		if 0 < n {
			data := buf[:n]
			println(strings.TrimSpace(string(data)))
		}
	}
}
