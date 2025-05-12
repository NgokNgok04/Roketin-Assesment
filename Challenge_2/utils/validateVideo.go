package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const MAXVIDEOSIZE = 100 * 1024 * 1024
func ValidateVideo(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MAXVIDEOSIZE {return errors.New("video file exceeds 100MB limit")}
	
	file, err := fileHeader.Open()
	if err != nil {return errors.New("failed to open uplaoded file")}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {return errors.New("failed to read uploaded file")}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "video/") {return errors.New("only video files are allowed")}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtension := map[string]bool {
		".mp4" : true,
		".mkv" : true,
		".webm": true,
	}

	if !allowedExtension[ext] {return errors.New("unsupported video file extension")}

	return nil
}