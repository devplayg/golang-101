package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	marand "math/rand"
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
	PayloadSize = 1024 * 1024
)

var (
	// Flags
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)

	// Channels
	senderQuitChan   = make(chan bool)
	receiverQuitChan = make(chan bool)
	responseChan = make(chan obj.Response)
)

func init() {
	// Initialize logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	// Options
	var (
		host  = cmdFlags.String("host", "127.0.0.1", "Host")
		port  = cmdFlags.Int("port", 8000, "Port")
		debug = cmdFlags.Bool("debug", false, "debug")
	)

	// Handle flags
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Debug
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("started debugging")
	}

	// Connect to server
	conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Error("failed to connect to server;", err)
		return
	}
	defer func() {
		// Close connection
		if conn != nil {
			conn.Close()
		}
		log.Info("closed connection")
	}()

	// Start receiver
	go startReceiver(conn)

	// Start sender
	go startSender(conn)

	// Wait for stop signal
	waitForSignals()
}

func startSender(conn net.Conn) error {
	defer func() {
		log.Debug("closing connection")
		close(senderQuitChan)
		conn.Close()
	}()
	var seq int64 = 1
	encoder := gob.NewEncoder(conn)
	for {
		// Create message
		msg := obj.Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      getData(PayloadSize+marand.Intn(1000)),
		}
		//log.Debug(marand.Intn(1000))
		merged := msg.Merge()
		hash := sha256.Sum256(merged)
		msg.Hash = hash[:]

		// Serialize message
		data, err := msg.Serialize()
		if err != nil {
			log.Error("failed to serialize;", err)
			return err
		}

		// Create message header
		msgHeader := obj.MessageHeader{
			PayloadSize: uint32(len(data)),
		}

		// Send message header
		err = encoder.Encode(&msgHeader)
		if err != nil {
			select {
			case <-receiverQuitChan:
				return nil
			default:
			}
			log.Debug("failed to send message header")
			return err
		}

		// Waiting for response to send payload
		log.Debug("waiting for signal to transmit payload")
		result := <-responseChan
		log.Debugf("permitted from server to send payload; code=%d", result.Code)

		// Send message
		reader := bytes.NewReader(data)
		n, err := io.Copy(conn, reader)
		if err != nil {
			log.Error("failed to send data;", err)
			continue
		}
		log.Debugf("sent payload, %d bytes", n)

		result = <-responseChan
		log.Debugf("sent data successfully; code=%d", result.Code)
	}

	return nil
}

func startReceiver(conn net.Conn) error {
	defer func() {
		close(receiverQuitChan)
		log.Debug("sent closing signal to receiver")
	}()

	response := obj.Response{}
	decoder := gob.NewDecoder(conn)
	for {
		log.Debug("receiver is listening")
		err := decoder.Decode(&response)
		if err != nil {
			select {
			case <-receiverQuitChan:
				// if detected stop signal
				return nil
			default:
			}
			log.Error("failed to decode response;", err)
			return err
		}
		log.Debugf("received message from server; responce code is %d", response.Code)
		responseChan <- response
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

func getData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}
