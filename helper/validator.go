//Package helper ...
package helper

import (
	"github.com/go-playground/validator/v10"
)

func ValidateModel(model interface{}) error {
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}
