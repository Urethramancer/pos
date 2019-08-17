package database

import (
	"time"
)

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

// AddClient returns the client ID on success.
func (db *Invoices) AddClient(c *Client) (int64, error) {
	q := "INSERT INTO public.clients (company, email, phone, address, vatid) VALUES($1,$2,$3,$4,$5) RETURNING id;"
	var id int64
	err := db.QueryRow(q, c.Company, c.Email, c.Phone, c.Address, c.VATID).Scan(&id)
	return id, err
}

// UpdateClient changes the details of a client.
func (db *Invoices) UpdateClient(c *Client) error {
	q := "UPDATE public.clients SET company=$2, email=$3, phone=$4, address=$5, vatid=$6 WHERE id=$1;"
	_, err := db.Exec(q, c.ID, c.Company, c.Email, c.Phone, c.Address, c.VATID)
	return err
}

// GetClient returns one client by ID.
func (db *Invoices) GetClient(id int64) *Client {
	q := "SELECT id,company,email,phone,address,vatid FROM public.clients WHERE id=$1 LIMIT 1;"
	row := db.QueryRow(q, id)
	var c Client
	err := row.Scan(&c.ID, &c.Company, &c.Email, &c.Phone, &c.Address, &c.VATID)
	if err != nil {
		return nil
	}

	return &c
}

// GetAllClients returns all clients.
func (db *Invoices) GetAllClients() ([]*Client, error) {
	q := "SELECT id,company,email,phone,address,vatid FROM public.clients;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*Client
	for rows.Next() {
		var c Client
		err = rows.Scan(&c.ID, &c.Company, &c.Email, &c.Phone, &c.Address, &c.VATID)
		if err != nil {
			return nil, err
		}

		list = append(list, &c)
	}

	return list, nil
}

// GetClients returns all matching clients.
func (db *Invoices) GetClients(keyword string) ([]*Client, error) {
	q := "SELECT id,company,email,phone,address,vatid FROM public.clients WHERE company LIKE '%' || $1 || '%';"
	rows, err := db.Query(q, keyword)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*Client
	for rows.Next() {
		var c Client
		err = rows.Scan(&c.ID, &c.Company, &c.Email, &c.Phone, &c.Address, &c.VATID)
		if err != nil {
			return nil, err
		}

		list = append(list, &c)
	}

	return list, nil
}

func (db *Invoices) RemoveClient(keyword string) error {
	return nil
}
