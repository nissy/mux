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

	m.Get("/a/b/:id/d", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/d | " + mux.URLParam(r, "id")))
	})

	m.Get("/a/b/:id/dd", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/dd | " + mux.URLParam(r, "id")))
	})

	http.ListenAndServe(":8080", m)
}
