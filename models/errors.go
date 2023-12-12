package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

// These error variables are exported to other packages as these start with capital
var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	//  io.ReadSeeker so that we can reset the file after reading some of it.
	// we only need to check the first 512 bytes from the file
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes) // read the first 512 bytes
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	// reset the file
	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, t := range allowedTypes {
		if contentType == t {
			return nil
		}
	}
	return FileError{
		Issue: fmt.Sprintf("invalid content type: %v", contentType),
	}

}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
		}
	}

	return nil
}
