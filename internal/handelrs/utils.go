package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type jsonUnmarshalSchemas interface {
	genDocBody | findDocBody
}

type validationSchema interface {
	jsonUnmarshalSchemas
}

type unmarshalTypeError struct {
	msg          string
	unmarshalErr *json.UnmarshalTypeError
}

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

var unmarshalErr *json.UnmarshalTypeError

func parseBody[T jsonUnmarshalSchemas](schema T, body io.ReadCloser) (T, error) {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&schema); err != nil {
		if errors.As(err, &unmarshalErr) {
			return schema, unmarshalTypeError{
				msg:          "wrong type provided for field",
				unmarshalErr: unmarshalErr,
			}
		}
	}
	return schema, nil
}

func sendResponse(w http.ResponseWriter, respBody map[string]string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(respBody)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
	}

	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := map[string]string{
		"resultdescription": "Bad",
		"error":             err.Error(),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
}

func isStructureEmpty[T validationSchema](strc T) bool {
	structSchema := new(T)
	return reflect.DeepEqual(strc, structSchema)
}

func split(str, separator string) []string {
	return strings.Split(str, separator)
}
