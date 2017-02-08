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

	characterParam    = ":"
	characterWildCard = "*"
)

var NotFound = http.NotFound

type (
	Mux struct {
		node *node
	}

	node struct {
		static nodeStatic
		param  nodeParam
	}

	nodeStatic map[route]http.HandlerFunc

	nodeParam struct {
		route       route
		params      []string
		handlerFunc http.HandlerFunc
	}

	route struct {
		method  string
		pattern string
	}
)

func NewMux() *Mux {
	return &Mux{
		node: &node{
			static: make(map[route]http.HandlerFunc),
		},
	}
}

func isParamPattern(pattern string) bool {
	for _, v := range []string{characterParam, characterWildCard} {
		if strings.Contains(pattern, v) {
			return true
		}
	}

	return false
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if isParamPattern(pattern) {
		return
	}

	mx.node.static[route{
		method:  method,
		pattern: pattern,
	}] = handlerFunc
}

func (n nodeStatic) routing(r *http.Request) http.HandlerFunc {
	if fn, ok := n[route{
		method:  r.Method,
		pattern: r.URL.Path,
	}]; ok {
		return fn
	}

	return nil
}

func (n nodeParam) routing(r *http.Request) http.HandlerFunc {
	return nil
}

func (m *Mux) handler(r *http.Request) http.HandlerFunc {
	if fn := m.node.static.routing(r); fn != nil {
		return fn
	}

	if fn := m.node.param.routing(r); fn != nil {
		return fn
	}

	return NotFound
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler(r)(w, r)
	return
}
