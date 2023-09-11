package utils

import (
	"io"
	"net/http"
	"strings"
)

func MakeHTTPRequest(method, endpoint string, requestBody []byte, header map[string][]string) ([]byte, int, error) {
	request, err := http.NewRequest(method, endpoint, strings.NewReader(string(requestBody)))
	if err != nil {
		return nil, 500, err
	}
	request.Header = header

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, 500, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}

	return respBody, resp.StatusCode, nil
}
