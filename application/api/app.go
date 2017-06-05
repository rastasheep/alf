package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"regexp"

	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type execution struct {
	ID    int    `json:"id"`
	Query string `json:"query"`
	Name  string `json:"name"`
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

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is the RESTful api")
}

func createExecution(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var e execution
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	re := regexp.MustCompile("(?i)(SET.*TRANSACTION)|(SET.*SESSION.*CHARACTERISTICS)")
	matched := re.MatchString(e.Query)
	if matched {
		log.Printf("blocked execution of query: %v", e.Query)
		respondWithError(w, http.StatusBadRequest, "pq: you are not allowed to change the characteristics of the current transaction")
		return
	}

	log.Printf("executing query: %v", e.Query)

	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = tx.Exec("SET TRANSACTION READ ONLY;")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	rows, err := tx.Queryx(e.Query)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		entry := make(map[string]interface{})

		err := rows.MapScan(entry)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		data = append(data, entry)
	}
	tx.Commit()

	respondWithJSON(w, http.StatusCreated, data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	db = NewDB(os.Getenv("DB_CONNECTION"))
	defer db.Close()

	router := httprouter.New()
	router.GET("/", indexHandler)
	router.POST("/executions", createExecution)

	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":3000", router)
}
