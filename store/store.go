package store

import (
	"context"

	"github.com/aau-network-security/haaukins-exercises/proto"
)

type Store interface {
	AddExercises(context.Context, []*proto.Exercise) error
	GetExercises(context.Context) ([]*proto.Exercise, error)
	GetExercisesByTags(context.Context, []string) ([]*proto.Exercise, error)
	GetExerciseByCategory(context.Context, string) ([]*proto.Exercise, error)

	AddCategory(context.Context, *proto.Category) error
	GetCategories(context.Context) ([]*proto.Category, error)
	GetCategoriesByName(context.Context, []string) ([]*proto.Category, error)
}
