package server

import (
	"context"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

func (s *Server) GetExercises(ctx context.Context, empty *pb.Empty) (*pb.GetExercisesResponse, error) {
	var exercises []*pb.Exercise
	for _, e := range s.store.GetExercises() {
		var children []*pb.ChildExercise

		//todo not the best way to manage it
		//maybe refactor the db struct
		for _, c := range e.Instance {
			for _, x := range c.Flags {
				children = append(children, &pb.ChildExercise{
					Tag:             string(x.Tag),
					Name:            x.Name,
					EnvFlag:         x.EnvVar,
					Points:          int32(x.Points),
					TeamDescription: x.Description,
				})
			}
		}

		exercises = append(exercises, &pb.Exercise{
			Tag:  string(e.Tag),
			Name: e.Name,
			//Category: s.store.GetCategoryName(),
			Children: children,
		})
	}
	return &pb.GetExercisesResponse{Exercises: exercises}, nil
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
