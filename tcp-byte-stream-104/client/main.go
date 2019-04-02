package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"github.com/devplayg/golang-101/tcp-byte-stream-104/obj"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"bytes"
	"io"
)
const (
	PayloadSize = 4
)

var (
	// Command flags
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)

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
	var (
		host  = cmdFlags.String("host", "127.0.0.1", "Host")
		port  = cmdFlags.Int("port", 8000, "Port")
		debug = cmdFlags.Bool("debug", false, "debug")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("start debugging")
		Debug = *debug
	}

	conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Error("failed to connect to server;", err)
		return
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
		log.Info("Closed connection")
	}()

	go run(conn)

	// Start server
	//go run(*host, *port)

	// Wait for stop signal
	waitForSignals()
}

func run(conn net.Conn) error {
	var seq int64 = 1
	for {
		// Create message
		m := obj.Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      getData(PayloadSize),
		}
		merged := m.Merge()
		hash := sha256.Sum256(merged)
		m.Hash = hash[:]

		// Serialize message
		data, err := m.Serialize()
		if err != nil {
			log.Error("failed to serialize;", err)
			return err
		}

		// Send message
		reader := bytes.NewReader(data)
		n, err := io.Copy(conn, reader)
		if err != nil {
			log.Error("failed to send data;", err)
			//continue
		}
		log.Debugf("written=%d", n)

		time.Sleep(5*time.Second)
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

func getData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}
