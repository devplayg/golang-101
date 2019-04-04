package main

import (
	"flag"
	"os"
	"net"
	"log"
	"bufio"
	"fmt"
)

func main() {
	var (
		cmdFlags = flag.NewFlagSet("", flag.ExitOnError)
		host = cmdFlags.String("h", "127.0.0.1", "Host")
		port = cmdFlags.String("p", "8000", "Port")
	)
	cmdFlags.Usage = func() {
		cmdFlags.PrintDefaults()
	}
	cmdFlags.Parse(os.Args[1:])

	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Fprintf(conn, text+"\n")
	}
}
