package gosnel

import (
	"crypto/rand"
	"os"
)

const (
	randomString = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_+"
)

// RandomString generates a random string length n from values in the const randomString
func (g *Gosnel) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomString)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

// CreateDirIfNotExist creates a new directory if it does not exist
func (g *Gosnel) CreateDirIfNotExist(path string) error {
	// owner can read/write/execute, group/others can read/execute
	const mode = 0755

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFileIfNotExist creates a new file if it does not exist
func (g *Gosnel) CreateFileIfNotExist(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}
