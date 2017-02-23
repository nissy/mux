package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.New()

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

	//m.Get("/a/:id", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/a/:id | Param id is " + mux.URLParam(r, "id")))
	//})

	m.Get("/a/b/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/b/:name | Param id is " + mux.URLParam(r, "name")))
	})

	m.Get("/a/:id/c", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/:id/c | Param id is " + mux.URLParam(r, "id")))
	})

	m.Get("/a/:id/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a/:id/:name | Param id is " + mux.URLParam(r, "id") + " / " + mux.URLParam(r, "name")))
	})

	//m.Get("/b/:id", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/b/:id | Param id is " + mux.URLParam(r, "id")))
	//})
	//
	//m.Get("/b/:id/c", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/b/:id/c | Param id is " + mux.URLParam(r, "id")))
	//})
	//
	//m.Get("/b/:id/:name", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("/b/:id/:name | Param id is " + mux.URLParam(r, "id") + " / " + mux.URLParam(r, "name")))
	//})

	http.ListenAndServe(":8080", m)
}
