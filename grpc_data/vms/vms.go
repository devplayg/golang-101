package main

import (
	"context"
	"fmt"
	pb "github.com/devplayg/golang-101/grpc_data"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"time"
)

type server struct {
	pb.UnimplementedEventReceiverServer
}

func (s *server) Send(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	t, _ := time.Parse(time.RFC3339, in.Date)
	for i, img := range in.Images {
		path := filepath.Join("./storage", t.Format("20060102150405")+"_"+strconv.Itoa(i)+".jpg")
		if err := ioutil.WriteFile(path, img, 0644); err != nil {
			log.Printf(err.Error())
			return &pb.EventResponse{
				Message: "failed " + in.Date,
			}, nil
		}
	}
	return &pb.EventResponse{
		Message: fmt.Sprintf("saved #%d", 1),
	}, nil
}

func (s *server) SendMany(ctx context.Context, in *pb.EventsRequest) (*pb.EventResponse, error) {
	log.Printf("received %d\n", len(in.Events))
	success := 0
	for _, e := range in.Events {
		t, _ := time.Parse(time.RFC3339, e.Date)
		for i, img := range e.Images {
			path := filepath.Join("./storage", t.Format("20060102150405")+"_"+strconv.Itoa(i)+".jpg")
			if err := ioutil.WriteFile(path, img, 0644); err != nil {
				log.Printf(err.Error())
				continue
			}
		}
		success++
	}
	time.Sleep(5 * time.Second)
	return &pb.EventResponse{
		Message: fmt.Sprintf("saved #%d", success),
		Count:   int32(success),
	}, nil
}

func main() {
	ln, err := net.Listen("tcp", ":8808")
	if err != nil {
		panic(err)
	}

	gRpcServer := grpc.NewServer()
	pb.RegisterEventReceiverServer(gRpcServer, &server{})
	if err := gRpcServer.Serve(ln); err != nil {
		panic(err)
	}
}
