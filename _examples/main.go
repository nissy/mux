package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	})

	m.Get("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a"))
	})

	m.Get("/a/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/:id | " + mux.URLParam(r, "id")))
	})

	m.Get("/a/b/:id/c", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/c | " + mux.URLParam(r, "id")))
	})

	m.Get("/a/b/:id/cc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/cc | " + mux.URLParam(r, "id")))
	})

	http.ListenAndServe(":8080", m)
}
