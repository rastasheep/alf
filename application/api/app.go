package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rastasheep/alf/migrations"
)

var (
	dbDataConnection  = flag.String("db-data-connection", "", "Connection to read only database")
	dbStoreConnection = flag.String("db-store-connection", "", "Connection to database for alf data")
	env               = flag.String("env", "development", "Application environment")
	port              = flag.String("port", "3000", "API server port")
	cacheSize         = flag.Int64("cache-size", 8, "Cache size in MB")
	pageSize          = flag.Int64("page-size", 20, "Number of items per page for API responses")

	maxSerial = 2147483647
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "", 0)

	s := NewServer(logger, *dbDataConnection, *dbStoreConnection, *cacheSize, *pageSize)
	defer func() {
		s.dbData.Close()
		s.dbStore.Close()
	}()

	logger.Println("migrating database")
	if err := migrations.Exec(s.dbStore); err != nil {
		log.Fatal(err)
	}

	s.Handle("/api/schema", s.schemaHandler)
	s.Handle("/api/executions", s.executionHandler)
	s.Handle("/api/results", s.resultHandler)
	s.Handle("/api/templates", s.templateHandler)

	logger.Printf("running server in %s mode on port %s", *env, *port)

	http.Handle("/", Adapt(s, Logger(logger)))

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
