package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func httpClient(method, url, auth string, body []byte, headers map[string]string) (resBody []byte, err error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return
	}

	switch {
	case len(auth) != 0:
		req.Header.Add("Authorization", auth)
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	response, err := client.Do(req)
	if err != nil {
		logger.Sugar().Info(err)
		return
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("status code: %v", response.StatusCode)
		return
	}

	resBody, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return nil, response.Body.Close()
}
