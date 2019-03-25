package main

import (
	"fmt"
	"strings"
	"bufio"
)

type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Printf("%q\n", p)
	return len(p), nil
}
func main() {
	s := strings.NewReader("onetwothree")
	w := new(Writer)
	bw := bufio.NewWriterSize(w, 3)
	bw.ReadFrom(s)
	err := bw.Flush()
	if err != nil {
		panic(err)
	}
}
