package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	// migrations
	dbType := gos.DB.DbType

	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := gos.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := gos.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users, tokens, remember_tokens cascade;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run up migrations
	doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// copy files over
	err = copyFileFromTemplate("templates/data/user.go.txt", gos.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", gos.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy over middlware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", gos.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", gos.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("    - users, tokens, and remember_tokens migrations created and executed")
	color.Yellow("    - user and token models created")
	color.Yellow("    - auth middleware created")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models in data/models.go, and to add appropriate middleware to your routes!")

	return nil
}
