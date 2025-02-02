package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var githubURL = "https://raw.githubusercontent.com/Axion-Security/Cosmic/main/Registery/"
var localMode = false

func FetchTools(url string) (map[string]Application, error) {
	if localMode { // For testing purposes when adding new tools
		return fetchLocalTools(url)
	}

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

func fetchLocalTools(file string) (map[string]Application, error) {
	decodedFile, err := url.QueryUnescape(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file URL: %v", err)
	}
	filePath := "Registery/" + decodedFile
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var tools map[string]Application
	err = json.Unmarshal(fileContent, &tools)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return tools, nil
}
