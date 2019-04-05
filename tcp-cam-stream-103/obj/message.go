package obj

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"gocv.io/x/gocv"
)

type Response struct {
	Code int
}

type MessageHeader struct {
	Rows        int
	Cols        int
	MatType     gocv.MatType
	PayloadSize uint32 // max: 4,294,967,295
}

// Message
type Message struct {
	Seq       int64 // wiill be deleted
	Timestamp int64 // wiill be deleted
	Data      []byte
	Hash      []byte
}

// Merge
func (m *Message) Merge() []byte {
	seq, _ := Int64ToByte(m.Seq)
	timestamp, _ := Int64ToByte(m.Timestamp)
	return bytes.Join(
		[][]byte{seq, timestamp, m.Data},
		[]byte(""),
	)
}

// Serialize
func (m *Message) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func Int64ToByte(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
