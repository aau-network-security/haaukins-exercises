package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/aau-network-security/haaukins-exercises/server"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

const (
	defaultConfigFile = "config.yml"
)

func main() {

	confFilePtr := flag.String("config", defaultConfigFile, "configuration file")
	flag.Parse()

	c, err := server.NewConfigFromFile(*confFilePtr)
	if err != nil {
		log.Fatalf("unable to read configuration file [%s]: %s", *confFilePtr, err)
	}

	s, err := server.NewServer(c)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts, err := s.GrpcOpts(c)
	if err != nil {
		log.Fatalf("failed to retrieve server options %s", err.Error())
	}

	gRPCServer := s.NewGRPCServer(opts...)
	pb.RegisterExerciseStoreServer(gRPCServer, s)
	log.Print("INFO server waiting for clients")
	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
