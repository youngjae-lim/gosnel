package gosnel

import (
	"fmt"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Gosnel struct {
	AppName string
	Debug   bool
	Version string
}

func (g *Gosnel) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := g.Init(pathConfig)
	if err != nil {
		return err
	}

	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	return nil
}

func (g *Gosnel) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create a folder if it doesn't exist
		err := g.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (g *Gosnel) checkDotEnv(path string) error {
	err := g.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}
