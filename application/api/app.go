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

	"github.com/rastasheep/alf/migrations"
)

var (
	dbDataConnection  = flag.String("db-data-connection", "", "Connection to read only database")
	dbStoreConnection = flag.String("db-store-connection", "", "Connection to database for alf data")
	env               = flag.String("env", "development", "Application environment")

	maxSerial = 2147483647
)

type Server struct {
	dbData  *sqlx.DB
	dbStore *sqlx.DB
	logger  *log.Logger
	perPage int
}

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func NewDbConnection(config string) *sqlx.DB {
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

	dbData := NewDbConnection(*dbDataConnection)
	dbStore := NewDbConnection(*dbStoreConnection)
	defer func() {
		dbData.Close()
		dbStore.Close()
	}()

	logger.Println("Migrating database")
	if err := migrations.Exec(dbStore); err != nil {
		log.Fatal(err)
	}

	respOptions := RespondOptions()
	server := Server{
		dbData:  dbData,
		dbStore: dbStore,
		logger:  logger,
		perPage: 20,
	}
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

	router.Handle("/executions/{id:[0-9]+}", Adapt(
		http.HandlerFunc(server.getExecution),
		JSONResponse(respOptions),
	)).Methods("GET")

	router.Handle("/executions/{id:[0-9]+}", Adapt(
		http.HandlerFunc(server.deleteExecution),
		JSONResponse(respOptions),
	)).Methods("DELETE")

	logger.Printf("Running api server in %s mode\n", *env)

	log.Fatal(http.ListenAndServe(":3000", Adapt(router, Logger(logger))))
}
