package store

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/aau-network-security/haaukins-exercises/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_NAME         = "exercise_store"
	EXER_COLLECTION = "exercise"
	CAT_COLLECTION  = "category"
)

type Store interface {
	GetExercises() ([]model.Exercise, error)
	GetExercisesByTags([]string) ([]model.Exercise, error)
	GetExerciseByCategory(string) ([]model.Exercise, error)
}

type store struct {
	db *mongo.Client
	m  sync.RWMutex
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

	return &store{db: client}, nil
}

func (s *store) GetExercises() ([]model.Exercise, error) {
	s.m.Lock()
	defer s.m.Unlock()

	collection := s.db.Database(DB_NAME).Collection(EXER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return s.getExercises(ctx, cur)
}

func (s *store) GetExercisesByTags(tags []string) ([]model.Exercise, error) {
	s.m.Lock()
	defer s.m.Unlock()

	collection := s.db.Database(DB_NAME).Collection(EXER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{"tag": bson.M{"$in": tags}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return s.getExercises(ctx, cur)
}

func (s *store) GetExerciseByCategory(cat string) ([]model.Exercise, error) {
	s.m.Lock()
	defer s.m.Unlock()

	collection := s.db.Database(DB_NAME).Collection(EXER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//todo get id from category tag
	id, _ := primitive.ObjectIDFromHex("5fa68a7bccd7b8d3fc142277")
	match := bson.D{{"$match", bson.D{{"cat", id}}}}
	lookup := bson.D{{"$lookup", bson.D{
		{"from", "category"},
		{"localField", "cat"},
		{"foreignField", "_id"},
		{"as", "exercise"},
	}}}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{match, lookup})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return s.getExercises(ctx, cur)
}

func (s *store) getExercises(ctx context.Context, cur *mongo.Cursor) ([]model.Exercise, error) {
	var exercises []model.Exercise
	for cur.Next(ctx) {
		var e model.Exercise
		err := cur.Decode(&e)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}
