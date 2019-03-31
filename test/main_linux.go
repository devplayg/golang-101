package main

import (
	"flag"
)
//
//const (
//	BufferSize = 1024
//)
//
//type Message struct {
//	Seq       int64
//	Timestamp int64
//	Data      []byte
//	Hash      []byte
//}
//
//func (m *Message) Merge() []byte {
//	seq := IntToHex(m.Seq)
//	ts := IntToHex(m.Timestamp)
//	merged := bytes.Join(
//		[][]byte{
//			seq,
//			ts,
//			m.Data,
//		},
//		[]byte(""),
//	)
//
//	return merged
//}
//
//func Deserialize(b []byte) (*Message, error) {
//	var m Message
//
//	decoder := gob.NewDecoder(bytes.NewReader(b))
//	err := decoder.Decode(&m)
//	if err != nil {
//		return nil, err
//	}
//
//	return &m, nil
//}

func main() {
	print("asdfasdf\n")
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		recvHost = cmdFlags.String("recvhost", "127.0.0.1", "Receieve host")
		recvPort = cmdFlags.String("recvport", "8000", "receieve port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	//cmdFlags.Parse(os.Args[1:])

	//// Start receiver
	//ln, err := net.Listen("tcp", *recvHost+":"+*recvPort)
	//if nil != err {
	//	panic(err)
	//}
	//defer ln.Close()
	//log.Printf("receiver started listening on %s:%s", *recvHost, *recvPort)
	//
	//// Accept connections
	//for {
	//	conn, err := ln.Accept()
	//	if nil != err {
	//		log.Println("failed to accept;", err)
	//		continue
	//	}
	//	go handleConnection(conn)
	//}
}
//
//func handleConnection(conn net.Conn) {
//	log.Printf("new connection %v", conn.RemoteAddr().String())
//	//buf := make([]byte, 4096)
//	sizeBuf := make([]byte, 10)
//
//	buf := make([]byte, BufferSize)
//
//	for {
//		// Read data size
//		n, err := conn.Read(sizeBuf)
//		if nil != err {
//			if io.EOF == err {
//				log.Printf("closed from client; %v", conn.RemoteAddr().String())
//				return
//			}
//			log.Printf("fail to receive data; err: %v", err)
//			return
//		}
//		if 0 < n {
//			// Read data
//			dataSize := int64(binary.BigEndian.Uint64(sizeBuf[:n]))
//			var read int64
//			data := make([]byte, 0)
//			for read < dataSize {
//				n, err := conn.Read(buf)
//				if err != nil {
//					if io.EOF == err {
//						log.Printf("closed from client; %v", conn.RemoteAddr().String())
//						return
//					}
//					log.Println("failed to read;", err)
//					return
//				}
//				if (dataSize - read) < BufferSize {
//					last := dataSize - read
//					read += last
//					data = append(data, buf[:last]...)
//
//				} else {
//					read += int64(n)
//					data = append(data, buf[:n]...)
//
//				}
//				//log.Printf("total=%d, read=%d, merged=%d, \n", dataSize, n, read)
//			}
//			m, err := Deserialize(data)
//			merged := m.Merge()
//			hash := sha256.Sum256(merged)
//			//m.Hash = hash[:]
//			if err != nil {
//				log.Println("failed to deserialize;", err)
//				//spew.Dump(data)
//				continue
//			}
//			if bytes.Equal(m.Hash, hash[:]) {
//				//log.Printf("==> seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, true)
//			} else {
//				log.Printf("==> seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, false)
//				spew.Dump(m)
//			}
//
//			if m.Seq%10000 == 0 {
//				log.Printf("==> seq=%d, timestamp=%d, equal=%v", m.Seq, m.Timestamp, true)
//			}
//
//		}
//
//		// Read data
//	}
//}
//
//func IntToHex(num int64) []byte {
//	buff := new(bytes.Buffer)
//	err := binary.Write(buff, binary.BigEndian, num)
//	if err != nil {
//		log.Println("failed to convert int to hex;", num)
//		return nil
//	}
//
//	return buff.Bytes()
//}
