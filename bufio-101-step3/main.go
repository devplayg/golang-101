package main

import (
	"fmt"
	"bufio"
)

type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Printf("%q\n", p)
	return len(p), nil
}
func main() {
	w := new(Writer)
	bw := bufio.NewWriterSize(w, 3)
	bw.Write([]byte("abcd"))
}

