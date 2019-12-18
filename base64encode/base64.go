package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: base64encode [file]")
		return
	}

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	fmt.Printf(encoded)
}
