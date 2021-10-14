package gosnel

import "os"

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
