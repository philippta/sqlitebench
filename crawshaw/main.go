package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

func main() {
	poolSize := flag.Int("poolsize", runtime.NumCPU()/2, "sqlite pool size")
	flag.Parse()

	db, err := sqlitex.Open("file:bench.db", 0, *poolSize)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn := db.Get(r.Context())
		defer db.Put(conn)

		var id, value int
		fn := func(stmt *sqlite.Stmt) error {
			id = int(stmt.GetInt64("id"))
			value = int(stmt.GetInt64("value"))
			return nil
		}
		if err := sqlitex.Exec(conn, "select * from foo where rowid = ? limit 1", fn, rand.Intn(100000)+1); err != nil {
			panic(err)
		}
		w.Write([]byte(`{"id":` + strconv.Itoa(id) + `,"value":` + strconv.Itoa(value) + `}`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
