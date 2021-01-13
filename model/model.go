package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag string

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Tag  Tag                `bson:"tag,omitempty"`
	Name string             `bson:"name,omitempty"`
}

//todo manage the status somehow
type Exercise struct {
	ID       primitive.ObjectID       `bson:"_id,omitempty"`
	Category primitive.ObjectID       `bson:"categoryid,omitempty"`
	Tag      Tag                      `bson:"tag,omitempty"`
	Name     string                   `bson:"name,omitempty"`
	Instance []ExerciseInstanceConfig `bson:"instance,omitempty"`
	Status   int                      `bson:"status"`
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
	Tag             Tag      `bson:"tag,omitempty"`
	Name            string   `bson:"name,omitempty"`
	EnvVar          string   `bson:"env,omitempty"`
	StaticFlag      string   `bson:"static_flag,omitempty"`
	Points          uint     `bson:"points,omitempty"`
	Category        string   `bson:"category,omitempty"`
	TeamDescription string   `bson:"td,omitempty"`
	OrgDescription  string   `bson:"od,omitempty"`
	PreRequisites   []string `bson:"reqs,omitempty"`
	Outcomes        []string `bson:"outc,omitempty"`
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
