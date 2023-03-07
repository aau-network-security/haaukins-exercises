package server

import (
	"context"

	"github.com/rs/zerolog/log"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
)

func (s *Server) GetExercises(ctx context.Context, empty *pb.Empty) (*pb.GetExercisesResponse, error) {
	log.Info().Msg("getting all exercises")

	exercises, err := s.store.GetExercises(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetExercisesResponse{Exercises: exercises}, nil
}

func (s *Server) GetExerciseByTags(ctx context.Context, request *pb.GetExerciseByTagsRequest) (*pb.GetExercisesResponse, error) {
	log.Info().Strs("exercise tags", request.Tag).Msg("getting exercises from database")

	exs, err := s.store.GetExercisesByTags(ctx, request.Tag)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve challenges by tags")
		return nil, err
	}
	return &pb.GetExercisesResponse{Exercises: exs}, nil
}

func (s *Server) GetExerciseByCategory(ctx context.Context, request *pb.GetExerciseByCategoryRequest) (*pb.GetExercisesResponse, error) {
	log.Info().Str("category", request.Category).Msg("getting exercise in category")

	exs, err := s.store.GetExerciseByCategory(ctx, request.Category)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve challenges by category")
		return nil, err
	}
	return &pb.GetExercisesResponse{Exercises: exs}, nil
}

func (s *Server) GetCategories(ctx context.Context, empty *pb.Empty) (*pb.GetCategoriesResponse, error) {
	log.Info().Msg("getting all categories")

	categs, err := s.store.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoriesResponse{Categories: categs}, nil
}

func (s *Server) GetCategoriesByName(ctx context.Context, in *pb.GetCategoriesByNameReq) (*pb.GetCategoriesResponse, error) {
	log.Info().Strs("categories", in.Name).Msg("getting categories by name")

	categs, err := s.store.GetCategoriesByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoriesResponse{Categories: categs}, nil
}

func (s *Server) AddCategory(ctx context.Context, request *pb.AddCategoryRequest) (*pb.Empty, error) {
	log.Info().Str("name", request.Category.Name).Msg("adding new category to database")

	if err := s.store.AddCategory(ctx, request.Category); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *Server) AddExercises(ctx context.Context, request *pb.AddExercisesRequest) (*pb.Empty, error) {
	log.Info().Int("amount", len(request.Exercises)).Msg("inserting exercises")
	if err := s.store.AddExercises(ctx, request.Exercises); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
