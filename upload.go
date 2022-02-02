package gosnel

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/youngjae-lim/gosnel/filesystems"
)

func (g *Gosnel) UploadFile(r *http.Request, destination, field string, fs filesystems.FS) error {
	fileName, err := g.getFileToUpload(r, field)
	if err != nil {
		g.ErrorLog.Println(err)
		return err
	}

	if fs != nil { // remote file system
		err = fs.Put(fileName, destination)
		if err != nil {
			g.ErrorLog.Println(err)
			return err
		}
	} else { // local file system
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			g.ErrorLog.Println(err)
			return err
		}
	}

	return nil
}

// getFileToUpload checks if the mime type of file to be uploaded is valid and returns a file path & name
// example: ./tmp/your_file.jpeg
func (g *Gosnel) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(10 << 20) // up to 20mb

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// NOTE: the DetectReader will move the reader offset that needes to be put back to the start of the file later.
	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	// move the reader offset to start of the file - i.e., move the reader offset at the start of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// only allowed mime types
	validMimeTypes := []string{
		"image/gif",
		"image/jpeg",
		"image/png",
		"application/pdf",
	}

	// check if the mime type  is valid
	if !inSlice(validMimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	// create a file to be uploaded in the ./tmp directory
	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy the file to the ./tmp directory
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

// inSlice checks wheter a value is in the slice and return true if it is.
func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
