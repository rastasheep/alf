package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/rastasheep/alf/execution"
	"github.com/rastasheep/alf/migrations"
	"github.com/rastasheep/alf/schema"
)

var (
	dbDataConnection  = flag.String("db-data-connection", "", "Connection to read only database")
	dbStoreConnection = flag.String("db-store-connection", "", "Connection to database for alf data")
	env               = flag.String("env", "development", "Application environment")
	cacheSize         = flag.Int64("cache-size", 8, "Cache size in MB")

	maxSerial = 2147483647
)

type Server struct {
	dbData       *sqlx.DB
	dbStore      *sqlx.DB
	logger       *log.Logger
	resultsCache *ResultsCache
	perPage      int
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
	logger := log.New(os.Stdout, "", 0)

	dbData := NewDbConnection(*dbDataConnection)
	dbStore := NewDbConnection(*dbStoreConnection)
	defer func() {
		dbData.Close()
		dbStore.Close()
	}()

	logger.Println("migrating database")
	if err := migrations.Exec(dbStore); err != nil {
		log.Fatal(err)
	}

	s := Server{
		dbData:  dbData,
		dbStore: dbStore,
		logger:  logger,
		perPage: 20,
	}
	s.initCache(*cacheSize * 100000)
	r := http.NewServeMux()

	schemaHandler := schema.NewSchemaHandler(s.logger, s.dbData)
	r.Handle("/schema", http.StripPrefix("/schema", schemaHandler))
	r.Handle("/schema/", http.StripPrefix("/schema", schemaHandler))

	executionHandler := execution.NewExecutionHandler(s.logger, s.dbStore, s.perPage)
	r.Handle("/executions", http.StripPrefix("/executions", executionHandler))
	r.Handle("/executions/", http.StripPrefix("/executions", executionHandler))

	//	router.Handle("/results", http.HandlerFunc(server.listResults)).Methods("GET")

	logger.Printf("running server in %s mode\n", *env)

	http.Handle("/", Adapt(r, Logger(logger)))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
