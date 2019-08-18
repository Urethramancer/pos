package database

import "time"

// Contact structure from the database.
type Contact struct {
	// ID in the database.
	ID int64
	// Name of person. Non-unique to allow multiple e-mails or phone numbers.
	Name string
	// Email to contact them at. May be blank.
	Email string
	// Phone number to contact them at. May be blank.
	Phone string
	// Client they work for.
	Client int64
	// Created timestamp.
	Created time.Time
}

// AddContact returns the contact ID on success.
func (db *Invoices) AddContact(c *Contact) (int64, error) {
	q := "INSERT INTO public.contacts (name,email,phone,client) VALUES($1,$2,$3,$4) RETURNING id;"
	var id int64
	err := db.QueryRow(q, c.Name, c.Email, c.Phone, c.Client).Scan(&id)
	return id, err
}

// UpdateContact changes the details of a contact.
func (db *Invoices) UpdateContact(c *Contact) error {
	q := "UPDATE public.contacts SET name=$2, email=$3, phone=$4, client=$5 WHERE id=$1;"
	_, err := db.Exec(q, c.ID, c.Name, c.Email, c.Phone, c.Client)
	return err
}

// GetContact returns one contact by ID.
func (db *Invoices) GetContact(id int64) *Contact {
	q := "SELECT id,name,email,phone,client,created FROM public.contacts WHERE id=$1 LIMIT 1;"
	row := db.QueryRow(q, id)
	var c Contact
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Client, &c.Created)
	if err != nil {
		return nil
	}

	return &c
}

// GetAllContacts returns all contacts.
func (db *Invoices) GetAllContacts() ([]*Contact, error) {
	q := "SELECT id,name,email,phone,client,created FROM public.contacts;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*Contact
	for rows.Next() {
		var c Contact
		err = rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Client, &c.Created)
		if err != nil {
			return nil, err
		}

		list = append(list, &c)
	}

	return list, nil
}

// GetContacts returns all matching contacts.
func (db *Invoices) GetContacts(keyword string) ([]*Contact, error) {
	q := "SELECT id,name,email,phone,client,created FROM public.contacts WHERE name LIKE '%' || $1 || '%';"
	rows, err := db.Query(q, keyword)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*Contact
	for rows.Next() {
		var c Contact
		err = rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Client, &c.Created)
		if err != nil {
			return nil, err
		}

		list = append(list, &c)
	}

	return list, nil
}

// RemoveContact from database.
func (db *Invoices) RemoveContact(id int64) error {
	q := "DELETE FROM public.contacts WHERE id=$1;"
	_, err := db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
