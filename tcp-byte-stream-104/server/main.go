package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/devplayg/golang-101/tcp-byte-stream-104/obj"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"io"
	"time"
)

const (
	ConnPoolSize = 10
)

var (
	// Command flags
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)

	// Connection
	connMap  sync.Map
	connPool = make(chan net.Conn, ConnPoolSize)
	listener net.Listener

	// Etc
	quit  = make(chan bool)
)

func init() {
	// Logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Handle command flags
	var (
		host  = cmdFlags.String("host", "127.0.0.1", "Host")
		port  = cmdFlags.Int("port", 8000, "Port")
		debug = cmdFlags.Bool("debug", false, "debug")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Check debug mode
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("started debugging")
	}

	var err error

	// Start listener
	listener, err = net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Error("failed to start server;", err)
		return
	}
	defer listener.Close()
	log.Debug("started server")

	// Start server
	go run()

	// Wait for stop signal
	waitForSignals()
	log.Info("stopped server")
}

func run() error {
	defer func() {
		log.Info("stopping server")
		close(quit)
	}()

	// Start listening
	var connID int64 = 1
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-quit: //
				return nil
			default:
			}
		}
		log.Debugf("connection #%d is connected", connID)

		connMap.Store(connID, conn)
		go handleConnection(conn, connID)
		connID++
	}

	return nil
}

func handleConnection(conn net.Conn, connID int64) error {
	var success_count int64
	t := time.Now()

	defer func() {
		log.Debugf("closing connection #%d", connID)
		log.Infof("time: %3.1f, count: %d, average=%3.1f", time.Since(t).Seconds(), success_count, float64(success_count) / time.Since(t).Seconds())
		conn.Close()
		connMap.Delete(connID)
	}()
	log.Debugf("connection #%d is ready", connID)

	//buf := make([]byte, 100)
	//data := make([]byte, 0)
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	var err error
	for {
		// Receive message header
		msgHeader := obj.MessageHeader{}
		log.Debug("ready to receive message header")
		err = decoder.Decode(&msgHeader)
		if err != nil {
			if err == io.EOF {
				log.Debug("disconnected from connection;", conn.RemoteAddr().String())
				return nil
			}
			log.Error("failed to receive header and decode;", err)
			return err
		}
		log.Debugf("received message header. payload is %d bytes", msgHeader.PayloadSize)

		// Response
		err = encoder.Encode(&obj.Response{Code:1})
		if err != nil {
			log.Error("failed to response;", err)
			continue
		}
		//log.Debug("sent response ")

		// Receive payload
		var received uint32
		buf := make([]byte, 100)
		data := make([]byte, 0)
		for received < msgHeader.PayloadSize {
			n, err := conn.Read(buf)
			if err != nil {
				log.Error("failed to receive payload;", err)
				break
			}
			received += uint32(n)
			data = append(data, buf[:n]...)
			//spew.Dump(data)
		}
		log.Debugf("received payload successfully. payload=%d, received=%d", msgHeader.PayloadSize, received)

		// Verify
		code := 0
		//spew.Dump(data)
		decoder := gob.NewDecoder(bytes.NewReader(data))
		var m obj.Message
		err = decoder.Decode(&m)
		if err != nil {
			log.Error("failed to decode payload;", err)
			code = -1
		} else {
			if m.Verify() {
				log.Debug("checksum ok")
				code = 200
				success_count++
			} else {
				log.Error("checksum error")
				code = -2
			}
		}

		// Response
		err = encoder.Encode(&obj.Response{Code:code})
		if err != nil {
			log.Error("failed to response;", err)
			continue
		}
	}
	//log.Debug("data: %v", data)



	//var err error
	//dataBuffer := make([]byte, 4)
	//for {
	//	log.Debug("ready to receive msg")
	//	msgLen, err := conn.Read(dataBuffer)
	//	if err != nil {
	//		if err == io.EOF {
	//			log.Debugf("disconnected connection #%d", connID)
	//			return err
	//		}
	//		log.Error("failed to received;", err)
	//		continue
	//	}
	//
	//	// Response
	//	_, err = conn.Write([]byte{0x01})
	//	if err != nil {
	//		if err == io.EOF {
	//			log.Debugf("disconnected connection #%d", connID)
	//			return err
	//		}
	//		log.Error(err)
	//	}
	//	log.Debugf("got_message=%d", msgLen)
	//}

	return nil
}

func waitForSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		log.Info("Signal received, shutting down...")
	}
}

func Deserialize(b []byte) (*obj.Message, error) {
	var m obj.Message

	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
