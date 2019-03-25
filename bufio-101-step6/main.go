package main

import (
	"strings"
	"bufio"
	"fmt"
)

func main() {
	s1 := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := bufio.NewReaderSize(s1, 16)
	b, _ := r.Peek(3)
	fmt.Printf("%q\n", b)
	r.Read(make([]byte, 1))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	r.Read(make([]byte, 15))
	fmt.Printf("%q\n", b)
}
