package main

import (
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

func doNew(appName string) {
	appName = strings.ToLower(appName)

	// sanitize the application name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[(len(exploded) - 1)]
	}

	log.Println("App name is", appName)

	// git clone the skeleton application
	color.Green("\tCloning repository...")
	_, err := git.PlainClone("./" + appName, false, &git.CloneOptions{
		URL: "git@github.com/youngjae-lim/gosnel-app.git",
		Progress: os.Stdout,
		Depth: 1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git directory

	// create a ready to go .env file

	// create a makefile

	// update the go.mod file

	// update existing .go files with correct name/imports

	// run do mod tidy in the project directory

}
