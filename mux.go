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
)

var (
	Params            = make(map[string]string)
	characterParam    = ":"
	characterWildCard = "*"
)

type (
	Mux struct {
		node     *node
		NotFound http.HandlerFunc
	}

	node struct {
		static nodeStatic
		param  nodeParam
	}

	nodeStatic map[routeStatic]http.HandlerFunc

	routeStatic struct {
		method  string
		pattern string
	}

	nodeParam struct {
		route       map[routeParam][]routeParamPattern
		params      []string
		handlerFunc http.HandlerFunc
	}

	routeParam struct {
		method   string
		dirIndex int
	}

	routeParamPattern struct {
		pattern     string
		handlerFunc http.HandlerFunc
	}
)

func NewMux() *Mux {
	return &Mux{
		node: &node{
			param: nodeParam{
				route: make(map[routeParam][]routeParamPattern),
			},
			static: make(map[routeStatic]http.HandlerFunc),
		},
		NotFound: http.NotFound,
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

func newRouteParam(method string, dirIndex int) routeParam {
	return routeParam{
		method:   method,
		dirIndex: dirIndex,
	}
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if isParamPattern(pattern) {
		rtp := newRouteParam(method, dirIndex(pattern))
		mx.node.param.route[rtp] = append(mx.node.param.route[rtp], routeParamPattern{
			pattern:     pattern,
			handlerFunc: handlerFunc,
		})

		return
	}

	mx.node.static[routeStatic{
		method:  method,
		pattern: pattern,
	}] = handlerFunc
}

func (n nodeStatic) routing(r *http.Request) http.HandlerFunc {
	if fn, ok := n[routeStatic{
		method:  r.Method,
		pattern: r.URL.Path,
	}]; ok {
		return fn
	}

	return nil
}

func dirIndex(dir string) int {
	return len(dirSplit(dir)) - 1
}

func dirSplit(dir string) (s []string) {
	for _, v := range strings.Split(dir, "/") {
		if len(v) > 0 {
			s = append(s, v)
		}
	}

	return s
}

func (n nodeParam) routing(r *http.Request) http.HandlerFunc {
	rDir := dirSplit(r.URL.Path)

	for _, v := range n.route[newRouteParam(r.Method, dirIndex(r.URL.Path))] {
		sDir := dirSplit(v.pattern)
		sDirIndex := len(dirSplit(v.pattern)) - 1

		for i, vv := range sDir {
			if strings.Contains(vv, characterWildCard) {
				if i == sDirIndex {
					return v.handlerFunc
				}

				continue
			}

			if string(vv[0]) == characterParam {
				Params[vv[1:]] = string(rDir[i])

				if i == sDirIndex {
					return v.handlerFunc
				}

				continue
			}

			if rDir[i] == vv {
				if i == sDirIndex {
					return v.handlerFunc
				}

				continue
			}

			break
		}

	}

	return nil
}

func (m *Mux) handler(r *http.Request) http.HandlerFunc {
	if fn := m.node.static.routing(r); fn != nil {
		return fn
	}

	if fn := m.node.param.routing(r); fn != nil {
		return fn
	}

	return m.NotFound
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler(r)(w, r)
	return
}
