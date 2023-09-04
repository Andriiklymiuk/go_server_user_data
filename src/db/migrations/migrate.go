package migrations

import (
	"flag"
	"fmt"
	"os"

	"github.com/andriiklymiuk/go_server_user_data/v2/src/utils"

	"github.com/fatih/color"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func CheckMigrationFlags(envConfig *utils.EnvConfig) bool {
	migrateUpCmd := flag.Bool("migrateUp", false, "Set to true to run db migrations up")
	migrateDownCmd := flag.Bool("migrateDown", false, "Set to true to run db migrations down")
	flag.Parse()

	if *migrateUpCmd {
		MigrateUpDb(envConfig)
		return true
	}

	if *migrateDownCmd {
		MigrateDownDb(envConfig)
		return true
	}

	return false
}

func setupMigration(envConfig *utils.EnvConfig) (*migrate.Migrate, error) {
	migrationPath := "file://src/db/migrations"

	var sslMode string
	if envConfig.DbHost == "localhost" || envConfig.DbHost == "127.0.0.1" {
		sslMode = "?sslmode=disable"
	}
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s%s",
		envConfig.DbUser,
		envConfig.DbPassword,
		envConfig.DbHost,
		envConfig.DbPort,
		envConfig.DbName,
		sslMode,
	)

	return migrate.New(migrationPath, databaseUrl)
}

func MigrateUpDb(envConfig *utils.EnvConfig) {
	m, err := setupMigration(envConfig)
	if err != nil {
		color.Red("Migration failed: %v", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			color.Yellow("No new migrations")
		} else {
			color.Red("Migration up failed: %v", err)
			os.Exit(1)
		}
	}

	color.Green("Migration successful")
}

func MigrateDownDb(envConfig *utils.EnvConfig) {
	m, err := setupMigration(envConfig)
	if err != nil {
		color.Red("Migration failed: %v", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		if err == migrate.ErrNoChange {
			color.Yellow("No migrations to undo")
		} else {
			color.Red("Migration down failed: %v", err)
			os.Exit(1)
		}
	}

	color.Green("Migration rollback successful")
}
