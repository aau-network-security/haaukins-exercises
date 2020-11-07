package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag string

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Tag  Tag                `bson:"tag,omitempty"`
	Name string             `bson:"name,omitempty"`
}

type Exercise struct {
	ID       primitive.ObjectID       `bson:"_id,omitempty"`
	Category primitive.ObjectID       `bson:"category,omitempty"`
	Tag      Tag                      `bson:"tag,omitempty"`
	Instance []ExerciseInstanceConfig `bson:"instance,omitempty"`
}

type ExerciseInstanceConfig struct {
	Image    string         `bson:"image,omitempty"`
	MemoryMB uint           `bson:"memory,omitempty"`
	CPU      float64        `bson:"cpu,omitempty"`
	Envs     []EnvVarConfig `bson:"env,omitempty"`
	Flags    []FlagConfig   `bson:"flags,omitempty"`
	Records  []RecordConfig `bson:"dns,omitempty"`
}

type FlagConfig struct {
	Tag         Tag    `bson:"tag,omitempty"`
	Name        string `bson:"name,omitempty"`
	EnvVar      string `bson:"env,omitempty"`
	StaticFlag  string `bson:"static_flag,omitempty"`
	Points      uint   `bson:"points,omitempty"`
	Description string `bson:"description,omitempty"`
}

type RecordConfig struct {
	Type  string `bson:"type,omitempty"`
	Name  string `bson:"name,omitempty"`
	RData string `bson:"rdata,omitempty"`
}

type EnvVarConfig struct {
	EnvVar string `bson:"env,omitempty"`
	Value  string `bson:"value,omitempty"`
}
