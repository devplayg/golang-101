// blog: Marcio

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	dataSize = 1024*1024
)

var (
	Debug bool
	timeout = time.Duration(time.Second)
)

type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
	Hash      []byte
}

func (m *Message) Merge() []byte {
	return bytes.Join(
		[][]byte{
			IntToHex(m.Seq),
			IntToHex(m.Timestamp),
			m.Data,
		},
		[]byte(""),
	)
}

func (m *Message) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	//rand.Seed(time.Now().UnixNano())
}

func main() {

	// Set flags
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host     = cmdFlags.String("host", "127.0.0.1", "Host")
		port     = cmdFlags.String("port", "8000", "Port")
		debug    = cmdFlags.Bool("debug", false, "debug")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Debug
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("start debugging")
		Debug = *debug
	}

	// Connect to server
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		log.Error("failed to connect to server;", err)
		return
	}
	defer func() {
		conn.SetDeadline(time.Now().Add(timeout))
		conn.Close()
	}()
	log.Infof("connected to server %s:%s", *host, *port)

	go send(conn)

	waitForSignals(conn)
}

func send(conn net.Conn) {
	var seq int64 = 1
	for {
		// Create message
		m := Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      getData(dataSize),
		}
		merged := m.Merge()
		hash := sha256.Sum256(merged)
		m.Hash = hash[:]

		data, err := m.Serialize()
		if err != nil {
			log.Error("failed to serialize;", err)
			continue
		}

		// Send data size
		dataSize := int64(len(data))
		dataSizeHex := IntToHex(dataSize)
		log.Debugf("dataSize=%d, hex=%v, data=%v", dataSize, dataSizeHex, data[70:90])
		_, err = conn.Write(dataSizeHex)
		if err != nil {
			log.Error("failed to send size;", err)
		} else {
			// Send data
			_, err = conn.Write(data)
			if err != nil {
				log.Error("failed to send data;", err)
			}

			if seq%100 == 0 {
				//log.WithFields(log.Fields{
				//	"data":   data,
				//	"length": len(data),
				//}).Debug()
				log.Debugf("[%3d] len=%-4d", m.Seq, dataSize)
			}

		}

		seq++
		time.Sleep(100 * time.Millisecond)
	}
}

func getData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

func waitForSignals(conn net.Conn) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		log.Info("Signal received, shutting down...")
	}
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
