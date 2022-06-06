package templatebuilder

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
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

func New(c config.TemplateBuilder) *TemplateBuilder {
	return &TemplateBuilder{
		Config: c,
	}
}

func (t *TemplateBuilder) BuildTemplate(template Template) error {
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

	file, err := os.Create(fileDir + fileName)
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
	var s scanner.Scanner
	var newTemplate string

	prepareUseDate := strings.Split(t.Template.InsertData.Use, ",")
	count := 1

	s.Init(strings.NewReader(*t.Template.Template))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if strings.Contains(txt, "_") {
			switch count {
			case 1:
				txt = strings.Replace(txt, "_", t.Template.InsertData.Text, 1)
				count++
			case 2:
				txt = strings.Replace(txt, "_", prepareUseDate[0], 1)
				count++
			case 3:
				txt = strings.Replace(txt, "_", prepareUseDate[1], 1)
				count++
			case 4:
				txt = strings.Replace(txt, "_", prepareUseDate[2], 1)
				count++
			}
		}
		newTemplate += txt
	}

	return &newTemplate
}

func (t *TemplateBuilder) generateDirName(docType string) string {
	return fmt.Sprint(
		t.Config.TemplateFolder,
		"/",
		t.Template.TemplateName,
		docType,
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
