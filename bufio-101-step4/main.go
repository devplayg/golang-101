package main

import (
	"fmt"
	"bufio"
)

type Writer1 int
func (*Writer1) Write(p []byte) (n int, err error) {
	fmt.Printf("writer#1: %q\n", p)
	return len(p), nil
}
type Writer2 int
func (*Writer2) Write(p []byte) (n int, err error) {
	fmt.Printf("writer#2: %q\n", p)
	return len(p), nil
}
func main() {
	w1 := new(Writer1)
	bw := bufio.NewWriterSize(w1, 2)
	bw.Write([]byte("ab"))
	bw.Write([]byte("cd")) // missing

	w2 := new(Writer2)
	bw.Reset(w2)
	bw.Write([]byte("ef"))
	bw.Flush()
}
