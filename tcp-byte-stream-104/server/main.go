package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/devplayg/golang-101/tcp-byte-stream-104/obj"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	ConnPoolSize = 10
)

var (
	// Command flags
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)

	// Connection
	//connMap  sync.Map // for future
	connPool = make(chan net.Conn, ConnPoolSize)

	// Etc
	quit = make(chan bool)
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

	// Debug
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("started debugging")
	}

	// Start server
	go startServer(*host, *port)

	// Wait for stop signal
	waitForSignals()
	log.Info("stopped server. ")
}

func startServer(host string, port int) error {
	// Start listener
	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Error("failed to start server;", err)
		return err
	}
	defer listener.Close()
	log.Info("started server")

	// Start listening
	var connID int64 = 1
	for {
		// Accept new connection
		conn, err := listener.Accept()

		if err != nil {
			select {
			case <-quit:
				return nil
			default:
			}

			continue
		}
		log.Debugf("connection #%d is connected from %s", connID, conn.RemoteAddr().String())
		//connMap.Store(connID, conn) // for future
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
		log.Infof("time: %3.1f, count: %d, average=%3.1f", time.Since(t).Seconds(), success_count, float64(success_count)/time.Since(t).Seconds())
		conn.Close()
		//connMap.Delete(connID)
	}()
	log.Debugf("connection #%d is ready", connID)

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
				log.Debug("disconnected connection %d;", connID, conn.RemoteAddr().String())
				return nil
			}
			log.Error("failed to receive header and decode;", err)
			return err
		}
		log.Debugf("received message header. payload is %d bytes", msgHeader.PayloadSize)

		// Response
		err = encoder.Encode(&obj.Response{Code: 1})
		if err != nil {
			log.Error("failed to response;", err)
			continue
		}

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
		}
		log.Debugf("received payload successfully. payload=%d, received=%d", msgHeader.PayloadSize, received)

		// Verify
		responseCode := 0
		decoder := gob.NewDecoder(bytes.NewReader(data))
		var m obj.Message
		err = decoder.Decode(&m)
		if err != nil {
			log.Error("failed to decode payload;", err)
			responseCode = -1
		} else {
			if m.Verify() {
				log.Debug("checksum ok")
				responseCode = 200
				success_count++
			} else {
				log.Error("checksum error")
				responseCode = -2
			}
		}

		// Response
		err = encoder.Encode(&obj.Response{Code: responseCode})
		if err != nil {
			log.Error("failed to response;", err)
			continue
		}
	}

	return nil
}

func waitForSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		log.Info("signal received, shutting down...")
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
