package serve

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func JSONRequest[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return data, err
	}
	return data, nil
}

func JSONValidRequest[T any](w http.ResponseWriter, r *http.Request) (T, error) {

	data, err := JSONRequest[T](w, r)
	if err != nil {
		return data, err
	}

	validate := validator.New()

	err = validate.Struct(data)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return data, err
	}

	return data, nil
}
