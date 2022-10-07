package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ftxt-3-3/candle"
	"ftxt-3-3/flag"
	"ftxt-3-3/model"

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

	cmp, err := getCandleMap()
	if err != nil {
		fmt.Println("###", err)
	}

	candleHandler := candle.NewCandleHandler(&cmp)
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

func getCandleMap() (model.CandleMap, error) {
	readFile, err := os.Open("./order_books.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	const (
		layout = "2006-01-02 15:04:05"
	)

	jst, _ := time.LoadLocation("Asia/Tokyo")

	mp := make(model.CandleMap)
	for _, line := range fileLines {
		lineParts := strings.Split(line, ",")
		timePart := lineParts[0]
		indx := strings.Index(timePart, " +0900")
		timePart = timePart[0:indx]
		timeStamp, err := time.ParseInLocation(layout, timePart, jst)
		if err != nil {
			log.Fatal("##", err)
		}
		price, err := strconv.Atoi(lineParts[2])
		if err != nil {
			log.Fatal("##", err)
		}

		slc, ok := mp[lineParts[1]]
		if !ok {
			slc = make([]model.Candle, 0)
			slc = append(slc, model.Candle{
				Time:  timeStamp,
				Price: price,
			})
		} else {
			slc = append(slc, model.Candle{
				Time:  timeStamp,
				Price: price,
			})
		}
		mp[lineParts[1]] = slc
	}

	return mp, nil
}
