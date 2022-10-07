package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ftxt-3-3/candle"
	"ftxt-3-3/flag"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

func main() {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"flag": &memdb.TableSchema{
				Name: "flag",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Flag"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	candleHandler := candle.NewCandleHandler()
	flagHandler := flag.NewFlagHandler(db)

	r := mux.NewRouter()
	r.HandleFunc("/candle", candleHandler.GetCandle).Methods("GET")
	r.HandleFunc("/flag", flagHandler.PutFlag).Methods("PUT")
	r.HandleFunc("/flag", flagHandler.GetFlag).Methods("Get")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
