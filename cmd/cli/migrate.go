package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// run the migration command
	switch arg2 {
	case "up":
		err := gos.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := gos.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := gos.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := gos.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = gos.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
