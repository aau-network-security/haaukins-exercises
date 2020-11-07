package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aau-network-security/haaukins-exercises/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddRandomData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			Username: "root",
			Password: "toor",
		}))
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Insert Category
	collection := client.Database(DB_NAME).Collection(CAT_COLLECTION)
	catID := primitive.NewObjectID()

	_, err = collection.InsertOne(ctx, model.Category{
		ID:   catID,
		Tag:  "binary",
		Name: "Binary",
	})

	if err != nil {
		return err
	}

	collection = client.Database(DB_NAME).Collection(EXER_COLLECTION)
	exID := primitive.NewObjectID()
	ex := model.Exercise{
		ID:       exID,
		Category: catID,
		Tag:      model.Tag("test"),
		Name:     "Test Exercise",
		Instance: []model.ExerciseInstanceConfig{
			{
				Image:    "Test",
				MemoryMB: 10,
				CPU:      10,
				Envs:     nil,
				Flags: []model.FlagConfig{
					{
						Tag:         "ex1",
						Name:        "Test",
						EnvVar:      "FLAG",
						Points:      10,
						Description: "this is a test",
					},
				},
				Records: nil,
			},
		},
	}

	_, err = collection.InsertOne(ctx, ex)
	if err != nil {
		return err
	}

	return nil

}

func TestStore_GetExercises(t *testing.T) {
	s, err := NewStore()
	if err != nil {
		t.Fatalf("Error creating the store: %v", err)
	}

	err = AddRandomData()
	if err != nil {
		t.Fatalf("Error adding random data to the db: %v", err)
	}

	exers, err := s.GetExercises()
	if err != nil {
		t.Fatalf("Error get exercises: %v", err)
	}
	fmt.Println(exers)
}
