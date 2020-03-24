# Sending binary files with gRPC

Generate gRPC service code

    protoc -I . --go_out=plugins=grpc:. data.proto
    
Run server

    go run server/server.go

Run client

    go run client/client.go