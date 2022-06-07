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
	req, _ := http.NewRequest("POST", convertWordToPdfUrl, bufBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.Do(req)
	if err != nil {
		return "", nil
	}

	if resp.StatusCode != 200 {
		return "", nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	return string(body), nil
}

func HttpResponse(url string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", nil
	}

	req.Header.Add("User-Agent", "hackerman")

	resp, err := c.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
