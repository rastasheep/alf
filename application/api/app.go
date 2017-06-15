package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbConnection = flag.String("db-connection", "", "Connection to read only database")
	env          = flag.String("env", "development", "Application environment")
)

type Server struct {
	db     *sqlx.DB
	logger *log.Logger
}

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func NewDB(config string) *sqlx.DB {
	log.Printf("Connecting to postgres: %s", config)
	db, _ := sqlx.Open("postgres", config)

	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	return db
}

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "[request] ", 0)

	db := NewDB(*dbConnection)
	defer db.Close()

	respOptions := RespondOptions()
	server := Server{db, logger}
	router := mux.NewRouter()

	router.Handle("/schema", Adapt(
		http.HandlerFunc(server.getSchema),
		JSONResponse(respOptions),
	)).Methods("GET")

	router.Handle("/executions", Adapt(
		http.HandlerFunc(server.listExecutions),
		JSONResponse(respOptions),
	)).Methods("GET")

	router.Handle("/executions", Adapt(
		http.HandlerFunc(server.createExecution),
		JSONResponse(respOptions),
	)).Methods("POST")

	logger.Printf("Running api server in %s mode\n", *env)

	log.Fatal(http.ListenAndServe(":3000", Adapt(router, Logger(logger))))
}
