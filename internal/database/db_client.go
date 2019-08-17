package database

// AddClient returns the client ID on success.
func (db *Invoices) AddClient(company, email, phone, address, vatid string) (int64, error) {
	q := "INSERT INTO public.clients (company, email, phone, address, vatid) VALUES($1,$2,$3,$4,$5) RETURNING id;"
	var id int64
	err := db.QueryRow(q, company, email, phone, address, vatid).Scan(&id)
	return id, err
}

// GetClient returns one client.
func (db *Invoices) GetClient(id int64) {

}

// GetClients returns all matching clients.
func (db *Invoices) GetClients(keyword string) {

}

func (db *Invoices) RemoveClient(keyword string) {

}
