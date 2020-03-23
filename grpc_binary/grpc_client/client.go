package main

import (
	"context"
	pb "github.com/devplayg/hello_grpc/binary"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)
var addr = "127.0.0.1:50051"

func main() {
	addr :=  "125.132.191.38:50051"
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
//func main() {
//	creds, err := credentials.NewServerTLSFromFile("../cert.pem", "../key.pem")
//	if err != nil {
//		log.Fatalf("failed to create credentials: %v", err)
//	}
//
//	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds)	)
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//
//	client := pb.NewGreeterClient(conn)
//	//images := []string{"gopher001.png", "gopher002.png", "gopher003.png"}
//	images := []string{"gopher001.png"}
//	for _, img := range images {
//		if err := send(client, img); err != nil {
//			log.Println(err.Error())
//		}
//	}
//}



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
