package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.New()

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	})

	m.Get("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a"))
	})

	m.Get("/a/b/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:name | Param id is " + mux.URLParam(r, "name")))
	})

	m.Get("/a/:id/c", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/:id/c | Param id is " + mux.URLParam(r, "id")))
	})

	m.Get("/a/:id/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/:id/:name | Param id is " + mux.URLParam(r, "id") + " / " + mux.URLParam(r, "name")))
	})

	http.ListenAndServe(":8080", m)
}
