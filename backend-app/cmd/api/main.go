package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// blank identifier
// To import a package solely for its side-effects (initialization), use the blank identifier as explicit package name:
// import _ "lib/math"

const version = "1.0.0"

// hold the app config
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// qn: when to capitalise ?
type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	// read port and env from the command line
	// read from flag
	// read from command line some int var and store in cfg.port
	// &(cfg.port) is address of type int: *int
	// called port in command line (flag name), default value, description
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	// change dsn from flag
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://kylekira@localhost/go_movies?sslmode=disable", "Postgres connection string")
	flag.Parse() // qn: why must parse?

	// qn: pipe?
	// write to os.stdout, prefix, logging properties: pipe together the integer
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// after main() then db.Close()
	defer db.Close()

	// get address of the application struct
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	// reference to http server
	// need a few members:

	srv := &http.Server{
		// Sprintf() format then return the string
		// qn: why not just cfg.port?; ans: becos need the colon
		Addr: fmt.Sprintf(":%d", cfg.port), // listen to localhost port 4000
		// what handler do you wanna use
		Handler: app.routes(),
		// if nothing in one minute, then stop
		IdleTimeout: time.Minute,
		// timeout when 10 sec pass
		ReadTimeout: 10 * time.Second,
		// timeout when 30 sec pass; too long to write then terminate
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	// start the server
	err = srv.ListenAndServe()

	// start the web server
	// format verb: https://pkg.go.dev/fmt
	// %d: integer value, base 10
	// second arg: no data
	// err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
	if err != nil {
		log.Println(err)
	}
}

// open the db
// take in config
func openDB(cfg config) (*sql.DB, error) {
	// return pointer DB struct
	// driver name, driver source name
	// tell connector which db driver to use
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// build context
	// get context that is always available, duration timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// delay the execution of the function until the nearby functions return
	defer cancel()

	// ping with context
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// db: connection tool
	return db, nil
}
