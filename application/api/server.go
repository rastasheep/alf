package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	"github.com/rastasheep/alf/executions"
	"github.com/rastasheep/alf/results"
	"github.com/rastasheep/alf/schema"
	"github.com/rastasheep/alf/templates"
)

type Server struct {
	dbData           *sqlx.DB
	dbStore          *sqlx.DB
	router           *http.ServeMux
	executionHandler http.Handler
	schemaHandler    http.Handler
	resultHandler    http.Handler
	templateHandler  http.Handler
}

func NewServer(logger *log.Logger, dbDataCon, dbStoreCon string, cacheSize, pageSize int64) *Server {
	s := &Server{}
	s.dbData = newDbConnection(dbDataCon)
	s.dbStore = newDbConnection(dbStoreCon)

	executionStore := &executions.ExecutionStore{
		DbStore: s.dbStore,
		DbData:  s.dbData,
	}
	resultCache := results.NewResultCache(logger, cacheSize*100000, executionStore.Execute)

	s.router = http.NewServeMux()

	s.executionHandler = &executions.ExecutionHandler{
		Store:       executionStore,
		ResultCache: resultCache,
		Logger:      logger,
		PageSize:    int(pageSize),
	}

	s.schemaHandler = &schema.SchemaHandler{
		Logger: logger,
		Store:  &schema.SchemaStore{s.dbData},
	}
	s.resultHandler = &results.ResultHandler{
		Logger:      logger,
		ResultCache: resultCache,
	}
	s.templateHandler = &templates.TemplateHandler{
		Logger: logger,
		Store:  &templates.TemplateStore{s.dbStore},
	}

	return s
}

func (s *Server) Handle(path string, handler http.Handler) {
	s.router.Handle(path, http.StripPrefix(path, handler))
	s.router.Handle(path+"/", http.StripPrefix(path, handler))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newDbConnection(config string) *sqlx.DB {
	log.Printf("Connecting to postgres: %s", config)
	db, _ := sqlx.Open("postgres", config)

	err := db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	return db
}
