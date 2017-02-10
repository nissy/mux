package mux

import (
	"context"
	"net/http"
	"strings"
)

const (
	GET         = "GET"
	POST        = "POST"
	PUT         = "PUT"
	DELETE      = "DELETE"
	HEAD        = "HEAD"
	OPTIONS     = "OPTIONS"
	ctxRouteKey = "mux"
)

var (
	characterColon    = ':'
	characterWildCard = '*'
	characterSlash    = '/'
	byteColon         = byte(characterColon)
	byteWildCard      = byte(characterWildCard)
	byteSlash         = byte(characterSlash)
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

	routeContext struct {
		URLParams map[string]string
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

func newRouteParam(method, pattern string) routeParam {
	return routeParam{
		method:   method,
		dirIndex: dirIndex(pattern),
	}
}

func newRouteContext() *routeContext {
	return &routeContext{
		URLParams: make(map[string]string),
	}
}

func routeContextExtract(r *http.Request) *routeContext {
	return r.Context().Value(ctxRouteKey).(*routeContext)
}

func isParamPattern(pattern string) bool {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == byteColon || pattern[i] == byteWildCard {
			return true
		}
	}

	return false
}

func URLParam(r *http.Request, key string) string {
	if ctx := routeContextExtract(r); ctx != nil {
		return ctx.URLParams[key]
	}

	return ""
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != byteSlash {
		panic("There is no leading slash")
	}

	if isParamPattern(pattern) {
		rtp := newRouteParam(method, pattern)
		mx.node.param.route[rtp] = append(mx.node.param.route[rtp], routeParamPattern{
			pattern:     pattern,
			handlerFunc: handlerFunc,
		})

		return
	}

	mx.node.static[newRouteStatic(method, pattern)] = handlerFunc
}

func dirIndex(dir string) (n int) {
	for i := 0; i < len(dir); i++ {
		if dir[i] == byteSlash {
			n++
		}
	}

	if n > 0 {
		return n - 1
	}

	return 0
}

func dirSplit(dir string) []string {
	return strings.FieldsFunc(dir, func(r rune) bool {
		return r == characterSlash
	})
}

func (n nodeStatic) lookup(r *http.Request) http.HandlerFunc {
	if fn, ok := n[newRouteStatic(r.Method, r.URL.Path)]; ok {
		return fn
	}

	return nil
}

func (n nodeParam) lookup(r *http.Request) http.HandlerFunc {
	rDirs := dirSplit(r.URL.Path)
	ctx := routeContextExtract(r)

	for _, v := range n.route[newRouteParam(r.Method, r.URL.Path)] {
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
				ctx.URLParams[vv[1:]] = rDirs[i]

				if i == nDirIndex {
					return v.handlerFunc
				}

				continue
			}

			if vv == rDirs[i] {
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

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(
		r.Context(), ctxRouteKey, newRouteContext()),
	)

	if fn := m.node.static.lookup(r); fn != nil {
		fn.ServeHTTP(w, r)
		return
	}

	if fn := m.node.param.lookup(r); fn != nil {
		fn.ServeHTTP(w, r)
		return
	}

	m.NotFound.ServeHTTP(w, r)
}
