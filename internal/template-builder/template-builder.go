package templatebuilder

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
)

const (
	wordFolderName = "/word"
	pdfFolderName  = "/pdf"

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

func New(c config.TemplateBuilder, t Template) *TemplateBuilder {
	return &TemplateBuilder{
		Config:   c,
		Template: t,
	}
}

func (t *TemplateBuilder) BuildTemplate() error {
	_, err := t.createWordFile()
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateBuilder) createWordFile() (*os.File, error) {
	fileDir := t.generateDirName(wordFolderName)
	fileName := t.generateFileName(wordFolderName)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fileDir + "/" + fileName)
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

	return file, nil
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

func (t *TemplateBuilder) generateDirName(docType string) string {
	return fmt.Sprint(
		t.Config.TemplateFolder,
		"/",
		t.Template.TemplateName,
		docType,
		"/",
		t.Template.FolderId,
	)
}

func (t *TemplateBuilder) generateFileName(docType string) string {
	currentTime := time.Now()

	switch docType {
	case wordFolderName:
		return fmt.Sprint(
			currentTime.Format(timeFormat),
			".doc",
		)
	case pdfFolderName:
		return fmt.Sprint(
			currentTime.Format(timeFormat),
			".pdf",
		)
	default:
		return ""
	}
}
