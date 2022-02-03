package gosnel

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (g *Gosnel) popConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (g *Gosnel) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = g.RootPath + "/migrations"

	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gosnel) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	return nil
}

func (g *Gosnel) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}
	return nil
}

func (g *Gosnel) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}
	return nil
}

func (g *Gosnel) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return nil
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}
	return nil
}
