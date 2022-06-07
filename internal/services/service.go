package services

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	url2 "net/url"
)

const convertWordToPdfUrl = "http://localhost:3000/forms/libreoffice/convert"

var c http.Client

func MultipartResponse(writer *multipart.Writer, bufBody io.Reader) (string, error) {
	req, err := http.NewRequest("POST", convertWordToPdfUrl, bufBody)
	if err != nil {
		return "", converterError{
			err: err,
		}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.Do(req)
	if err != nil {
		return "", converterError{
			err: err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", convertServiceRespError{
			url:        convertWordToPdfUrl,
			statusCode: resp.StatusCode,
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", converterError{
			err: err,
		}
	}

	return string(body), nil
}

func HttpResponse(url string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", templateError{
			err: err,
		}
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", templateError{
			err: err,
		}
	}

	req.Header.Add("User-Agent", "hackerman")

	resp, err := c.Do(req)
	if err != nil {
		return "", templateError{
			err: err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", templateRespError{
			url:        url,
			statusCode: resp.StatusCode,
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", templateError{
			err: err,
		}
	}

	return string(body), nil
}
