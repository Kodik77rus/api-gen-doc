package services

import (
	"io/ioutil"
	"net/http"
	url2 "net/url"
)

func HttpResponse(url string) (string, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return "", err
	}

	c := http.Client{}

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
