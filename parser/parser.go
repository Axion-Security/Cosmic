package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var githubURL = "https://raw.githubusercontent.com/Axion-Security/Cosmic/main/Registery/"

func FetchTools(url string) (map[string]Application, error) {
	resp, err := http.Get(githubURL + url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var tools map[string]Application
	err = json.Unmarshal(body, &tools)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return tools, nil
}
