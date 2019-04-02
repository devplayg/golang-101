package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
	"bytes"
	"io"
	"encoding/gob"
	"github.com/devplayg/golang-101/tcp-byte-stream-104/obj"
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
	Debug = false
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
		debug = cmdFlags.Bool("debug", true, "debug")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Check debug mode
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("started debugging")
		Debug = *debug
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
		//listener.Close()
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
		log.Debugf("new connection detected #%d", connID)

		connMap.Store(connID, conn)
		go handleConnection(conn, connID)
		connID++
	}

	return nil
}

func handleConnection(conn net.Conn, connID int64) error {
	defer func() {
		log.Debugf("Closing connection #", connID)
		conn.Close()
		connMap.Delete(connID)
	}()

	log.Debugf("Conn #%d is ready", connID)
	for {
		var buf bytes.Buffer
		writer := io.Writer(&buf)
		// accept file from client & write to new file
		msgLen, err := io.Copy(writer, conn)
		if err != nil {
			log.Error("failed to received;", err)
			return err
		}
		log.Debugf("meglen=%d", msgLen)

		//m, err := Deserialize(buf.Bytes())
		//if err != nil {
		//	log.Error("failed to deserialize;", err)
			//			log.WithFields(log.Fields{
			//				"dataSize": dataSize,
			//				"realSize": len(data),
			//			}).Debug("failed to deserialize;", err)
			//			continue
			//return
		//}
		//if m.Verify() {
		//	log.Debugf("seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, true)
		//} else {
		//	log.Infof("seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, false)
		//}
		time.Sleep(5 * time.Second)
	}

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

