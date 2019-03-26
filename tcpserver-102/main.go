package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
)

type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
}

func main() {
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host     = cmdFlags.String("h", "127.0.0.1", "Host")
		port     = cmdFlags.String("p", "8000", "Port")
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
	log.Printf("server started listening on %s:%s", *host, *port)

	// Accept connections
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println(errors.Wrap(err, "failed to accept"))
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("new connection %v", conn.RemoteAddr().String())
	buf := make([]byte, 100)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println(errors.Wrapf(err, "closed from client; %v", conn.RemoteAddr().String()))
				continue
			}
			log.Println(errors.Wrapf(err, "fail to receive data; err: %v", err))
			return
		}

		data := buf[:n]
		decoder := gob.NewDecoder(bytes.NewReader(data))
		var m Message
		err = decoder.Decode(&m)
		if err != nil {
			log.Println(errors.Wrapf(err, "failed to decode", err))
			continue
		}

		log.Printf("[%3d] len=%-4d data=%-6s, %v", m.Seq, n, m.Data, m.Data)
	}
}
