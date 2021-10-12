package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"

	_ "modernc.org/sqlite"
)

func main() {
	poolSize := flag.Int("poolsize", runtime.NumCPU()/2, "sqlite pool size")
	flag.Parse()

	db, err := sql.Open("sqlite", "file:bench.db?_pragma=journal_mode=WAL&_pragma=synchronous=normal")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(*poolSize)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var id, value int
		err := db.QueryRow("select * from foo where rowid = ?", rand.Intn(100000)+1).Scan(&id, &value)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(`{"id":` + strconv.Itoa(id) + `,"value":` + strconv.Itoa(value) + `}`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
