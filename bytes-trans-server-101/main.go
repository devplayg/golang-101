package main

import (
	"encoding/gob"
	"flag"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
	Rows      int
	Cols      int
	MatType   gocv.MatType
}

func main() {
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		recvHost = cmdFlags.String("recvhost", "127.0.0.1", "Receieve host")
		recvPort = cmdFlags.String("recvport", "8000", "receieve port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Start receiver
	ln, err := net.Listen("tcp", *recvHost+":"+*recvPort)
	if nil != err {
		panic(err)
	}
	defer ln.Close()
	log.Printf("receiver started listening on %s:%s", *recvHost, *recvPort)

	// Accept connections
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println("failed to accept;", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("new connection %v", conn.RemoteAddr().String())
	//buf := make([]byte, 4096)

	for {
		// Read data

		m := &Message{}
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(m)
		if err != nil {
			if err == io.EOF {
				log.Printf("closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Println("failed to decode;", err)
			continue
		}

		log.Printf("[%d]", m.Seq)

	}
}
