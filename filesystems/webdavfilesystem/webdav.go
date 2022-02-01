package webdavfilesystem

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/studio-b12/gowebdav"
	"github.com/youngjae-lim/gosnel/filesystems"
)

type WebDAV struct {
	Host string
	User string
	Pass string
}

func (w *WebDAV) getCredentials() *gowebdav.Client {
	c := gowebdav.NewClient(w.Host, w.User, w.Pass)
	return c
}

func (w *WebDAV) Put(fileName, folder string) error {
	client := w.getCredentials()

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = client.WriteStream(fmt.Sprintf("%s/%s", folder, path.Base(fileName)), file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (w *WebDAV) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	client := w.getCredentials()

	files, err := client.ReadDir(prefix)
	if err != nil {
		return listing, err
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			b := float64(file.Size())
			kb := b / 1024
			mb := kb / 1024
			current := filesystems.Listing{
				LastModified: file.ModTime(),
				Key:          file.Name(),
				Size:         mb,
				IsDir:        file.IsDir(),
			}
			listing = append(listing, current)
		}
	}
	return listing, nil
}

func (w *WebDAV) Delete(itemsToDelete []string) bool {
	return true
}

func (w *WebDAV) Get(destination string, items ...string) error {
	return nil
}
