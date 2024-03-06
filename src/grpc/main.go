package grpc

import (
	"log"
	"net"

	infc "github.com/predixus/analytics_framework/protobufs/go"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Cannnot create listener: %s", err)
	}
}
