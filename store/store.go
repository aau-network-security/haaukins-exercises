package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aau-network-security/haaukins-exercises/model"
)

type Store interface {
	GetExercises() []model.Exercise
	GetExercisesByTags([]string) ([]model.Exercise, error)
	GetExerciseByCategory(string) ([]model.Exercise, error)
	GetCategories() []model.Category
	GetCategoryName(primitive.ObjectID) string
}

type store struct {
	db     *mongo.Client
	m      sync.RWMutex
	categs map[model.Tag]model.Category
	exs    map[model.Tag]model.Exercise
}

//todo pass config file to modify the connection parameters
func NewStore() (Store, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			Username: "root",
			Password: "toor",
		}))
	if err != nil {
		return nil, err
	}

	s := &store{
		db:     client,
		categs: make(map[model.Tag]model.Category),
		exs:    make(map[model.Tag]model.Exercise),
	}

	if err = s.initStore(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *store) GetExercises() []model.Exercise {
	s.m.RLock()
	defer s.m.RUnlock()

	var exercises []model.Exercise
	for _, e := range s.exs {
		exercises = append(exercises, e)
	}

	return exercises
}

func (s *store) GetCategories() []model.Category {
	s.m.RLock()
	defer s.m.RUnlock()

	var categ []model.Category
	for _, c := range s.categs {
		categ = append(categ, c)
	}

	return categ
}

func (s *store) GetCategoryName(obj primitive.ObjectID) string {
	s.m.RLock()
	defer s.m.RUnlock()

	for _, c := range s.categs {
		if c.ID == obj {
			return string(c.Tag)
		}
	}

	return ""
}

func (s *store) GetExercisesByTags(tags []string) ([]model.Exercise, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	var exercises []model.Exercise
	for _, t := range tags {
		e, ok := s.exs[model.Tag(t)]
		if !ok {
			return nil, fmt.Errorf("exercise [%s] not found", t)
		}
		exercises = append(exercises, e)
	}

	return exercises, nil
}

func (s *store) GetExerciseByCategory(cat string) ([]model.Exercise, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	obj, ok := s.categs[model.Tag(cat)]
	if !ok {
		return nil, fmt.Errorf("category not found")
	}

	var exercises []model.Exercise
	for _, e := range s.exs {
		if e.Category == obj.ID {
			exercises = append(exercises, e)
		}
	}

	return exercises, nil
}
