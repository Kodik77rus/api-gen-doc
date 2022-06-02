package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

var (
	errEmptyBody      = errors.New("empty body")
	errEmptyUseField  = errors.New("empty \"use\" field")
	errEmptyTextField = errors.New("empty \"text\" field")
)

type unmarshalTypeError struct {
	msg          string
	unmarshalErr *json.UnmarshalTypeError
}

type requestBody struct {
	Use  string `json:"use"`
	Text string `json:"text"`
}

func GenDocHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.ErrBodyNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var parsedbody requestBody
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&parsedbody); err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(
				w,
				unmarshalTypeError{
					msg:          "wrong type provided for field",
					unmarshalErr: unmarshalErr,
				},
				http.StatusBadRequest,
			)
			return
		}

		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := bodyValidator(&parsedbody); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	sendRespone(w, parsedbody, http.StatusOK)
}

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

func (rb *requestBody) IsStructureEmpty() bool {
	return reflect.DeepEqual(rb, requestBody{}) //
}

func bodyValidator(rb *requestBody) error {
	if rb.IsStructureEmpty() {
		return errEmptyBody
	}

	if len(rb.Text) == 0 {
		return errEmptyTextField
	}

	if len(rb.Use) == 0 {
		return errEmptyUseField
	}

	return nil
}

func sendRespone(w http.ResponseWriter, b requestBody, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp = map[string]string{
		"resultdescription": "Ok",
		"resultdata":        fmt.Sprintf("%s,  %s", b.Use, b.Use),
	}

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp["message"] = err.Error()

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}