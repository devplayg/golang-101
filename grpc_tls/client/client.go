package main

import (
	"context"
	pb "github.com/devplayg/golang-101/grpc_binary"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

var addr = "127.0.0.1:50051"

func main() {

	opts := grpc.WithInsecure()

	creds, err := credentials.NewClientTLSFromFile("../../cert.pem", "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts = grpc.WithTransportCredentials(creds)

	// conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithInsecure())
	//conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithInsecure())
	//
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()

	conn, err := grpc.Dial(addr, opts)
	if err != nil {
		panic(err)
	}
	//defer conn.Close()
	//
	//client := pb.NewGreeterClient(conn)
	client := pb.NewGreeterClient(conn)

	images := []string{"gopher001.png"}
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
