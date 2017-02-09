package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	m.Entry(mux.GET, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo static world"))
	})

	m.Entry(mux.GET, "/a/b/c/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo asterisk world"))
	})

	m.Entry(mux.GET, "/a/b/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo param id is " + mux.Params["id"] + " world"))
	})

	http.ListenAndServe(":8080", m)
}
