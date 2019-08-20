package database

import "time"

// Task structure from the database.
type Task struct {
	// ID from database.
	ID int64
	// Name of task.
	Name string
	// Description of work involved.
	Description string
	// Hours worked.
	Hours float64
	// Cost multiplier.
	Cost float64
	// Created timestamp.
	Created time.Time
	// Start is the last time the task was started.
	Start time.Time
	// Stop is before Start if the task is running.
	Stop time.Time
	// Job reference.
	Job int64
}

// AddTask returns the task ID on success.
func (db *Invoices) AddTask(t *Task) (int64, error) {
	q := "INSERT INTO public.tasks (name,description,cost,job) VALUES($1,$2,$3,$4) RETURNING id;"
	var id int64
	err := db.QueryRow(q, t.Name, t.Description, t.Cost, t.Job).Scan(&id)
	return id, err
}

// UpdateTask changes the details of a task.
func (db *Invoices) UpdateTask(t *Task) error {
	q := "UPDATE public.tasks SET client=$2, currency=$3 WHERE id=$1;"
	_, err := db.Exec(q, t.ID, t.Name, t.Description, t.Cost, t.Job)
	return err
}

// UpdateTaskHours changes the amount worked.
func (db *Invoices) UpdateTaskHours(id int64, h float64) error {
	q := "UPDATE public.tasks SET hours=$2 WHERE id=$1;"
	_, err := db.Exec(q, id, h)
	return err
}

// UpdateTaskCost changes the cost per hour.
func (db *Invoices) UpdateTaskCost(id int64, c float64) error {
	q := "UPDATE public.tasks SET cost=$2 WHERE id=$1;"
	_, err := db.Exec(q, id, c)
	return err
}

// RemoveTask from database.
func (db *Invoices) RemoveTask(id int64) error {
	q := "DELETE FROM public.tasks WHERE id=$1;"
	_, err := db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
