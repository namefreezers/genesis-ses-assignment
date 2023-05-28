package util

import (
	"fmt"
	"io"
	"net/http"
)

func Request_and_get_resp_body(api_url string) ([]byte, error) {
	resp, err := http.Get(api_url)
	if err != nil {
		return nil, fmt.Errorf("can't `get` url: %v. Error: %v", api_url, err.Error())
	}

	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body from `get` response: %v. Error: %v", api_url, err.Error())
	}

	return body_bytes, nil
}
