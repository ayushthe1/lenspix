package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

// Open will open a SQL connection with the provided Postgres database. Callers op Open need to ensure that the connection is eventually closed via the db.Close() method
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

// Not used anymore. The configs are loaded from the env variables
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		Database: "database",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

//  fmt.Printf is used to print formatted strings to the standard output, while fmt.Sprintf is used to format strings and store the result in a new string variable. Both functions use format strings and arguments, but they differ in their output destinations.
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Migrate(db *sql.DB, dir string) error {

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// The pressly/goose library uses global variables, which is why we can call a function like SetDialect and the changes from that function call are persisted later when we call the Up function

	// run the migrations
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {

	if dir == "" {
		dir = "."
	}

	goose.SetBaseFS(migrationsFS)
	defer func() {
		//  Ensure that we remove the FS on the off chance some other part of our app uses goose for migrations and doesn't want to use our FS.
		goose.SetBaseFS(nil)
	}()

	// We are able to use the existing Migrate function because goose uses global variables, so when we call Migrate it will use the base file system we just set. Then once our MigrateFS function is complete it will unset the base file system by passing in nil in the deferred code.
	return Migrate(db, dir)
}
