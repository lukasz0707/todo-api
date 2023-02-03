package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// func init() {
// 	validate.RegisterValidation("notWhiteSpace", notWhiteSpace)
// }

// Validate validates the input struct
func Validate(payload interface{}) error {
	err := validate.Struct(payload)

	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(
				errs,
				fmt.Sprintf("`%v` doesn't satisfy the `%v` constraint", err.Field(), err.Tag()),
			)
		}

		return errors.New(strings.Join(errs, ","))
	}

	return nil
}

// // String should not contain white space
// func notWhiteSpace(f1 validator.FieldLevel) bool {
// 	return !strings.Contains(f1.Field().String(), " ")

// }
