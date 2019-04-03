package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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
	PayloadSize = 1024 * 1024
)

var (
	// Flags
	cmdFlags = flag.NewFlagSet("", flag.ExitOnError)

	// Channels
	senderQuit   = make(chan bool)
	receiverQuit = make(chan bool)
	response     = make(chan obj.Response)
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
		log.Debug("started debugging")
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
		log.Info("closed connection")
	}()

	go startReceiver(conn)
	go startSender(conn)

	// Wait for stop signal
	waitForSignals()
}

func startSender(conn net.Conn) error {
	defer func() {
		log.Debug("closing connection")
		close(senderQuit)
		conn.Close()
	}()
	var seq int64 = 1
	encoder := gob.NewEncoder(conn)
	for {
		// Create message
		msg := obj.Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      getData(PayloadSize),
		}
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
			case <-receiverQuit:
				return nil
			default:
			}
			log.Debug("failed to send message header")
			return err
		}

		// Waiting for response to send payload
		log.Debug("waiting for signal to transmit payload")
		result := <-response
		log.Debugf("permitted from server to send payload; code=%d", result.Code)

		//header, err := utils.Serialize(msgHeader)
		//if err != nil {
		//	log.Error("failed to serialize;", err)
		//	continue
		//}
		//_, err = conn.Write(header)
		//if err != nil {
		//	log.Error(err)
		//	continue
		//}
		//<-response

		// Send message
		reader := bytes.NewReader(data)
		n, err := io.Copy(conn, reader)
		if err != nil {
			log.Error("failed to send data;", err)
			continue
		}
		log.Debugf("sent payload, %d bytes", n)
		//spew.Dump(data)

		result = <-response
		log.Debugf("sent data successfully; code=%d", result.Code)

		//time.Sleep(5 * time.Second)
		//time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func startReceiver(conn net.Conn) error {
	defer func() {
		log.Debug("closing connection")
		close(receiverQuit)
		//conn.Close()
	}()

	resp := obj.Response{}
	decoder := gob.NewDecoder(conn)
	for {
		log.Debug("receiver is listening")
		err := decoder.Decode(&resp)
		if err != nil {
			select {
			case <-receiverQuit:
				return nil
			default:
			}
			log.Error("failed to decode response;", err)
			return err
			//select {
			//case <-quit:
			//	return nil
			//default:
			//}
			//
			//continue
		}

		log.Debugf("got message. respcode is %d", resp.Code)
		response <- resp
	}

	//m := &Message{}

	//reponse := make([]byte, 1024)
	//reader := bufio.NewReader(conn)
	//
	//reader.Read(response)

	//reader.
	//var err error
	//for err != nil {
	//	reader.ReadByte()
	//}

	//var buf bytes.Buffer

	//bytes.bu
	//response
	//bytes.NewReader()
	//conn.re
	//var buf bytes.Buffer
	//bufio.new
	//reader  := bufio.NewReader(conn)
	//buf := make([]byte, 1024)
	//var b bytes.Buffer
	//bytes.NewReader()
	//for {
	//reader.Reset()
	//reader.re
	//n, err := conn.Read(buf)
	//if err != nil {
	//	log.Error("failed to read response;", err)
	//	return err
	//}
	//}

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
