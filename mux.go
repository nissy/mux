package mux

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var NotFound = http.NotFound

type Mux struct {
	statics map[route]http.HandlerFunc
	params  map[route]http.HandlerFunc
}

type route struct {
	method  string
	pattern string
}

func NewMux() *Mux {
	return &Mux{
		statics: make(map[route]http.HandlerFunc),
	}
}

func isParamPattern(pattern string) (ok bool) {
	if ok = strings.Contains(pattern, ":"); !ok {
		return strings.Contains(pattern, "*")
	}

	return ok
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if isParamPattern(pattern) {
		return
	}

	mx.statics[route{
		method:  method,
		pattern: pattern,
	}] = handlerFunc
}

func (m *Mux) match(r *http.Request) http.HandlerFunc {
	u, err := url.Parse(r.RequestURI)

	if err != nil {
		panic(err)
	}

	//static
	rt := route{
		method:  r.Method,
		pattern: u.Path,
	}

	if _, ok := m.statics[rt]; ok {
		return m.statics[rt]
	}

	return NotFound
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.match(r)(w, r)
	return
}
