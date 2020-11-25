package server

import (
	"context"
	"log"

	"github.com/aau-network-security/haaukins-exercises/model"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

func (s *Server) parseExercise(exercisesStore []model.Exercise) []*pb.Exercise {
	var exercises []*pb.Exercise

	for _, e := range exercisesStore {
		var instance []*pb.ExerciseInstance

		for _, c := range e.Instance {
			var children []*pb.ChildExercise
			var envs []*pb.EnvVariable
			var records []*pb.Records

			for _, x := range c.Flags {
				children = append(children, &pb.ChildExercise{
					Tag:                  string(x.Tag),
					Name:                 x.Name,
					EnvFlag:              x.EnvVar,
					Points:               int32(x.Points),
					TeamDescription:      x.TeamDescription,
					OrganizerDescription: x.TeamDescription,
					Prerequisite:         x.PreRequisites,
					Outcome:              x.Outcomes,
				})
			}

			for _, v := range c.Envs {
				envs = append(envs, &pb.EnvVariable{
					Name:  v.EnvVar,
					Value: v.Value,
				})
			}

			for _, r := range c.Records {
				records = append(records, &pb.Records{
					Type: r.Type,
					Name: r.Name,
					Data: r.RData,
				})
			}

			instance = append(instance, &pb.ExerciseInstance{
				Image:    c.Image,
				Memory:   int32(c.MemoryMB),
				Cpu:      float32(c.CPU),
				Envs:     envs,
				Records:  records,
				Children: children,
			})
		}

		exercises = append(exercises, &pb.Exercise{
			Tag:      string(e.Tag),
			Name:     e.Name,
			Status:   int32(e.Status),
			Category: s.store.GetCategoryName(e.Category),
			Instance: instance,
		})
	}
	return exercises
}

func (s *Server) GetExercises(ctx context.Context, empty *pb.Empty) (*pb.GetExercisesResponse, error) {
	log.Print("GET all exercises")
	return &pb.GetExercisesResponse{Exercises: s.parseExercise(s.store.GetExercises())}, nil
}

func (s *Server) GetExerciseByTags(ctx context.Context, request *pb.GetExerciseByTagsRequest) (*pb.GetExercisesResponse, error) {
	tags := request.Tag
	log.Printf("GET exercises by tag %s", tags)
	exs, err := s.store.GetExercisesByTags(tags)
	if err != nil {
		return nil, err
	}
	return &pb.GetExercisesResponse{Exercises: s.parseExercise(exs)}, nil
}

func (s *Server) GetExerciseByCategory(ctx context.Context, request *pb.GetExerciseByCategoryRequest) (*pb.GetExercisesResponse, error) {
	category := request.Category
	log.Printf("GET exercises by category %s", category)
	exs, err := s.store.GetExerciseByCategory(category)
	if err != nil {
		return nil, err
	}
	return &pb.GetExercisesResponse{Exercises: s.parseExercise(exs)}, nil
}

func (s *Server) GetCategories(ctx context.Context, empty *pb.Empty) (*pb.GetCategoriesResponse, error) {

	log.Print("GET all categories")
	var categs []*pb.GetCategoriesResponse_Category
	for _, c := range s.store.GetCategories() {
		categs = append(categs, &pb.GetCategoriesResponse_Category{
			Tag:  string(c.Tag),
			Name: c.Name,
		})
	}

	return &pb.GetCategoriesResponse{Categories: categs}, nil
}

func (s *Server) AddCategory(ctx context.Context, request *pb.AddCategoryRequest) (*pb.ResponseStatus, error) {

	log.Printf("ADD category [%s]", request.Tag)
	if err := s.store.AddCategory(request.Tag, request.Name); err != nil {
		return nil, err
	}
	return &pb.ResponseStatus{}, nil
}

func (s *Server) AddExercise(ctx context.Context, request *pb.AddExerciseRequest) (*pb.ResponseStatus, error) {
	log.Printf("ADD exercise [%s]", request.Tag)
	if err := s.store.AddExercise(request.Tag, request.Content, request.CategoryTag); err != nil {
		return nil, err
	}
	return &pb.ResponseStatus{}, nil
}

func (s *Server) UpdateExercise(ctx context.Context, request *pb.UpdateExerciseRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}

func (s *Server) UpdateCategory(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ResponseStatus, error) {
	panic("implement me")
}
