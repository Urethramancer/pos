package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Invoices is the name of our DB type to make it clear this isn't generic.
type Invoices struct {
	*sql.DB
}

// Open a database session.
func Open(host, port, user, password, dbname string, ssl bool) (*Invoices, error) {
	sslmode := ""
	if ssl {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	inv := &Invoices{}
	var err error
	inv.DB, err = sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	// Test the connection to the host and the database.
	err = inv.Ping()
	if err != nil {
		return nil, err
	}

	return inv, nil
}
