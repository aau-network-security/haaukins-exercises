package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/aau-network-security/haaukins-exercises/server"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"

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
		log.Fatal().Err(err).Str("config file", *confFilePtr).Msg("Failed to read configuration file")
	}

	s, err := server.NewServer(*c)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create new exercise database server")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatal().Err(err).Uint("port", c.Port).Msg("failed listen on port")
	}

	opts, err := s.GrpcOpts(*c)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server option")
	}

	gRPCServer := s.NewGRPCServer(opts...)
	pb.RegisterExerciseStoreServer(gRPCServer, s)
	reflection.Register(gRPCServer)

	log.Info().Uint("port", c.Port).Msg("serving clients on port")
	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve grpc")
	}
}
