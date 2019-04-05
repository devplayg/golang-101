package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"flag"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"github.com/devplayg/golang-101/tcp-cam-stream-103/obj"
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
	responseChan     = make(chan obj.Response)
	exitChan         = make(chan os.Signal)
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
		camID = cmdFlags.Int("cam", 0, "Camera ID")
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
	defer conn.Close()

	// Open camera
	webCam, err := gocv.OpenVideoCapture(*camID)
	if err != nil {
		log.Errorf("failed to open camera #%d", camID)
		return
	}
	defer webCam.Close()

	// Start signal receiver
	go startReceiver(conn)

	// Start stream sender
	go startSender(conn, webCam)

	// Wait for stop signal
	waitForSignals()
}

func startSender(conn net.Conn, webCam *gocv.VideoCapture) error {
	img := gocv.NewMat()
	defer func() {
		log.Debug("closing connection")
		close(senderQuitChan)
		conn.Close()
		img.Close()
		exitChan <- os.Interrupt
	}()
	var seq int64 = 1
	encoder := gob.NewEncoder(conn)

	for {
		// Capture image
		if ok := webCam.Read(&img); !ok {
			msg := "failed to read image from camera"
			log.Debug(msg)
			return errors.New(msg)
		}
		if img.Empty() {
			continue
		}

		// Create message
		msg := obj.Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      img.ToBytes(),
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
			Rows: img.Rows(),
			Cols: img.Cols(),
			MatType: img.Type(),
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
		exitChan <- os.Interrupt
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
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-exitChan:
		log.Info("signal received, shutting down...")
	}
}

func getData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}
