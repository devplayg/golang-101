package main

import (
	"context"
	"fmt"
	pb "github.com/devplayg/golang-101/grpc_data"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var addr = ":8808"

func init() {
	rand.Seed(time.Now().UnixNano())
	os.Mkdir("storage", 0755)

}

type server struct {
	pb.UnimplementedEventReceiverServer
}

func (s *server) Send(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	t, _ := time.Parse(time.RFC3339, in.Date)
	for i, img := range in.Images {
		path := filepath.Join("storage", t.Format("20060102150405")+"_"+getRandString(5)+"_"+strconv.Itoa(i)+".data")
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

func getRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	gRpcServer := grpc.NewServer()
	pb.RegisterEventReceiverServer(gRpcServer, &server{})
	if err := gRpcServer.Serve(ln); err != nil {
		panic(err)
	}
}
