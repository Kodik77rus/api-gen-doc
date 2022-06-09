package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	unmarshalErr *json.UnmarshalTypeError
	pathError    *fs.PathError
)

type jsonUnmarshalSchemas interface {
	genDocBody | findDocBody
}

type validationSchema interface {
	jsonUnmarshalSchemas
}

type responseBody interface {
	map[string]string | map[string][]string
}

type unmarshalTypeError struct {
	msg          string
	unmarshalErr *json.UnmarshalTypeError
}

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

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

func sendResponse[T responseBody](w http.ResponseWriter, respBody T, httpStatusCode int) {
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

func isStructureEmpty[T validationSchema](strct T) bool {
	structSchema := new(T)
	return reflect.DeepEqual(strct, structSchema)
}

func split(str, separator string) []string {
	return strings.Split(str, separator)
}

func generateFilePath(etc ...string) string {
	return filepath.Join(etc...)
}

func getTemplateName(url string) string {
	templateName := split(url, "/")
	return templateName[len(templateName)-1]
}
