package main

import (
	"crypto/rand"
	"encoding/gob"
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
}

const (
	dataSize = 8
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {

	// Set flags
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host     = cmdFlags.String("host", "127.0.0.1", "Host")
		port     = cmdFlags.String("port", "8000", "Port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Connect to server
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		log.Println("failed to connect to server", err)
		return
	}
	defer conn.Close()

	go send(conn)

	waitForSignals()
}

func send(conn net.Conn) {
	img := gocv.NewMat()
	defer img.Close()

	var seq int64 = 1
	//var buf bytes.Buffer
	for {
		// Create message
		m := Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      getData(dataSize),
		}
		encoder := gob.NewEncoder(conn)
		err := encoder.Encode(m)
		if err != nil {
			log.Println("failed to send;", err)
			continue
		}
		//buf.Reset()
		//encoder := gob.NewEncoder(&buf)

		//
		//// Encode
		//err := encoder.Encode(m)
		//if err != nil {
		//	log.Println("failed to encode;", err)
		//	continue
		//}
		//
		//// Send
		//n, err := conn.Write(buf.Bytes())
		//if err != nil {
		//	log.Println("failed to send data;", err)
		//	return
		//}
		//
		//log.Printf("[%3d] len=%-4d", m.Seq, n)
		log.Printf("[%3d]", m.Seq)
		seq++
		time.Sleep(5 * time.Second)
	}
}

func getData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

func waitForSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		fmt.Print("Signal received, shutting down...")
	}
}
