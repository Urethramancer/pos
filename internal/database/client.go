package database

import "time"

// Client structure translated from the database.
type Client struct {
	// ID is autogenerated.
	ID int64
	// Company is the primary string reference.
	Company string
	// Email is the primary contact. Add to the contacts table for more.
	Email string
	// Phone if applicable.
	Phone string
	// Address to print on invoices.
	Address string
	// VATID or org. no.
	VATID string
	// Created timestamp.
	Created time.Time
}
