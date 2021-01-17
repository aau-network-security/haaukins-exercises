package store

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	AddCategory(string, string) error
	AddExercise(string) error
	UpdateCache() error
}

type store struct {
	db     *mongo.Client
	m      sync.RWMutex
	categs map[model.Tag]model.Category
	exs    map[model.Tag]model.Exercise
}

func NewStore(host string, port uint, user string, pass string) (Store, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)).
		SetAuth(options.Credential{
			Username: user,
			Password: pass,
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
		Tag:  categoryTag, // FR,
		Name: name,        // Forensics,
	}
	_, err := collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}

	s.categs[categoryTag] = category
	return nil
}

// AddExercise updates the challenge if any
// Otherwise adds to collection
func (s *store) AddExercise(content string) error {
	s.m.Lock()
	defer s.m.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	excollect := s.db.Database(DB_NAME).Collection(EXER_COLLECTION)

	var ex model.Exercise
	if err := bson.UnmarshalExtJSON([]byte(content), false, &ex); err != nil {
		fmt.Printf("ERROR: %v", err)
		return err
	}
	tags := strings.Split(string(ex.Tag), "_")
	if len(tags) != 2 {
		return fmt.Errorf("The tag of the challenge does not match with requirements !")
	}
	catTag := strings.ToUpper(tags[0])
	tag := model.Tag(tags[1])
	collection := s.db.Database(DB_NAME).Collection(CAT_COLLECTION)
	var categ model.Category

	c := collection.FindOne(ctx, bson.M{"tag": bson.M{"$eq": catTag}})
	if err := c.Decode(&categ); err != nil {
		return err
	}

	ex.ID = primitive.NewObjectID()
	ex.Tag = tag
	ex.Category = categ.ID

	if err := checkExerciseFields(ex); err != nil {
		return err
	}

	_, ok := s.exs[tag]
	if ok {
		log.Printf("exercise is already exists, updating the document ! ")
		// update the exercise if it is already in the database
		excollect.DeleteOne(ctx, bson.M{"tag": bson.M{"$eq": tag}})
	}

	_, err := excollect.InsertOne(ctx, ex)
	if err != nil {
		return err
	}

	s.exs[tag] = ex

	return nil
}

//Used when someone manually change something in the DB
func (s *store) UpdateCache() error {
	s.m.Lock()
	defer s.m.Unlock()

	s.categs = make(map[model.Tag]model.Category)
	s.exs = make(map[model.Tag]model.Exercise)
	return s.initStore()
}
