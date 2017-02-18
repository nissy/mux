package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	//m.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/ | static"))
	//})
	//
	//m.Get("/a", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/a | static"))
	//})
	//
	//m.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/* | asterisk"))
	//})
	//
	//m.Get("/a/:id", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/a/:id | Param id is " + mux.URLParam(r, "id")))
	//})
	//
	//m.Get("/a/b/:id/d", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/a/b/:id/d | Param id is " + mux.URLParam(r, "id")))
	//})

	m.Get("/a/b/:id/d", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/d | Param id is " + mux.URLParam(r, "id")))
	})

	m.Get("/a/b/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id | Param id is " + mux.URLParam(r, "id")))
	})

	m.Get("/a/b/:id/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:id/:name | Param id is " + mux.URLParam(r, "id") + " " + mux.URLParam(r, "name")))
	})

	http.ListenAndServe(":8080", m)
}
