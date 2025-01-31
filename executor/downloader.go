package executor

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var ApplicationDirectory string

func init() {
	dir, _ := os.UserConfigDir()
	currentTime := time.Now().Format("15-04-05")
	ApplicationDirectory = fmt.Sprintf("%s\\%s\\%s", dir, "Cosmic", currentTime)

	if _, err := os.Stat(ApplicationDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(ApplicationDirectory, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func extractZip(zipFile, extractDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %v", err)
	}
	defer r.Close()

	if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create extraction directory: %v", err)
	}

	for _, f := range r.File {
		filePath := filepath.Join(extractDir, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in ZIP: %v", err)
		}
		defer rc.Close()

		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to write file content: %v", err)
		}
	}

	return nil
}

func DownloadFile(isCompressed bool, downloadUrl string) (string, bool) {
	u, err := url.Parse(downloadUrl)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		return "", false
	}

	resp, err := http.Get(downloadUrl)
	if err != nil {
		log.Printf("Failed to download file: %v", err)
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download file: server returned %s", resp.Status)
		return "", false
	}

	originalFilename := path.Base(u.Path)
	if originalFilename == "" || originalFilename == "." || originalFilename == "/" {
		originalFilename = "downloaded_file"
	}

	var timestamp = time.Now().Format("15-04-05")
	outputFilename := timestamp + "_" + originalFilename

	outputFilePath := filepath.Join(ApplicationDirectory, outputFilename)

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", outputFilePath, err)
		return "", false
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
		return "", false
	}

	if isCompressed {
		extractDir := strings.TrimSuffix(outputFilePath, filepath.Ext(outputFilePath))
		if err := extractZip(outputFilePath, extractDir); err != nil {
			log.Printf("Failed to extract ZIP file: %v", err)
			return "", false
		}

		info, err := os.Stat(extractDir)
		if err != nil || !info.IsDir() {
			return extractDir, false
		}

		return extractDir, true
	}

	return outputFilePath, false
}
