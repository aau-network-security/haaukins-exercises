package main

import (
	"context"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

func (s *Server) GetExercises(ctx context.Context, empty *pb.Empty) (*pb.GetExercisesResponse, error) {
	panic("implement me")
}

func (s *Server) GetExerciseByTags(ctx context.Context, request *pb.GetExerciseByTagsRequest) (*pb.GetExercisesResponse, error) {
	panic("implement me")
}

func (s *Server) GetExerciseByCategory(ctx context.Context, request *pb.GetExerciseByCategoryRequest) (*pb.GetExercisesResponse, error) {
	panic("implement me")
}

func (s *Server) AddExercise(ctx context.Context, request *pb.AddExerciseRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}

func (s *Server) UpdateExercise(ctx context.Context, request *pb.UpdateExerciseRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}
