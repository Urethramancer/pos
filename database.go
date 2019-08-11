package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Urethramancer/signor/log"
	_ "github.com/lib/pq"
)

func testDBExistence(cfg *Config) error {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	return db.Ping()
}

func testDBHost(cfg *Config) error {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	return db.Ping()
}

func ensureDBExists(cfg *Config) error {
	var err error
	// Open connection to host with database name and ping it.
	err = testDBExistence(cfg)
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

// Invoices is the name of our DB type to make it clear this isn't generic.
type Invoices struct {
	db *sql.DB
}

func OpenDB(host, port, user, password, dbname string, ssl bool) (*Invoices, error) {
	sslmode := ""
	if ssl {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	inv := &Invoices{}
	var err error
	inv.db, err = sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	// Test the connection to the host and the database.
	err = inv.db.Ping()
	if err != nil {
		return nil, err
	}

	return inv, nil
}
