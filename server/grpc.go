package server

import (
	"context"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

func (s *Server) GetExercises(ctx context.Context, empty *pb.Empty) (*pb.GetExercisesResponse, error) {
	var exercises []*pb.Exercise

	for _, e := range s.store.GetExercises() {
		var instance []*pb.ExerciseInstance

		for _, c := range e.Instance {
			var children []*pb.ChildExercise

			for _, x := range c.Flags {
				children = append(children, &pb.ChildExercise{
					Tag:             string(x.Tag),
					Name:            x.Name,
					EnvFlag:         x.EnvVar,
					Points:          int32(x.Points),
					TeamDescription: x.TeamDescription,
				})
			}

			instance = append(instance, &pb.ExerciseInstance{
				Image:    c.Image,
				Memory:   int32(c.MemoryMB),
				Cpu:      float32(c.CPU),
				Children: children,
			})
		}

		exercises = append(exercises, &pb.Exercise{
			Tag:      string(e.Tag),
			Name:     e.Name,
			Category: s.store.GetCategoryName(e.Category),
			Instance: instance,
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

func (s *Server) GetCategories(ctx context.Context, empty *pb.Empty) (*pb.GetCategoriesResponse, error) {
	panic("implement me")
}

func (s *Server) AddCategory(ctx context.Context, request *pb.AddCategoryRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}

func (s *Server) AddExercise(ctx context.Context, request *pb.AddExerciseRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}

func (s *Server) UpdateExercise(ctx context.Context, request *pb.UpdateExerciseRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}

func (s *Server) UpdateCategory(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}
