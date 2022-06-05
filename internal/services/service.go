package services

import (
	"io/ioutil"
	"net/http"
)

func HttpResponse(url string) (string, error) {
	c := http.Client{}

	req, err := http.NewRequest("GET", url, nil)

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
