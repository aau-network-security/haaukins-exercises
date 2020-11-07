package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
	"github.com/aau-network-security/haaukins-exercises/store"
	"google.golang.org/grpc"
)

//todo add authentication
type Server struct {
	store *store.Store
}

func main() {

	//todo create a read from config.yml file

	sst, err := store.NewStore()
	if err != nil {
		panic(err)
	}

	s := &Server{store: sst}

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterExerciseStoreServer(gRPCServer, s)
	fmt.Println("waiting client")
	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
