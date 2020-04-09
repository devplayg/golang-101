package main

import (
	"context"
	"fmt"
	pb "github.com/devplayg/golang-101/grpc_data"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	addr, imgDir string
	images       = make([]string, 0)
)

func init() {
	if len(os.Args) < 3 {
		os.Exit(0)
	}
	addr = os.Args[1]
	imgDir = os.Args[2]
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewEventReceiverClient(conn)
	for {

		files, _ := readDir(imgDir)
		for _, path := range files {
			if err := send(client, path); err != nil {
				fmt.Println("[error] " + err.Error())
				continue
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func readDir(dir string) ([]string, error) {
	println(dir)
	files := make([]string, 0)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		if !f.Mode().IsRegular() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files, err
}

func send(client pb.EventReceiverClient, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event := generateEvent(path)
	res, err := client.Send(ctx, event)
	if err != nil {
		return err
	}

	log.Println(res.Message)
	return nil
}

//
//func sendMany(client pb.EventReceiverClient) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	count := rand.Intn(3) + 1
//	events := make([]*pb.EventRequest, 0)
//	for i := 0; i < count; i++ {
//		events = append(events, generateEvent())
//	}
//	eventReq := &pb.EventsRequest{
//		Events: events,
//	}
//	res, err := client.SendMany(ctx, eventReq)
//	if err != nil {
//		return err
//	}
//	log.Printf("msg=%s, sent=%d, success=%d\n", res.Message, count, res.Count)
//	return nil
//}

func generateEvent(path string) *pb.EventRequest {
	img, _ := ioutil.ReadFile(path)

	//n := rand.Intn(2)
	//var b bool
	//if n == 1 {
	//	b = true
	//}

	return &pb.EventRequest{
		Date:      time.Now().Format(time.RFC3339),
		Camera:    fmt.Sprintf("DS0000%d", rand.Intn(9)+1),
		EventType: int32(rand.Intn(10)),
		Images:    [][]byte{img},
		Gloves:    false,
		Helmet:    false,
		Shoes:     false,
	}
}
