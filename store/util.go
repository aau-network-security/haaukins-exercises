package store

import (
	"fmt"

	"github.com/aau-network-security/haaukins-exercises/model"
)

var (
	MissingTagNameErr        = fmt.Errorf("error! tag or exercise name is missing")
	MissingInstanceErr       = fmt.Errorf("error! there should be at the least an instance")
	MissingImageErr          = fmt.Errorf("error! instance image empty")
	MissingFlagConfigErr     = fmt.Errorf("error! there should be at the least a child in an instance")
	MissingExerciseFieldsErr = fmt.Errorf("error! not all the required fields in a child are present")
)

func checkExerciseFields(ex model.Exercise) error {

	if ex.Name == "" || string(ex.Tag) == "" {
		return MissingTagNameErr
	}

	if len(ex.Instance) == 0 {
		return MissingInstanceErr
	}

	//todo the checks can be extended to other variable as well
	flags := 0
	for _, i := range ex.Instance {
		if i.Image == "" {
			return MissingImageErr
		}

		flags += len(i.Flags)
		for _, f := range i.Flags {
			if string(f.Tag) == "" || f.Name == "" || (f.EnvVar == "" && f.StaticFlag == "") {
				return MissingExerciseFieldsErr
			}
		}
	}

	if flags == 0 {
		return MissingFlagConfigErr
	}

	return nil
}
