package database

import "time"

// Job structure from the database.
type Job struct {
	// ID from database.
	ID int64
	// Client this job is for.
	Client int64
	// Currency to charge in.
	Currency string
	// Cost total of all tasks so far in that currency.
	Cost float64
	// Created timestamp.
	Created time.Time
}

// AddJob returns the job ID on success.
func (db *Invoices) AddJob(j *Job) (int64, error) {
	q := "INSERT INTO public.jobs (client,currency) VALUES($1,$2) RETURNING id;"
	var id int64
	err := db.QueryRow(q, j.Client, j.Currency).Scan(&id)
	return id, err
}

// UpdateJob changes the details of a job.
func (db *Invoices) UpdateJob(j *Job) error {
	q := "UPDATE public.jobs SET client=$2, currency=$3 WHERE id=$1;"
	_, err := db.Exec(q, j.ID, j.Client, j.Currency, j.Cost)
	return err
}

// GetJob returns one job by ID.
func (db *Invoices) GetJob(id int64) *Job {
	q := "SELECT id,client,currency,cost,created FROM public.jobs WHERE id=$1 LIMIT 1;"
	row := db.QueryRow(q, id)
	var j Job
	err := row.Scan(&j.ID, &j.Client, &j.Currency, &j.Cost, &j.Created)
	if err != nil {
		return nil
	}

	return &j
}

// GetJobsFor returns all jobs by client ID.
func (db *Invoices) GetJobsFor(id int64) *Job {
	q := "SELECT id,client,currency,cost,created FROM public.jobs WHERE client=$1;"
	row := db.QueryRow(q, id)
	var j Job
	err := row.Scan(&j.ID, &j.Client, &j.Currency, &j.Cost, &j.Created)
	if err != nil {
		return nil
	}

	return &j
}

// GetAllJobs returns all jobs, sorted by client.
func (db *Invoices) GetAllJobs() ([]*Job, error) {
	q := "SELECT id,client,currency,cost,created FROM public.jobs ORDER BY client ASC;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*Job
	var j Job
	for rows.Next() {
		err := rows.Scan(&j.ID, &j.Client, &j.Currency, &j.Cost, &j.Created)
		if err != nil {
			return nil, err
		}

		list = append(list, &j)
	}

	return list, nil
}

// RemoveJob from database.
func (db *Invoices) RemoveJob(id int64) error {
	q := "DELETE FROM public.jobs WHERE id=$1;"
	_, err := db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
