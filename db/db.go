package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

//DB ...
type DB struct {
	*sql.DB
}

var db *gorp.DbMap

//Init ...
func Init() *gorp.DbMap {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createArticlesTabel() string {
	return `
	CREATE TABLE IF NOT EXISTS public.articles
	(
	    id serial NOT NULL,
	    deleted boolean NOT NULL DEFAULT false,
	    subject character varying(160) COLLATE pg_catalog."default" NOT NULL,
	    body text COLLATE pg_catalog."default" NOT NULL,
	    summary character varying(300) COLLATE pg_catalog."default",
	    tags json,
	    pcc json,
	    created_by character varying COLLATE pg_catalog."default",
	    modified_by character varying COLLATE pg_catalog."default",
	    created_date timestamp with time zone,
	    modified_date timestamp with time zone,
	    attachments json,
			custom_fields json,
			tags_string text,
	    CONSTRAINT articles_pkey PRIMARY KEY (id)
	)
	WITH (
	    OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public.articles
	    OWNER to admin;
	`
}

func createTagsTabel() string {
	return `
	CREATE TABLE IF NOT EXISTS public.tags
	(
	    id serial NOT NULL,x
	    name text COLLATE pg_catalog."default" NOT NULL,
	    CONSTRAINT tags_pkey PRIMARY KEY (id)
	)
	WITH (
	    OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public.tags
	    OWNER to admin;
	`
}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	//Create Article table if not exists
	_, err = db.Exec(createArticlesTabel())
	if err != nil {
		log.Println(err)
	}
	//Create Tags table if not exists
	_, err = db.Exec(createTagsTabel())
	if err != nil {
		log.Println(err)
	}

	return dbmap, nil
}

func GetDB() *gorp.DbMap {
	return db
}
