package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Urethramancer/signor/log"
	_ "github.com/lib/pq"
)

const (
	createdb = `CREATE DATABASE {DBNAME}
	WITH OWNER = {OWNER}
	ENCODING = 'UTF8'
	LC_COLLATE = 'en_US.UTF-8'
	LC_CTYPE = 'en_US.UTF-8';
`
	createtables = `
SET search_path={DBNAME},public;

-- Clients are the core of your business, hopefully.
CREATE TABLE public.clients
(
	id serial NOT NULL,
	-- Contact person to address an invoice to. Leave blank if it's addressed to a company.
	contact character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT "",
	-- Company the bill is for. May be blank if contact is specified.
	company character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT "",
	-- Email of contact/accounting department at the receiving company.
	email character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT "",
	-- Physical address of the recipient.
	address text COLLATE pg_catalog."default" NOT NULL DEFAULT "",
	-- This is a government-issued ID number, where applicable.
	companyid character varying(30) COLLATE pg_catalog."default",
	CONSTRAINT clients_pkey PRIMARY KEY (id)
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.clients OWNER to {OWNER};

-- A job is the pre-imvoice representation of work both in progress and finished.
CREATE TABLE public.job
(
	id serial NOT NULL,
	-- Short name or description of the collection of tasks/items.
	name character varying(100) COLLATE pg_catalog."default",
	-- Client this is done for.
	client bigint NOT NULL,
	-- Timestamp when it was ordered. Not necessarily the same as when actual work began.
	started timestamp with time zone NOT NULL,
	-- Currency code for the currency it's charged in.
	currency character varying(10) COLLAGE pg_catalog."default",
	-- Total cost of all jobs/items so far.
	cost numeric(16,2) NOT NULL DEFAULT 0,
	CONSTRAINT job_pkey PRIMARY KEY (id),
	-- All jobs in progress are cancelled when a client is removed.
	CONSTRAINT job_client_fkey FOREIGN KEY (client)
		REFERENCES public.clients (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.job OWNER to {OWNER};

-- Tasks are the entries making up a service-oriented invoice (work order).
CREATE TABLE public.task
(
	id serial NOT NULL,
	-- Short name of a task, usually one or more keywords.
	name character varying(50) COLLATE pg_catalog."default" NOT NULL,
	-- In-depth description of task.
	description text COLLATE pg_catalog."default" NOT NULL,
	-- Hours on this task so far.
	hours double precision NOT NULL,
	-- The cost per hour is in the currency of the job.
	cost numeric(16,2) NOT NULL DEFAULT 0,
	-- Start and stop timestamps are used for real-time tracking.
    start timestamp with time zone NOT NULL,
    stop timestamp with time zone,
    done boolean NOT NULL,
    job bigint NOT NULL,
	CONSTRAINT task_pkey PRIMARY KEY (id),
	-- All tasks are deleted when the job they're attached to is cancelled.
    CONSTRAINT task_job_fkey FOREIGN KEY (job)
        REFERENCES public.job (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.task OWNER to {OWNER};

-- Finally the table with invoices, which relies on all the rest.
CREATE TABLE public.invoices
(
	-- Doubles as the invoice number.
	id serial NOT NULL,
	-- Currency code.
	currency character varying(10) COLLAGE pg_catalog."default",
	-- Total cost of items/tasks in the job.
	total numeric(16,2) NOT NULL DEFAULT 0,
	-- How much of the total is VAT, sales tax etc.
	vatamount numeric(16,2) NOT NULL DEFAULT 0,
	-- Actual VAT percentage charged.
	vat smallint NOT NULL DEFAULT 0,
	-- Date the invoice was created and sent.
	sent timestamp with time zone NOT NULL,
	-- Due by this date.
	due timestamp with time zone NOT NULL,
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.task OWNER to {OWNER};	
`
)

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
	// Open once with a database name to connect to.
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	// Ping will now check for the existence of the named DB.
	err = db.Ping()
	db.Close()
	m := log.Default.Msg
	if err == nil {
		m("Database %s exists.", cfg.DBName)
		return nil
	}

	// Since it didn't exist, connect to the host itself with no database name.
	info = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err = sql.Open("postgres", info)
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

	return nil
}
