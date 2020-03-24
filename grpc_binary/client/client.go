package main

import (
	"context"
	pb "github.com/devplayg/golang-101/grpc_binary"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

var addr = "127.0.0.1:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	images := []string{"gopher001.png", "gopher002.png", "gopher003.png"}
	for _, img := range images {
		if err := send(client, img); err != nil {
			log.Println(err.Error())
		}
	}
}

func send(client pb.GreeterClient, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	res, err := client.SayHello(ctx, &pb.DataRequest{Name: filepath.Base(path), Data: data})
	if err != nil {
		return err
	}
	log.Println(res.Message)
	return nil
}
