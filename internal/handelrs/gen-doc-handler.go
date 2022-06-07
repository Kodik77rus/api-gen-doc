package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/Kodik77rus/api-gen-doc/internal/services"

	// "strings"

	templatebuilder "github.com/Kodik77rus/api-gen-doc/internal/template-builder"
	"github.com/julienschmidt/httprouter"
)

var (
	errEmptyBody      = errors.New("empty body")
	errEmptyUseField  = errors.New("empty \"use\" field")
	errEmptyTextField = errors.New("empty \"text\" field")
	errEmptyURLField  = errors.New("empty \"Url\" field")
	errNegativeId     = errors.New("negative id")
	errBadUseData     = errors.New("field \"use\" must consist of 3 values separated by commas")
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

func GetGenDocHandler() httprouter.Handle {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var parsedBody requestBody
		var unmarshalErr *json.UnmarshalTypeError

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&parsedBody); err != nil {
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

		if err := bodyValidator(&parsedBody); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		template, err := services.HttpResponse(parsedBody.UrlTemplate)
		if err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		templateName := strings.Split(parsedBody.UrlTemplate, "/")
		if len(templateName) != 3 {
			errorResponse(w, errBadUseData, http.StatusBadRequest)
			return
		}

		t := templatebuilder.Template{
			FolderId:     parsedBody.RecordID,
			TemplateName: templateName[len(templateName)-1],
			Template:     &template,
			InsertData: templatebuilder.InsertData{
				Text: parsedBody.Text,
				Use:  parsedBody.Use,
			},
		}

		if err := templatebuilder.New(conf.TemplateBuilder, t).BuildTemplate(); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		sendResponse(w, parsedBody, http.StatusOK)
	}
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

func sendResponse(w http.ResponseWriter, b requestBody, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"resultdescription": "Ok",
		"resultdata":        fmt.Sprintf("%s, %s, %s", b.Text, b.Use, b.Use),
	}

	jsonResp, err := json.Marshal(resp)
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

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

func (rb *requestBody) IsStructureEmpty() bool {
	return reflect.DeepEqual(rb, requestBody{}) //
}
