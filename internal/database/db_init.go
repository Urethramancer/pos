package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Urethramancer/pos/internal/config"
	"github.com/Urethramancer/signor/log"
)

// TestDBExistence checks for the host and the database.
func TestDBExistence(cfg *config.Config) error {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	return db.Ping()
}

// TestDBHost only tests for the host.
func TestDBHost(cfg *config.Config) error {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	return db.Ping()
}

// EnsureDBExists creates the database and the correct tables.
func EnsureDBExists(cfg *config.Config) error {
	var err error
	// Open connection to host with database name and ping it.
	err = TestDBExistence(cfg)
	if err == nil {
		// No error means the database already exists. Anything else is a reason to try creating a new one.
		return nil
	}

	m := log.Default.Msg

	// Since it didn't exist, connect to the host itself with no database name.
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	// Create the database itself with basic settings.
	defer db.Close()
	m("No database. Creating.")
	q := strings.ReplaceAll(createdb, "{DBNAME}", cfg.DBName)
	q = strings.ReplaceAll(q, "{OWNER}", cfg.Username)
	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	// Now we can open a connection with the database name again.
	info = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err = sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	// At this point it's finally safe to create the tables, having them added to the correct schema.
	q = strings.ReplaceAll(createtables, "{DBNAME}", cfg.DBName)
	q = strings.ReplaceAll(q, "{OWNER}", cfg.Username)
	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	// And finally correct the sequence for the invoices.
	q = fmt.Sprintf("ALTER SEQUENCE invoices_id_seq RESTART WITH %d;", cfg.FirstInvoice)
	_, err = db.Exec(q)
	return err
}
