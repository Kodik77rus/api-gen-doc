package services

import "fmt"

type templateRespError struct {
	url        string
	statusCode int
}

type convertServiceRespError struct {
	url        string
	statusCode int
}

type converterError struct {
	err error
}

type templateError struct {
	err error
}

func (e converterError) Error() string {
	return fmt.Sprintf("failed converter err: %v", e.err)
}

func (e templateError) Error() string {
	return fmt.Sprintf("get template err: %v", e.err)
}

func (e templateRespError) Error() string {
	return fmt.Sprintf("failed get template status code from %v = %d", e.url, e.statusCode)
}

func (e convertServiceRespError) Error() string {
	return fmt.Sprintf("failed convert word file status code from %v = %d", e.url, e.statusCode)
}
