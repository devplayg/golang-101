package main

// https://github.com/grpc/grpc-go/blob/master/examples/features/encryption/TLS/server/main.go

import (
	"context"
	pb "github.com/devplayg/golang-101/grpc_binary"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.DataRequest) (*pb.DataResponse, error) {
	if err := ioutil.WriteFile(in.Name, in.Data, 0644); err != nil {
		return &pb.DataResponse{
			Message: err.Error(),
		}, err
	}
	return &pb.DataResponse{
		Message: "saved " + in.GetName(),
	}, nil
}

func main() {
	creds, err := credentials.NewServerTLSFromFile("../../cert.pem", "../../key.pem")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	gRpcServer := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(gRpcServer, &server{})
	if err := gRpcServer.Serve(ln); err != nil {
		panic(err)
	}
}
