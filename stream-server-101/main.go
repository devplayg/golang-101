package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/hybridgroup/mjpeg"
	"github.com/pkg/errors"
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
		cmdFlags   = flag.NewFlagSet("", flag.ExitOnError)
		recvHost   = cmdFlags.String("h", "127.0.0.1", "Receieve host")
		recvPort   = cmdFlags.String("rport", "8000", "receieve port")
		streamHost = cmdFlags.String("shost", "127.0.0.1", "Stream host")
		streamPort = cmdFlags.String("spport", "8080", "Stream port")
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

	// Open streaming
	stream := mjpeg.NewStream()
	http.Handle("/", stream)
	go func() {
		err := http.ListenAndServe(*streamHost+":"+*streamPort, nil)
		if err != nil {
			log.Println(errors.Wrap(err, "failed to open stream port"))
		}
	}()
	//log.Println("===============================")
	log.Printf("streamer started listening on %s:%s", *streamHost, *streamPort)

	// Accept connections
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println(errors.Wrap(err, "failed to accept"))
			continue
		}
		go handleConnection(conn, stream)
	}
}

func handleConnection(conn net.Conn, stream *mjpeg.Stream) {
	log.Printf("new connection %v", conn.RemoteAddr().String())
	buf := make([]byte, 1024*1000)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println(errors.Wrapf(err, "closed from client; %v", conn.RemoteAddr().String()))
				return
			}
			log.Println(errors.Wrapf(err, "failed to receive data; err: %v", err))
			return
		}

		data := buf[:n]
		decoder := gob.NewDecoder(bytes.NewReader(data))
		var m Message
		err = decoder.Decode(&m)
		if err != nil {
			log.Println(errors.Wrap(err, "failed to decode"))
			continue
		}

		img, err := gocv.NewMatFromBytes(m.Rows,m.Cols,m.MatType,m.Data)
		if err != nil {
			log.Println(errors.Wrap(err, "failed to back to mat"))
			continue
		}
		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)

		log.Printf("[%3d] len=%-4d, type=%d", m.Seq, n, m.MatType)
	}
}
