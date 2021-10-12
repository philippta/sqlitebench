package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	poolSize := flag.Int("poolsize", runtime.NumCPU()/2, "sqlite pool size")
	flag.Parse()

	db, err := sql.Open("sqlite3", "file:bench.db?_journal=WAL&_synchronous=NORMAL")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(*poolSize)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var id, value int
		if err := db.QueryRow("select * from foo where rowid = ?", rand.Intn(100000)+1).Scan(&id, &value); err != nil {
			panic(err)
		}
		w.Write([]byte(`{"id":` + strconv.Itoa(id) + `,"value":` + strconv.Itoa(value) + `}`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
