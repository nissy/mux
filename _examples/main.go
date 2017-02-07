package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	m.Entry(mux.GET, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo World !"))
	})

	http.ListenAndServe(":8080", m)
}
