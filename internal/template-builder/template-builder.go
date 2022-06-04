package templatebuilder

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type Template struct {
	FolderId   int
	Template   *string
	InsertData struct {
		Text string
		Use  string
	}
}

type TemplateBuilder struct {
	OutPutFolder string
}

func New(outPutFolder string) *TemplateBuilder {
	return &TemplateBuilder{
		OutPutFolder: outPutFolder,
	}
}

func NewTemplate(id int, template *string, text string, use string) *Template {
	var data = struct {
		Text string
		Use  string
	}{
		Text: text,
		Use:  use,
	}

	return &Template{
		FolderId:   id,
		Template:   template,
		InsertData: data,
	}
}

func (t *TemplateBuilder) BuildTemplate(template Template) error {
	filePath := t.OutPutFolder + strconv.Itoa(template.FolderId)

	replacer := strings.NewReplacer(
		"Text", template.InsertData.Text,
		"Use", template.InsertData.Use,
	)

	replacesedData := replacer.Replace(*template.Template)

	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(generateFileName(filePath))
	if err != nil {
		return err
	}

	if _, err := file.WriteString(replacesedData); err != nil {
		return err
	}

	return nil
}

func generateFileName(filePath string) string {
	currentTime := time.Now()

	return fmt.Sprint(
		filePath,
		"/",
		strings.Replace(currentTime.Format(timeFormat), " ", "_", 1),
		".doc",
	)
}
