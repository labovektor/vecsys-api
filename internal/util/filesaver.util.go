package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func FileSaver(fh *multipart.FileHeader, name string, saveTo ...string) (string, error) {
	// Open the uploaded file
	srcFile, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("could not open uploaded file: %v", err)
	}
	defer srcFile.Close()

	destDir := "./__public/"
	url := "/public/"
	if len(saveTo) > 0 {
		destDir += saveTo[0]
		url += saveTo[0]
	}

	// Ensure the destination directory exists
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("could not create destination directory: %v", err)
		}
	}

	// Generate nama file baru
	newFileName := generateNewFilename(name, fh.Filename)

	// Create the destination file path and url
	destPath := filepath.Join(destDir, newFileName)
	urlPath := filepath.Join(url, newFileName)

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("could not create destination file: %v", err)
	}
	defer destFile.Close()

	// Copy the uploaded file's content to the destination file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("could not copy file: %v", err)
	}

	return fmt.Sprintf("%s?t=%d", urlPath, time.Now().Unix()), nil
}

func generateNewFilename(userID string, filename string) string {
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", userID, ext)
	return newFilename
}
