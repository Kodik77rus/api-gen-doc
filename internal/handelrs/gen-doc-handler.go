package handlers

import (
	"errors"
	"fmt"
	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/Kodik77rus/api-gen-doc/internal/services"
	"log"
	"net/http"

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

type genDocBody struct {
	Use         string `json:"use"`
	Text        string `json:"text"`
	UrlTemplate string `json:"urlTemplate"`
	RecordID    int    `json:"recordID"`
}

func GetGenDocHandler() httprouter.Handle {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var body genDocBody

		body, err := parseBody(body, r.Body)
		if err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		if err := genDocBodyValidator(&body); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		template, err := services.HttpResponse(body.UrlTemplate)
		if err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		}

		templateName := split(body.UrlTemplate, "/")

		t := templatebuilder.Template{
			FolderId:     body.RecordID,
			TemplateName: templateName[len(templateName)-1],
			Template:     &template,
			InsertData: templatebuilder.InsertData{
				Text: body.Text,
				Use:  body.Use,
			},
		}

		if err := templatebuilder.
			New(conf.TemplateBuilder, t).
			BuildTemplate(); err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"resultdescription": "Ok",
			"resultdata": fmt.Sprintf(
				"%s, %s", body.Text, body.Use,
			),
		}

		sendResponse(w, resp, http.StatusOK)
	}
}

func genDocBodyValidator(rb *genDocBody) error {
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

	validUseData := split(rb.Use, ",")
	if len(validUseData) != 3 {
		return errBadUseData
	}

	return nil
}
