package store

import (
	"context"
	"log"
	"time"

	"github.com/aau-network-security/haaukins-exercises/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DB_NAME         = "exercise_store"
	EXER_COLLECTION = "exercise"
	CAT_COLLECTION  = "category"
)

//init the store cache
func (s *store) initStore() error {

	log.Print("INFO initialising the Cache")
	exs, err := s.getExercises()
	if err != nil {
		return err
	}

	for _, e := range exs {
		s.exs[e.Tag] = e
	}

	cat, err := s.getCategory()
	if err != nil {
		return err
	}

	for _, c := range cat {
		s.categs[c.Tag] = c
	}

	log.Printf("INFO exercises present in the cache [%d]", len(s.exs))
	log.Printf("INFO categories present in the cache [%d]", len(s.categs))
	return nil
}

//Get Categories from the DB
func (s *store) getCategory() ([]model.Category, error) {

	collection := s.db.Database(DB_NAME).Collection(CAT_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categs []model.Category
	for cur.Next(ctx) {
		var c model.Category
		err := cur.Decode(&c)
		if err != nil {
			return nil, err
		}
		categs = append(categs, c)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return categs, nil
}

//Get exercises from the DB
func (s *store) getExercises() ([]model.Exercise, error) {

	collection := s.db.Database(DB_NAME).Collection(EXER_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return s.getExercisesFromDB(ctx, cur)
}

func (s *store) getExercisesFromDB(ctx context.Context, cur *mongo.Cursor) ([]model.Exercise, error) {
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
