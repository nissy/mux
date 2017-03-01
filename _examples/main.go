package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	m.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users"))
	})

	m.Get("/users/:name", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:name | Param id is " + mux.URLParam(r, "name")))
	})

	m.Get("/users/taro/man", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/taro/man"))
	})

	m.Get("/users/hanako/woman", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/hanako/woman"))
	})

	m.Get("/users/akira/:sex", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/akira/:sex | Param id is " + mux.URLParam(r, "sex")))
	})

	m.Get("/users/:name/:sex", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:name/:sex | Param id is " + mux.URLParam(r, "name") + " " + mux.URLParam(r, "sex")))
	})

	m.Get("/users/:name/woman", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:name/woman | Param id is " + mux.URLParam(r, "name")))
	})

	http.ListenAndServe(":8080", m)
}
