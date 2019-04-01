package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"io"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"time"
)

const (
	BufferSize = 1024
)

var (
	Debug bool
)

// Message object
type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
	Hash      []byte
}

// Merge data
func (m *Message) Merge() []byte {
	seq := IntToHex(m.Seq)
	ts := IntToHex(m.Timestamp)
	merged := bytes.Join(
		[][]byte{
			seq,
			ts,
			m.Data,
		},
		[]byte(""),
	)

	return merged
}

// Verify data
func (m *Message) Verify() bool {
	data := m.Merge()
	hash := sha256.Sum256(data)
	if bytes.Equal(m.Hash, hash[:]) {
		return true
	}
	return false
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func main() {
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host     = cmdFlags.String("host", "127.0.0.1", "receieve host")
		port     = cmdFlags.String("port", "8000", "receieve port")
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

	// Start receiver
	ln, err := net.Listen("tcp", *host+":"+*port)
	if nil != err {
		log.Panic(err)
	}
	defer ln.Close()
	log.Infof("receiver started listening on %s:%s", *host, *port)

	// Accept connections
	closedConn := make(chan net.Conn, 4)
	go func(ch chan net.Conn) {
		for {
			conn, err := ln.Accept()
			if nil != err {
				log.Error("failed to accept;", err)
				continue
			}
			go handleConnection(conn, ch)
		}
	}(closedConn)


	for {
		select {
		case conn:= <-closedConn :
			err := conn.Close()
			time.Sleep(1*time.Second)
			if err != nil {
				log.Error("failed to close connection", err)
			}
		}
	}
}

func handleConnection(conn net.Conn, closeConn chan<- net.Conn) {
	log.Infof("new connection %v", conn.RemoteAddr().String())
	sizeBuf := make([]byte, 8)


	for {
		buf := make([]byte, BufferSize)
		// Read data size
		n, err := conn.Read(sizeBuf)
		//log.Debug("read")
		if nil != err {
			if io.EOF == err {
				log.Errorf("closed from client; %v", conn.RemoteAddr().String())
				closeConn <- conn
				return
			}
			log.Errorf("fail to receive data; err: %v", err)
			closeConn <- conn
			return
		}
		if 0 < n {
			dataSize := int64(binary.BigEndian.Uint64(sizeBuf[:n]))
			log.Debugf("read=%d, data=%v, size=%d", n, sizeBuf[:n], dataSize)

			var read int64
			data := make([]byte, 0)
			for read < dataSize {
				n, err := conn.Read(buf)
				if err != nil {
					if io.EOF == err {
						log.Errorf("closed from client; %v", conn.RemoteAddr().String())
						closeConn <- conn
						return
					}
					log.Error("failed to read;", err)
					closeConn <- conn
					return
				}
				if (dataSize - read) < BufferSize {
					last := dataSize - read
					read += last
					data = append(data, buf[:last]...)

				} else {
					read += int64(n)
					data = append(data, buf[:n]...)

				}
				//log.Debugf("total=%d, read=%d, merged=%d, \n", dataSize, n, read)
			}
			m, err := Deserialize(data)
			if err != nil {
				//log.Error("failed to deserialize;", err)
				log.WithFields(log.Fields{
					"dataSize": dataSize,
					"realSize": len(data),
				}).Debug("failed to deserialize;", err)
				continue
			}
			if m.Verify() {
				log.Debugf("seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, true)
			} else {
				log.Infof("seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, false)
			}

			//if m.Seq%100 == 0 {
				//log.WithFields(log.Fields{
				//	"sizeData": sizeBuf[:n],
				//	"size":     dataSize,
				//}).Debug()
				//log.Debugf("seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, true)
			//}

		}
	}
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Error("failed to convert int to hex;", num)
		return nil
	}

	return buff.Bytes()
}

func Deserialize(b []byte) (*Message, error) {
	var m Message

	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
