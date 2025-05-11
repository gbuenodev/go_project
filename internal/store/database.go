package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type DBConfig struct {
	Provider string
	Driver   string
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
	SSL      string
}

func Open(dbConfig *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
		dbConfig.SSL,
	)

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db: ping %w", err)
	}

	fmt.Printf(`----- Connected to Database -----
Provider: %s
Port: %d
DB: %s
User: %s
---------------------------------
`,
		dbConfig.Provider,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.User,
	)

	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
