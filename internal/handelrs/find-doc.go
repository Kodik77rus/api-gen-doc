package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/julienschmidt/httprouter"
)

const downloadFileRoute = "http://localhost:8080/api/download/"

var (
	templateDir = []string{"word", "pdf"}

	errFolderNotFound = errors.New("not found id folder")
)

type findDocBody struct {
	UrlTemplate string `json:"urlTemplate"`
	RecordID    int    `json:"recordID"`
}

type asyncDirReader struct {
	wordFiles *[]string
	pdfFiles  *[]string

	err error
}

func newAsyncDirReader() *asyncDirReader { return &asyncDirReader{} }

func FindDocs() httprouter.Handle {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var (
			body findDocBody
			wg   sync.WaitGroup
		)

		body, err := parseBody(body, r.Body)
		if err != nil {
			errorResponse(w, err, http.StatusBadRequest)
		}

		if err := findDocBodyValidator(&body); err != nil {
			errorResponse(w, err, http.StatusBadRequest)
			return
		}

		route := generateFilePath(
			conf.TemplateBuilder.TemplateFolder,
			getTemplateName(body.UrlTemplate),
			strconv.Itoa(body.RecordID),
		)

		reader := newAsyncDirReader()

		wg.Add(len(templateDir))

		for _, v := range templateDir {
			go reader.readDir(generateFilePath(route, v), &wg)
		}

		wg.Wait()

		if reader.err != nil {
			errorResponse(w, reader.err, http.StatusBadRequest)
			return
		}

		sendResponse(w, reader.printResult(), http.StatusOK)
	}
}

func (a *asyncDirReader) readDir(pathToDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.ReadDir(pathToDir)
	if errors.As(err, &pathError) {
		a.setError(err)
		return
	}

	fls := new([]string)

	for _, file := range files {
		*fls = append(*fls,
			generateDownloadPath(
				generateFilePath(pathToDir, file.Name())))
	}

	switch path.Base(pathToDir) {
	case templateDir[0]:
		a.setWordFiles(fls)
	case templateDir[1]:
		a.setPdfFiles(fls)
	}
}

func (a *asyncDirReader) setWordFiles(files *[]string) {
	a.wordFiles = files
}

func (a *asyncDirReader) setPdfFiles(files *[]string) {
	a.pdfFiles = files
}

func (a *asyncDirReader) setError(err error) {
	a.err = err
}

func (a *asyncDirReader) printResult() map[string][]string {
	return map[string][]string{
		templateDir[0]: *a.wordFiles,
		templateDir[1]: *a.pdfFiles,
	}
}

func generateDownloadPath(fileName string) string {
	urlStr, _ := url.Parse(
		downloadFileRoute + fileName)
	return strings.Replace(urlStr.String(), "../template/", "", 1)
}

func findDocBodyValidator(b *findDocBody) error {
	if isStructureEmpty(*b) {
		return errEmptyBody
	}

	if b.RecordID < 0 {
		return errNegativeId
	}

	if len(b.UrlTemplate) == 0 {
		return errEmptyURLField
	}

	return nil
}
