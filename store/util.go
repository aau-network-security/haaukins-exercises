package store

import (
	"fmt"

	"github.com/aau-network-security/haaukins-exercises/model"
)

var (
	errorMissingTagName        = fmt.Errorf("error! tag or exercise name is missing")
	errorMissingInstance       = fmt.Errorf("error! there should be at the least an instance")
	errorMissingImage          = fmt.Errorf("error! instance image empty")
	errorMissingFlagConfig     = fmt.Errorf("error! there should be at the least a child in an instance")
	errorMissingExerciseFields = fmt.Errorf("error! not all the required fields in a child are present")
)

func checkExerciseFields(ex model.Exercise) error {

	if ex.Name == "" || string(ex.Tag) == "" {
		return errorMissingTagName
	}

	if len(ex.Instance) == 0 {
		return errorMissingInstance
	}

	//todo the checks can be extended to other variable as well
	flags := 0
	for _, i := range ex.Instance {
		if i.Image == "" {
			return errorMissingImage
		}

		flags += len(i.Flags)
		for _, f := range i.Flags {
			if string(f.Tag) == "" || f.Name == "" || f.EnvVar == "" {
				return errorMissingExerciseFields
			}
		}
	}

	if flags == 0 {
		return errorMissingFlagConfig
	}

	return nil
}
