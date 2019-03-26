package obj
//
//import (
//	"bytes"
//	"encoding/gob"
//)
//
//type Message struct {
//	Seq       int64
//	Timestamp int64
//	Type      uint8
//	Data      []byte
//}
//
//func (m *Message) Serialize() ([]byte, error) {
//	var buf bytes.Buffer
//	encoder := gob.NewEncoder(&buf)
//
//	err := encoder.Encode(m)
//	if err != nil {
//		return nil, err
//	}
//
//	return buf.Bytes(), nil
//}
//
//func (m *Message) Serialize2() ([]byte, error) {
//	var buf bytes.Buffer
//	if err := gob.NewEncoder(&buf).Encode(m); err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}
//
//func Deserialize(b []byte) (*Message, error) {
//	var m Message
//
//	reader := bytes.NewReader(b)
//	decoder := gob.NewDecoder(reader)
//
//	err := decoder.Decode(&m)
//	if err != nil {
//		return nil, err
//	}
//
//	return &m, nil
//}
//
//func Deserialize2(b []byte) (*Message, error) {
//	var m Message
//	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&m); err != nil {
//		return nil, err
//	}
//	return &m, nil
//}
