package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/aau-network-security/haaukins-exercises/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBName             = "exercise_store"
	ExerciseCollection = "exercises"
	CategoryCollection = "categories"
)

type store struct {
	db *mongo.Client
}

func NewStore(ctx context.Context, host string, port uint, user string, pass string) (*store, error) {

	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)).
			SetAuth(options.Credential{Username: user, Password: pass}))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	s := &store{
		db: client,
	}

	return s, nil
}

func (s *store) GetExercises(ctx context.Context) ([]*proto.Exercise, error) {
	var exercises []*proto.Exercise

	coll := s.db.Database(DBName).Collection(ExerciseCollection)

	cur, err := coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &exercises); err != nil {
		return nil, err
	}

	return exercises, nil
}

func (s *store) GetExercisesByTags(ctx context.Context, tags []string) ([]*proto.Exercise, error) {
	var exercises []*proto.Exercise

	coll := s.db.Database(DBName).Collection(ExerciseCollection)
	cur, err := coll.Find(ctx, bson.M{"tag": bson.M{"$in": tags}})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &exercises); err != nil {
		return nil, err
	}
	out := handlePrivacyUniverse(exercises)

	return out, nil
}

func (s *store) GetExerciseByCategory(ctx context.Context, cat string) ([]*proto.Exercise, error) {
	var exercises []*proto.Exercise

	coll := s.db.Database(DBName).Collection(ExerciseCollection)

	cur, err := coll.Find(ctx, bson.M{"category": bson.M{"$eq": cat}})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &exercises); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (s *store) GetCategories(ctx context.Context) ([]*proto.Category, error) {
	var cats []*proto.Category

	coll := s.db.Database(DBName).Collection(CategoryCollection)

	cur, err := coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &cats); err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *store) GetCategoriesByName(ctx context.Context, names []string) ([]*proto.Category, error) {
	var cats []*proto.Category

	coll := s.db.Database(DBName).Collection(CategoryCollection)

	cur, err := coll.Find(ctx, bson.M{"name": bson.M{"$in": names}})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &cats); err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *store) AddCategory(ctx context.Context, category *proto.Category) error {
	collection := s.db.Database(DBName).Collection(CategoryCollection)

	_, err := collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) AddExercises(ctx context.Context, exs []*proto.Exercise) error {
	coll := s.db.Database(DBName).Collection(ExerciseCollection)

	for _, ex := range exs {
		filter := bson.M{"tag": bson.M{"$eq": ex.Tag}}
		if _, err := coll.ReplaceOne(ctx, filter, ex, options.Replace().SetUpsert(true)); err != nil {
			return err
		}
	}
	return nil
}

func handlePrivacyUniverse(in []*proto.Exercise) []*proto.Exercise {
	var exercises []*proto.Exercise
	var privEnvs []string
	platforms := make(map[string]bool)
	var flags []*proto.Flags
	var dns []*proto.DNSRecord
	var hosts []*proto.Host

	for _, i := range in {
		if i.Category != "Privacy Universe" {
			exercises = append(exercises, i)
			continue
		}
		privEnvs = append(privEnvs, i.PrivacyEnv)
		for _, h := range i.Hosts {
			flags = append(flags, h.Flags...)
		}

		for _, p := range i.Platforms {
			platforms[p] = true
		}
	}
	if len(platforms) >= 1 {
		for k := range platforms {
			dns = append(dns, &proto.DNSRecord{
				Name: fmt.Sprintf("api.%s.hkn", k),
				Type: "A",
			})
			hosts = append(hosts, &proto.Host{
				Image: fmt.Sprintf("registry.gitlab.com/haaukins/privacy-universe/%s", k),
				Dns: []*proto.DNSRecord{
					{
						Name: fmt.Sprintf("%s.hkn", k),
						Type: "A",
					},
				},
			})
		}
		privacyAPI := proto.Host{
			Image: "registry.gitlab.com/haaukins/privacy-universe/privacy-api",
			Environment: []*proto.EnvVariable{
				{
					Name:  "CHALLS",
					Value: strings.Join(privEnvs, ","),
				},
			},
			Dns:   dns,
			Flags: flags,
		}

		hosts = append(hosts, &privacyAPI)

		privacyUniverse := proto.Exercise{
			Name:  "Privacy Universe",
			Tag:   "pu",
			Hosts: hosts,
		}

		exercises = append(exercises, &privacyUniverse)
	}

	return exercises
}
