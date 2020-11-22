package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aau-network-security/haaukins-exercises/model"
)

type Store interface {
	GetExercises() []model.Exercise
	GetExercisesByTags([]string) ([]model.Exercise, error)
	GetExerciseByCategory(string) ([]model.Exercise, error)
	AddCategory(string, string) error
	AddExercise(string, string, string) error
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

func (s *store) AddCategory(tag string, name string) error {
	s.m.Lock()
	defer s.m.Unlock()

	categoryTag := model.Tag(tag)
	_, ok := s.categs[categoryTag]
	if ok {
		return fmt.Errorf("category already exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(DB_NAME).Collection(CAT_COLLECTION)

	category := model.Category{
		Tag:  categoryTag,
		Name: name,
	}
	_, err := collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}

	s.categs[categoryTag] = category
	return nil
}

func (s *store) AddExercise(tag string, content string, catTag string) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.exs[model.Tag(tag)]
	if ok {
		return fmt.Errorf("exercise already exists")
	}

	collection := s.db.Database(DB_NAME).Collection(CAT_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categ model.Category
	c := collection.FindOne(ctx, bson.M{"tag": bson.M{"$eq": catTag}})
	if err := c.Decode(&categ); err != nil {
		return err
	}

	var ex model.Exercise
	if err := bson.UnmarshalExtJSON([]byte(content), false, &ex); err != nil {
		return err
	}

	ex.ID = primitive.NewObjectID()
	ex.Category = categ.ID

	if err := checkExerciseFields(ex); err != nil {
		return err
	}

	collection = s.db.Database(DB_NAME).Collection(EXER_COLLECTION)
	_, err := collection.InsertOne(ctx, ex)
	if err != nil {
		return err
	}

	s.exs[model.Tag(tag)] = ex

	return nil
}
