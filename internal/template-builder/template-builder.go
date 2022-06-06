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

	fmt.Printf(fileDir)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fileDir + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	newTemplate := t.replaceTemplateData()

	_, errr := file.WriteString(*newTemplate)
	if errr != nil {
		return nil, errr
	}

	return file, nil
}

func (t *TemplateBuilder) replaceTemplateData() *string {
	var replacedData string

	insertData := make([]string, 4, 4)
	useData := strings.Split(t.Template.InsertData.Use, ",")

	for i, v := range useData {
		if i == 0 {
			insertData[i] = t.Template.InsertData.Text
		}
		insertData[i+1] = v
	}

	for _, v := range insertData {
		replacedData = strings.Replace(*t.Template.Template, "_", v, 1)
	}

	return &replacedData
}

func (t *TemplateBuilder) generateDirName(docType string) string {
	return fmt.Sprint(
		t.Config.TemplateFolder, //../template
		"/",
		t.Template.TemplateName, //
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
