package templatebuilder

import (
	"encoding/xml"
	"github.com/Kodik77rus/api-gen-doc/internal/config"
)

const (
	wordFolderName = "word"
	pdfFolderName  = "pdf"
	timeFormat     = "2006-01-02 15:04:05"
)

type TemplateBuilder struct {
	Config config.TemplateBuilder
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

type parsedXml struct {
	textField string `xml:"asd"`
}

func (t *TemplateBuilder) BuildTemplate(template Template) error {
	var a parsedXml

	if err := xml.Unmarshal([]byte(*template.Template), &a); err != nil {
		return err
	}

	//err := os.MkdirAll(filePath, os.ModePerm)
	//if err != nil {
	//	return err
	//}
	//
	//file, err := os.Create(generateFileName(filePath))
	//if err != nil {
	//	return err
	//}
	//

	return nil
}

//func (t *TemplateBuilder) generateFileName(docType string) string {
//	return fmt.Sprint()
//}
