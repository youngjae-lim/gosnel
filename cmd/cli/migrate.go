package main

func doMigrate(arg2, arg3 string) error {
	// NOTE: won't be using this functionality eventually
	// dsn := getDSN()

	// exit gracefully when database type is not set or database config file does not exist
	checkForDB()

	tx, err := gos.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// run the migration command
	switch arg2 {
	case "up":
		err := gos.RunPopMigrations(tx)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := gos.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := gos.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := gos.PopMigrationReset(tx)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
