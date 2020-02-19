//Package helper ...
package helper

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func SetupResponse(res http.ResponseWriter, contentType string, body []byte, statusCode int) {
	res.Header().Set("Content-Type", contentType)
	res.WriteHeader(statusCode)
	if _, err := res.Write(body); err != nil {
		log.Println(err)
	}
}

func ValidateModel(model interface{}) error {
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}
