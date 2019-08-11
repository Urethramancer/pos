package main

const (
	createdb = `CREATE DATABASE {DBNAME}
	WITH OWNER = {OWNER}
	ENCODING = 'UTF8'
	LC_COLLATE = 'en_US.UTF-8'
	LC_CTYPE = 'en_US.UTF-8';
`
	createtables = `SET search_path={DBNAME},public;

-- We'll trigger creation timestamp setting in a few places.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	NEW.created = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Clients are the core of your business, hopefully.
CREATE TABLE public.clients
(
	id serial NOT NULL,
	-- Company the bill is for. May be a person for sole proprietorships.
	company character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Email of contact/accounting department at the receiving company.
	email character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Default phone number. Add another via contacts.
	phone character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Physical address of the recipient.
	address text COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- This is a government-issued ID number, where applicable.
	companyid character varying(30) COLLATE pg_catalog."default",
	created timestamp with time zone,
	CONSTRAINT clients_pkey PRIMARY KEY (id),
	CONSTRAINT company_unique UNIQUE (company)
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.clients OWNER to {OWNER};

CREATE TRIGGER set_clients_timestamp
	BEFORE INSERT ON public.clients
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- Contacts for companies. Sometimes a company has more than one.
CREATE TABLE public.contacts
(
	id serial NOT NULL,
	-- Name of the person.
	name character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- E-mail they usually respond to.
	email character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Phone number, if applicable.
	phone character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Client company they work at/for.
	client bigint NOT NULL,
	-- Time this contact was added.
	created timestamp with time zone,
	CONSTRAINT contact_client_fkey FOREIGN KEY (client)
	-- All contacts disappear when a client is removed.
	REFERENCES public.clients (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.contacts OWNER to {OWNER};

CREATE TRIGGER set_contacts_timestamp
	BEFORE INSERT ON public.contacts
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


-- A job is the pre-imvoice representation of work both in progress and finished.
CREATE TABLE public.jobs
(
	id serial NOT NULL,
	-- Short name or description of the collection of tasks/items.
	name character varying(100) COLLATE pg_catalog."default",
	-- Client this is done for.
	client bigint NOT NULL,
	-- Timestamp when it was ordered. Not necessarily the same as when actual work began.
	created timestamp with time zone,
	-- Currency code for the currency it's charged in.
	currency character varying(10) COLLATE pg_catalog."default",
	-- Total cost of all jobs/items so far.
	cost numeric(16,2) NOT NULL DEFAULT 0,
	CONSTRAINT job_pkey PRIMARY KEY (id),
	-- All jobs in progress are cancelled when a client is removed.
	CONSTRAINT job_client_fkey FOREIGN KEY (client)
		REFERENCES public.clients (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.jobs OWNER to {OWNER};

CREATE TRIGGER set_jobs_timestamp
	BEFORE INSERT ON public.jobs
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- Tasks are the entries making up a service-oriented invoice (work order).
CREATE TABLE public.tasks
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
	-- Created is the time this task was added, but it's not used for hours tracking.
	created timestamp with time zone,
	-- Start and stop timestamps are used for real-time tracking.
	start timestamp with time zone,
	stop timestamp with time zone,
	done boolean NOT NULL DEFAULT false,
	job bigint NOT NULL,
	CONSTRAINT task_pkey PRIMARY KEY (id),
	-- All tasks are deleted when the job they're attached to is cancelled.
	CONSTRAINT task_job_fkey FOREIGN KEY (job)
		REFERENCES public.jobs (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.tasks OWNER to {OWNER};

CREATE TRIGGER set_tasks_timestamp
	BEFORE INSERT ON public.tasks
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- Finally the table with invoices, which relies on all the rest.
CREATE TABLE public.invoices
(
	-- Doubles as the invoice number.
	id serial NOT NULL,
	-- Currency code.
	currency character varying(10) COLLATE pg_catalog."default",
	-- Total cost of items/tasks in the job.
	total numeric(16,2) NOT NULL DEFAULT 0,
	-- How much of the total is VAT, sales tax etc.
	vatamount numeric(16,2) NOT NULL DEFAULT 0,
	-- Actual VAT percentage charged.
	vat smallint NOT NULL DEFAULT 0,
	-- Date the invoice was created.
	created timestamp with time zone,
	-- Due by this date.
	due timestamp with time zone
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.invoices OWNER to {OWNER};

CREATE TRIGGER set_invoices_timestamp
	BEFORE INSERT ON public.invoices
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
`
)
