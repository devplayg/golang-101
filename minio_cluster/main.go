package main

import (
	"github.com/minio/minio-go/v6"
	"log"
)

func main() {
	endpoint := "192.168.0.5"
	accessKeyID := "unisem"
	secretAccessKey := "unisem"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup
}
