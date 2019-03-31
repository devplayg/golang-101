package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
}

func init() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// Set flags
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host     = cmdFlags.String("h", "127.0.0.1", "Host")
		port     = cmdFlags.String("p", "8000", "Port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	// Connect to server
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Send data
	var seq int64 = 1
	var buf bytes.Buffer
	for {
		buf.Reset()
		encoder := gob.NewEncoder(&buf)
		m := Message{
			Seq:       seq,
			Timestamp: time.Now().Unix(),
			Data:      []byte(getRandString(3 + rand.Intn(13))), // Random data
		}

		if err := encoder.Encode(m); err == nil {
			if n, err := conn.Write(buf.Bytes()); err == nil {
				log.Printf("[%3d] len=%-4d data=%-6s, %v", m.Seq, n, m.Data, m.Data)
			} else {
				log.Println("failed to write data;", err)
			}
		} else {
			log.Println("failed to decode", err)
		}

		time.Sleep(1000 * time.Millisecond)
		seq++
	}
}

func getRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
