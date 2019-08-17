package database

// AddClient returns the client ID on success.
func (db *Invoices) AddClient(company, email, phone, address, vatid string) (int64, error) {
	q := "INSERT INTO public.clients (company, email, phone, address, vatid) VALUES($1,$2,$3,$4,$5) RETURNING id;"
	var id int64
	err := db.QueryRow(q, company, email, phone, address, vatid).Scan(&id)
	return id, err
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

// GetClients returns all matching clients.
func (db *Invoices) GetClients(keyword string) []Client {
	var c []Client
	return c
}

// UpdateClient changes the details of a client.
func (db *Invoices) UpdateClient(companyname, email, phone, address, vatid string) error {
	return nil
}

func (db *Invoices) RemoveClient(keyword string) error {
	return nil
}
