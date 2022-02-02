package gosnel

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

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

func (g *Gosnel) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(10 << 20) // up to 20mb

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}
