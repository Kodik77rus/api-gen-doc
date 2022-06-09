package templatebuilder

import "fmt"

type createWordFileErr struct {
	err error
}

type convertWordToPdf struct {
	err error
}

type savePdf struct {
	err error
}

func (e createWordFileErr) Error() string {
	return fmt.Sprintf("failed create word file: %v", e.err)
}

func (e convertWordToPdf) Error() string {
	return fmt.Sprintf("failed conver word to pdf: %v", e.err)
}

func (e savePdf) Error() string {
	return fmt.Sprintf("failed save pdf: %v", e.err)
}
