package obj

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/devplayg/golang-101/tcp-byte-stream-104/utils"
)

// Message
type Message struct {
	Seq       int64
	Timestamp int64
	Data      []byte
	Hash      []byte
}

// Merge
func (m *Message) Merge() []byte {
	seq, _ := utils.Int64ToByte(m.Seq)
	timestamp, _ := utils.Int64ToByte(m.Timestamp)
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
