package main

import (
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrations []*gormigrate.Migration

func main() {
	dbConnectionUrl, exists := os.LookupEnv("DB_CONNECTION_URL")
	if !exists {
		log.Fatal("DB_CONNECTION_URL env variable does not exist")
	}

	db, err := gorm.Open(postgres.Open(dbConnectionUrl))
	if err != nil {
		log.Fatalf("migration error: %v", err)
	}

	err = run(db, os.Args)
	if err != nil {
		log.Fatalf("migration error: %v", err)
	}
}

func run(db *gorm.DB, args []string) error {
	cmd := ""
	if len(args) > 1 {
		cmd = args[1]
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)

	switch cmd {
	case "migrate":
		if len(args) == 2 { // nolint:mnd
			err := m.Migrate()
			if err != nil {
				return err
			}
		} else if len(args) == 3 { // nolint:mnd
			toMigrationID := args[2]
			err := m.MigrateTo(toMigrationID)
			if err != nil {
				return err
			}
		}
		log.Print("migration migrate ok")
	case "rollback":
		if len(args) == 2 { // nolint:mnd
			err := m.RollbackLast()
			if err != nil {
				return err
			}
		} else if len(args) == 3 { // nolint:mnd
			toMigrationID := args[2]
			err := m.RollbackTo(toMigrationID)
			if err != nil {
				return err
			}
		}
		log.Print("migration rollback ok")
	default:
		log.Printf("migration run without any action")
		return nil
	}

	return nil
}

func addMigration(id string, migrate, rollback func(tx *gorm.DB) error) {
	m := &gormigrate.Migration{
		ID: id,
		Migrate: func(tx *gorm.DB) error {
			log.Printf("start migration %s\n", id)
			err := migrate(tx)
			if err != nil {
				return err
			}
			log.Printf("end migration %s\n", id)
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			log.Printf("start rollback %s\n", id)
			err := rollback(tx)
			if err != nil {
				return err
			}
			log.Printf("end rollback %s\n", id)
			return nil
		},
	}
	migrations = append(migrations, m)
}
