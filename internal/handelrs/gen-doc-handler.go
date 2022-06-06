package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/Kodik77rus/api-gen-doc/internal/services"
	"log"
	"net/http"
	"reflect"
	"strings"

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

		t := templatebuilder.Template{
			FolderId:     parsedBody.RecordID,
			TemplateName: templateName[len(templateName)],
			Template:     &template,
			InsertData: templatebuilder.InsertData{
				Text: parsedBody.Text,
				Use:  parsedBody.Use,
			},
		}

		if err := templatebuilder.New(conf.TemplateBuilder).BuildTemplate(t); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
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
