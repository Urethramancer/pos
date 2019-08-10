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

CREATE TABLE public.clients
(
    id serial NOT NULL,
    contact character varying(100) COLLATE pg_catalog."default" NOT NULL,
    company character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    address text COLLATE pg_catalog."default" NOT NULL,
    companyid character varying(30) COLLATE pg_catalog."default",
    CONSTRAINT clients_pkey PRIMARY KEY (id)
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.clients OWNER to {OWNER};

CREATE TABLE public.job
(
    id serial NOT NULL,
    name character varying(100) COLLATE pg_catalog."default",
    client bigint NOT NULL,
    started timestamp with time zone NOT NULL,
    CONSTRAINT job_pkey PRIMARY KEY (id),
    CONSTRAINT job_client_fkey FOREIGN KEY (client)
        REFERENCES public.clients (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
) WITH (OIDS = FALSE) TABLESPACE pg_default;
ALTER TABLE public.job OWNER to {OWNER};

CREATE TABLE public.task
(
    id serial NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    hours double precision NOT NULL,
    start timestamp with time zone NOT NULL,
    "end" timestamp with time zone,
    done boolean NOT NULL,
    job bigint NOT NULL,
    CONSTRAINT task_pkey PRIMARY KEY (id),
    CONSTRAINT task_job_fkey FOREIGN KEY (job)
        REFERENCES public.job (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
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
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}

	err = db.Ping()
	db.Close()
	m := log.Default.Msg
	if err == nil {
		m("Database %s exists.", cfg.DBName)
		return nil
	}

	info = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	db, err = sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	m("No database. Creating.")
	q := strings.ReplaceAll(createdb, "{DBNAME}", cfg.DBName)
	q = strings.ReplaceAll(q, "{OWNER}", cfg.Username)
	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	info = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err = sql.Open("postgres", info)
	if err != nil {
		return err
	}

	defer db.Close()
	q = strings.ReplaceAll(createtables, "{DBNAME}", cfg.DBName)
	q = strings.ReplaceAll(q, "{OWNER}", cfg.Username)
	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
