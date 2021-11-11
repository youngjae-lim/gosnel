package gosnel

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func (g *Gosnel) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
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
