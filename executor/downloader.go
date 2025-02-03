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

func ClearDownloadFolder() {
	dir, _ := os.UserConfigDir()
	ApplicationDirectoryX := fmt.Sprintf("%s\\%s", dir, "Cosmic")

	files, err := os.ReadDir(ApplicationDirectoryX)
	if err != nil {
		log.Printf("Failed to read directory: %v", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(ApplicationDirectoryX, file.Name())
		if file.IsDir() {
			if err := os.RemoveAll(filePath); err != nil {
				log.Printf("Failed to remove folder %s: %v", filePath, err)
			}
		}
	}
}

func extractZip(zipFile, extractDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %v", err)
	}
	defer r.Close()

	absExtractDir, err := filepath.Abs(extractDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for extraction directory: %v", err)
	}

	if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create extraction directory: %v", err)
	}

	for _, f := range r.File {
		filePath := filepath.Join(extractDir, f.Name)

		// Validate file path to prevent Zip Slip vulnerability
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for file: %v", err)
		}
		if !strings.HasPrefix(absFilePath, absExtractDir+string(filepath.Separator)) {
			return fmt.Errorf("invalid file path: %s escapes extraction directory", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory: %v", err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in ZIP: %v", err)
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create output file: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
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
