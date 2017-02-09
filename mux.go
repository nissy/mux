package mux

import (
	"net/http"
	"strings"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

var Params = make(map[string]string)

var (
	characterColon    = ":"
	characterWildCard = "*"
	characterSlash    = "/"
	byteColon         = []byte(characterColon)[0]
	byteWildCard      = []byte(characterWildCard)[0]
	byteSlash         = []byte(characterSlash)[0]
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

func newRouteStatic(method, pattern string) routeStatic {
	return routeStatic{
		method:  method,
		pattern: pattern,
	}
}

func newRouteParam(method string, dirIndex int) routeParam {
	return routeParam{
		method:   method,
		dirIndex: dirIndex,
	}
}

func isParamPattern(pattern string) bool {
	for _, v := range []string{characterColon, characterWildCard} {
		if strings.Contains(pattern, v) {
			return true
		}
	}

	return false
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != byteSlash {
		panic("There is no leading slash")
	}

	if isParamPattern(pattern) {
		rtp := newRouteParam(method, dirIndex(pattern))
		mx.node.param.route[rtp] = append(mx.node.param.route[rtp], routeParamPattern{
			pattern:     pattern,
			handlerFunc: handlerFunc,
		})

		return
	}

	mx.node.static[newRouteStatic(method, pattern)] = handlerFunc
}

func (n nodeStatic) routing(r *http.Request) http.HandlerFunc {
	if fn, ok := n[newRouteStatic(r.Method, r.URL.Path)]; ok {
		return fn
	}

	return nil
}

func dirIndex(dir string) (n int) {
	for i := 0; i < len(dir); i++ {
		if dir[i] == byteSlash {
			n++
		}
	}

	return n - 1
}

func dirSplit(dir string) (ds []string) {
	for _, v := range strings.Split(dir, "/") {
		if len(v) > 0 {
			ds = append(ds, v)
		}
	}

	return ds
}

func (n nodeParam) routing(r *http.Request) http.HandlerFunc {
	rDirs := dirSplit(r.URL.Path)

	for _, v := range n.route[newRouteParam(r.Method, dirIndex(r.URL.Path))] {
		nDirs := dirSplit(v.pattern)
		nDirIndex := len(nDirs) - 1

		for i, vv := range nDirs {
			if vv[0] == byteWildCard {
				if i == nDirIndex {
					return v.handlerFunc
				}

				continue
			}

			if vv[0] == byteColon {
				Params[vv[1:]] = string(rDirs[i])

				if i == nDirIndex {
					return v.handlerFunc
				}

				continue
			}

			if rDirs[i] == vv {
				if i == nDirIndex {
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
