package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"fmt"
)

func main() {
	var i int64
	var length int

	for i< math.MaxUint16 {

		//fmt.Printf("n=%d, len=%d\n", i, length)
		if length != len(IntToHex(i)) {
			fmt.Printf("n=%d, len=%d\n", i, length)
			length = len(IntToHex(i))
		}
		i++
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
