package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cshep4/grpc-course/grpc-hello-world-sevalla/proto"
)

func main() {
	ctx := context.Background()

	host, ok := os.LookupEnv("GRPC_HOST")
	if !ok {
		host = "localhost:50051"
	}

	log.Printf("gRPC host: %s", host)

	tlsCredentials := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.NewClient(host,
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := client.SayHello(ctx, &proto.SayHelloRequest{Name: "Chris"})
		if err != nil {
			log.Printf("gRPC error: %s", err)
			http.Error(w, err.Error(), 500)
		}

		// return file contents to user
		if _, err := w.Write([]byte(res.GetMessage())); err != nil {
			log.Printf("write error: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Printf("starting http server on address: :%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
