package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	templatebuilder "github.com/Kodik77rus/api-gen-doc/internal/template-builder"
	"github.com/julienschmidt/httprouter"
)

var (
	errEmptyBody      = errors.New("empty body")
	errEmptyUseField  = errors.New("empty \"use\" field")
	errEmptyTextField = errors.New("empty \"text\" field")
	errEmptyURLField  = errors.New("empty \"Url\" field")
	errNegativeId     = errors.New("negative id")
)

type requestBody struct {
	Use         string `json:"use"`
	Text        string `json:"text"`
	UrlTemplate string `json:"urlTemplate"`
	RecordID    int    `json:"recordID"`
}

type unmarshalTypeError struct {
	msg          string
	unmarshalErr *json.UnmarshalTypeError
}

func GenDocHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	template, err := getTempalte(parsedbody.UrlTemplate)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	prepareTempalte := templatebuilder.NewTemplate(
		parsedbody.RecordID,
		&template,
		parsedbody.Text,
		parsedbody.Use,
	)

	if err := templatebuilder.New("../").BuildTemplate(*prepareTempalte); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
	}

	sendRespone(w, parsedbody, http.StatusOK)
}

func getTempalte(url string) (string, error) {
	c := http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", nil
	}

	req.Header.Add("User-Agent", "hackerman")
	resp, err := c.Do(req)

	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func bodyValidator(rb *requestBody) error {
	if rb.IsStructureEmpty() {
		return errEmptyBody
	}

	if rb.RecordID < 0 {
		return errNegativeId
	}

	if len(rb.UrlTemplate) == 0 {
		return errEmptyURLField
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
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"resultdescription": "Ok",
		"resultdata":        fmt.Sprintf("%s,  %s", b.Use, b.Use),
	}

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := map[string]string{
		"resultdescription": "Bad",
		"error":             err.Error(),
	}

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

func (rb *requestBody) IsStructureEmpty() bool {
	return reflect.DeepEqual(rb, requestBody{}) //
}
