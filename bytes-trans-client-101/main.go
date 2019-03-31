// blog: Marcio

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	dataSize = 80
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
			log.Println("failed to serialize;", err)
			continue
		}

		// Send data size
		dataSize := int64(len(data))
		_, err = conn.Write(IntToHex(dataSize))
		if err != nil {
			log.Println("failed to send data;", err)
		}

		// Send data
		_, err = conn.Write(data)
		if err != nil {
			log.Println("failed to send;", err)
		}

		if seq%10000 == 0 {
			log.Printf("[%3d] len=%-4d", m.Seq, dataSize)
		}

		seq++
		time.Sleep(10 * time.Millisecond)
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

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
