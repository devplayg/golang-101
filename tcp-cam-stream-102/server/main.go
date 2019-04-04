package main

import (
	"flag"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"log"
	"net"
	"net/http"
	"os"
	"encoding/gob"
	"io"
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
		recvHost   = cmdFlags.String("recvhost", "127.0.0.1", "Receieve host")
		recvPort   = cmdFlags.String("recvport", "8000", "receieve port")
		streamHost = cmdFlags.String("streamhost", "127.0.0.1", "Stream host")
		streamPort = cmdFlags.String("streampport", "8080", "Stream port")
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
			log.Println("failed to open stream port;", err)
		}
	}()
	log.Printf("streamer started listening on %s:%s", *streamHost, *streamPort)

	// Accept connections
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println("failed to accept;", err)
			continue
		}
		go handleConnection(conn, stream)
	}
}

func handleConnection(conn net.Conn, stream *mjpeg.Stream) {
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
		//spew.Dump(m)

		//n, err := conn.Read(buf)
		//if err != nil {
		//	if err == io.EOF {
		//		log.Printf("closed from client; %v", conn.RemoteAddr().String())
		//		return
		//	}
		//	log.Println("failed to receive data;", err)
		//	return
		//}
		//
		//log.Printf("length=%d", n)

		//// Decode
		//decoder := gob.NewDecoder(bytes.NewReader(buf[:n]))
		//var m Message
		//err = decoder.Decode(&m)
		//if err != nil {
		//	log.Println("failed to decode;", err)
		//	continue
		//}
		//log.Printf("[%3d] len=%-4d, type=%d", m.Seq, n, m.MatType)
		//
		//img, err := gocv.NewMatFromBytes(m.Rows, m.Cols, m.MatType, m.Data)
		//if err != nil {
		//	log.Println("failed to back to mat;", err)
		//	continue
		//}
		//buf, _ := gocv.IMEncode(".jpg", img)
		//stream.UpdateJPEG(buf)

	}
}
