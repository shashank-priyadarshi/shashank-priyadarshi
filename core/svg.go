package core

import (
	"fmt"
	"net/http"
	"os"
)

func FetchStatsSVG(srcs map[string]string) (err error) {
	for key, src := range srcs {
		var body []byte
		body, err = httpClient(http.MethodGet, src, "", nil, nil)
		if err != nil {
			logger.Sugar().Errorf("error making http request to %s while fetching %s: %s\n", src, key, err.Error())
			return
		}

		if err = os.WriteFile(fmt.Sprintf("./assets/images/%s.svg", key), body, 0644); err != nil {
			logger.Sugar().Errorf("error writing stats svg file: %s\n", err.Error())
			return
		}
	}

	return
}
