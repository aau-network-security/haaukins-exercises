package store

import (
	"context"
	"testing"
	"time"

	"github.com/aau-network-security/haaukins-exercises/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	NExercises  = []string{"ftp", "xxs", "xxe", "sql", "mitm", "crypto", "shad", "rand", "ccs"}
	NCategories = []string{"forensics", "binary"}
)

const (
	testHost = "localhost"
	testPort = 27017
	testUser = "root"
	testPass = "toor"
)

//This function will add some data in the DB. it will return if in the DB there are already some categories
//this has be done in order to avoid multiple insertion by calling this function in each test
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

	collection := client.Database(DB_NAME).Collection(CAT_COLLECTION)

	//Skip the function if there are already categories and exercises in the DB
	countDocuments, err := collection.CountDocuments(ctx, nil, nil)
	if countDocuments > 0 {
		return nil
	}

	//Insert Category
	var id primitive.ObjectID
	for _, c := range NCategories {
		id = primitive.NewObjectID()
		_, err = collection.InsertOne(ctx, model.Category{
			ID:   id,
			Tag:  model.Tag(c),
			Name: c,
		})
		if err != nil {
			return err
		}
	}

	collection = client.Database(DB_NAME).Collection(EXER_COLLECTION)

	//NB all the exercises will have the same category binary
	for _, e := range NExercises {

		ex := model.Exercise{
			Category: id,
			Tag:      model.Tag(e),
			Name:     e,
			Instance: []model.ExerciseInstanceConfig{
				{
					Image:    e,
					MemoryMB: 10,
					CPU:      10,
					Envs:     nil,
					Flags: []model.FlagConfig{
						{
							Tag:         model.Tag(e),
							Name:        e,
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
	}

	return nil

}

func TestStore_GetExercises(t *testing.T) {

	if err := AddRandomData(); err != nil {
		t.Fatalf("Error adding random data to the db: %v", err)
	}

	s, err := NewStore(testHost, testPort, testUser, testPass)
	if err != nil {
		t.Fatalf("Error creating the store: %v", err)
	}

	exers := s.GetExercises()
	if len(exers) != len(NExercises) {
		t.Fatalf("Expected number of challenges %d, got %d", len(NExercises), len(exers))
	}
}

func TestStore_GetExercisesByTags(t *testing.T) {
	if err := AddRandomData(); err != nil {
		t.Fatalf("Error adding random data to the db: %v", err)
	}

	s, err := NewStore(testHost, testPort, testUser, testPass)
	if err != nil {
		t.Fatalf("Error creating the store: %v", err)
	}

	tt := []struct {
		name     string
		tags     []string
		expected int
		err      bool
	}{
		{name: "Normal Get exercises by tags not empty", tags: NExercises[:4], expected: 4},
		{name: "Normal Get exercises by tags empty", tags: []string{}, expected: 0},
		{name: "Invalid tags", tags: []string{"random"}, expected: 0, err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			exers, err := s.GetExercisesByTags(tc.tags)
			if err != nil && !tc.err {
				t.Fatalf("Error get exercises: %v", err)
			}
			if len(exers) != tc.expected {
				t.Fatalf("Expected number of challenges %d, got %d", tc.expected, len(exers))
			}
		})
	}
}

func TestStore_GetExerciseByCategory(t *testing.T) {
	if err := AddRandomData(); err != nil {
		t.Fatalf("Error adding random data to the db: %v", err)
	}

	s, err := NewStore(testHost, testPort, testUser, testPass)
	if err != nil {
		t.Fatalf("Error creating the store: %v", err)
	}

	tt := []struct {
		name     string
		categ    string
		expected int
		err      bool
	}{
		{name: "Normal Get exercises by category", categ: NCategories[1], expected: len(NExercises)},
		{name: "Normal 2 Get exercises by category", categ: NCategories[0], expected: 0},
		{name: "Invalid category", categ: "test", expected: 0, err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			exers, err := s.GetExerciseByCategory(tc.categ)
			if err != nil && !tc.err {
				t.Fatalf("Error get exercises: %v", err)
			}
			if len(exers) != tc.expected {
				t.Fatalf("Expected number of challenges %d, got %d", tc.expected, len(exers))
			}
		})
	}
}
