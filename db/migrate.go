package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	// go-lang migrate library - run database migrations
	_ "github.com/golang-migrate/migrate/source/file"
	// pq library - postgresql driver
	_ "github.com/lib/pq"

	"log"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/config"
	"strings"
)

const migrationsPath = "file://./db/migrations"

func Run(cfg config.Config) {
	db := initDatabaseFromConfig(cfg)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error while creating postgres driver - %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatalf("Error while initialising database instance - %v", err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("No migrations")
		return
	}

	if err != nil {
		log.Fatalf("Failed to run migrations - %v", err)
	}

	log.Println("Migrations run successfully")
}

func initDatabaseFromConfig(cfg config.Config) *sql.DB {
	conn := connectionString(cfg)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("Invalid database configuration")
	}
	if db.Ping() != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Printf("Connection to database successful - %s\n", cfg.DBName)
	return db
}

func connectionString(cfg config.Config) string {
	hostPort := strings.Split(cfg.DBAddr, ":")
	return fmt.Sprintf("dbname=%s user=%s password='%s' host=%s port=%s sslmode=disable",
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
		hostPort[0],
		hostPort[1],
	)
}
