package mux

import (
	"net/http"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"

	ParamCharacter    = ":"
	WildcardCharacter = "*"
)

var NotFound = http.NotFound

type Mux struct {
	static map[route]http.HandlerFunc
	//dynamic map[route]http.HandlerFunc
}

type route struct {
	method  string
	pattern string
}

func NewMux() *Mux {
	return &Mux{
		static: make(map[route]http.HandlerFunc),
	}
}

func isDynamicPattern(pattern string) bool {
	for _, v := range []string{ParamCharacter, WildcardCharacter} {
		if strings.Contains(pattern, v) {
			return true
		}
	}

	return false
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if isDynamicPattern(pattern) {
		return
	}

	mx.static[route{
		method:  method,
		pattern: pattern,
	}] = handlerFunc
}

func (m *Mux) handler(r *http.Request) http.HandlerFunc {
	// static route
	if fn, ok := m.static[route{
		method:  r.Method,
		pattern: r.URL.Path,
	}]; ok {
		return fn
	}

	// dynamic route

	return NotFound
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler(r)(w, r)
	return
}
