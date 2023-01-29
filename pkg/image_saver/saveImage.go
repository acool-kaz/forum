package image_saver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidImage = errors.New("invalid image")

func SaveImages(ctx context.Context, savePath, dir string, fileHeader *multipart.FileHeader) (string, error) {
	savePath += dir

	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("image saver: save images: %w", err)
	}

	if !strings.Contains(fileHeader.Header["Content-Type"][0], "image") {
		return "", fmt.Errorf("image saver: save images: %w", ErrInvalidImage)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("image saver: save images: %w", err)
	}
	defer file.Close()

	temp := strings.Split(fileHeader.Filename, ".")
	fileType := temp[len(temp)-1]

	if isInvalidImageType(fileType) {
		return "", fmt.Errorf("image saver: save images: %w", ErrInvalidImage)
	}

	fileName := uuid.NewString()
	savePath += "/" + fileName + "." + fileType

	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("image saver: save images: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("image saver: save images: %w", err)
	}

	return savePath[1:], nil
}

func isInvalidImageType(imageType string) bool {
	validImageType := []string{"jpeg", "jpg", "png"}
	for _, t := range validImageType {
		if t == imageType {
			return false
		}
	}
	return true
}
