package main

import (
	"context"
	"fmt"
	pb "github.com/devplayg/golang-101/grpc_data"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

var (
	addr   = "127.0.0.1:8808"
	imgDir = filepath.Join("../images")
	images = make([]string, 0)
)

func init() {
	// seed
	rand.Seed(time.Now().UnixNano())

	// read data
	files, err := ioutil.ReadDir(imgDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".jpg") {
			images = append(images, filepath.Join(imgDir, f.Name()))
		}
	}
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewEventReceiverClient(conn)
	for {
		if err := sendMany(client); err != nil {
			log.Println(err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

func send(client pb.EventReceiverClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	event := generateEvent()
	res, err := client.Send(ctx, event)
	if err != nil {
		return err
	}

	log.Println(res.Message)
	return nil
}

func sendMany(client pb.EventReceiverClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count := rand.Intn(3) + 1
	events := make([]*pb.EventRequest, 0)
	for i := 0; i < count; i++ {
		events = append(events, generateEvent())
	}
	eventReq := &pb.EventsRequest{
		Events: events,
	}
	res, err := client.SendMany(ctx, eventReq)
	if err != nil {
		return err
	}
	log.Printf("msg=%s, sent=%d, success=%d\n", res.Message, count, res.Count)
	return nil
}

func generateEvent() *pb.EventRequest {
	img1, _ := ioutil.ReadFile(images[rand.Intn(len(images))])
	img2, _ := ioutil.ReadFile(images[rand.Intn(len(images))])

	n := rand.Intn(2)
	var b bool
	if n == 1 {
		b = true
	}

	return &pb.EventRequest{
		Date:      time.Now().Format(time.RFC3339),
		Camera:    fmt.Sprintf("DS0000%d", rand.Intn(9)+1),
		EventType: int32(rand.Intn(10)),
		Images:    [][]byte{img1, img2},
		Gloves:    b,
		Helmet:    !b,
		Shoes:     b,
	}
}
