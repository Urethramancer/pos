package main

import (
	"github.com/Urethramancer/signor/files"
)

// Config JSON.
type Config struct {
	// Name of the user (full).
	Name string `json:"name"`
	// Email of the user.
	Email string `json:"email"`
	// Company name on invoices.
	Company string `json:"company"`
	// Company ID on invoices.
	CompanyID string `json:"companyid"`
	// VAT percentage.
	VAT string `json:"vat"`
	// Address on invoices.
	Address string `json:"address"`
	// Host of Postgres DB.
	Host string `json:"host"`
	// Port of DB.
	Port string `json:"port"`
	// DBName is usually "invoices".
	DBName string `json:"dbname"`
	// Username to connect to DB.
	Username string `json:"username"`
	// Password of DB user.
	Password string `json:"password"`
	// FirstInvoice is the number to start with on DB creation.
	FirstInvoice int `json:"first_invoice"`
	// InvoicePrefix goes before any year prefix and the number.
	InvoicePrefix string `json:"invoice_prefix"`
	// YearPrefix prepends the invoice number with the year.
	YearPrefix bool `json:"year_prefix"`
}

func (cfg *Config) Load(fn string) error {
	return files.LoadJSON(fn, cfg)
}

func (cfg *Config) Save(fn string) error {
	return files.SaveJSON(fn, cfg)
}
