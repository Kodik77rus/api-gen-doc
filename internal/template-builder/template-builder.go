package templatebuilder

import (
	"bytes"
	"fmt"
	"github.com/Kodik77rus/api-gen-doc/internal/services"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type TemplateBuilder struct {
	Config   config.TemplateBuilder
	Template Template
}

type Template struct {
	FolderId     int
	TemplateName string
	Template     *string
	InsertData   InsertData
}

type InsertData struct {
	Text string
	Use  string
}

type wordDoc struct {
	dir        string
	pathToFile string
}

type pdfDoc struct {
	dir        string
	pathToFile string
	content    *string
}

func New(c config.TemplateBuilder, t Template) *TemplateBuilder {
	return &TemplateBuilder{
		Config:   c,
		Template: t,
	}
}

func newWordDoc(t *TemplateBuilder) *wordDoc {
	fullPath := fmt.Sprint(
		t.Config.TemplateFolder,
		"/",
		t.Template.TemplateName,
		"/",
		t.Template.FolderId,
		"/word/",
		time.Now().Format(timeFormat),
		".doc",
	)

	dir, _ := filepath.Split(fullPath)

	return &wordDoc{
		dir:        dir,
		pathToFile: fullPath,
	}
}

func newPdfDoc(wordFileName string, content *string) *pdfDoc {
	replacer := strings.NewReplacer("word", "pdf", ".doc", ".pdf")
	path := replacer.Replace(wordFileName)

	dir, _ := filepath.Split(path)

	return &pdfDoc{
		dir:        dir,
		pathToFile: path,
		content:    content,
	}
}

func (t *TemplateBuilder) BuildTemplate() error {
	wordFile, err := t.createWordFile()
	if err != nil {
		return err
	}

	pdfFile, err := t.convertWordToPdf(wordFile)
	if err != nil {
		return err
	}

	if err := t.savePdf(pdfFile); err != nil {
		return err
	}

	return nil
}

func (t *TemplateBuilder) createWordFile() (*wordDoc, error) {
	word := newWordDoc(t)

	err := os.MkdirAll(word.dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(word.pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	insertData := t.prepareData()
	for _, v := range insertData {
		t.insertData(v)
	}

	_, e := file.WriteString(*t.Template.Template)
	if e != nil {
		return nil, e
	}

	return word, nil
}

func (t *TemplateBuilder) convertWordToPdf(wordFile *wordDoc) (*pdfDoc, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("filename", wordFile.pathToFile)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(wordFile.pathToFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("o", wordFile.pathToFile); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	pdfFile, err := services.MultipartResponse(writer, body)
	if err != nil {
		return nil, err
	}

	return newPdfDoc(wordFile.pathToFile, &pdfFile), nil
}

func (t *TemplateBuilder) savePdf(pdf *pdfDoc) error {
	err := os.MkdirAll(pdf.dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(pdf.pathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, e := file.WriteString(*pdf.content)
	if e != nil {
		return e
	}

	return nil
}

func (t *TemplateBuilder) insertData(v string) {
	d := strings.Replace(*t.Template.Template, "_", v, 1)
	t.setTemplate(&d)
}

func (t *TemplateBuilder) setTemplate(template *string) {
	t.Template.Template = template
}

func (t *TemplateBuilder) prepareData() []string {
	var insertData []string

	useData := strings.Split(t.Template.InsertData.Use, ",")

	insertData = append(insertData, t.Template.InsertData.Text)

	for i := 0; i < 3; i++ {
		insertData = append(insertData, useData[i])
	}

	return insertData
}
