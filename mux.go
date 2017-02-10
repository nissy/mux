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
	charColon    = ':'
	charWildCard = '*'
	charSlash    = '/'
	byteColon    = byte(charColon)
	byteWildCard = byte(charWildCard)
	byteSlash    = byte(charSlash)
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
		dirs        []string
		dirIndex    int
	}

	routeContext struct {
		params params
	}

	params []param

	param struct {
		key   string
		value string
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

func newRouteParam(method string, index int) routeParam {
	return routeParam{
		method:   method,
		dirIndex: index,
	}
}

func newRouteContext() *routeContext {
	return &routeContext{}
}

func routeContextExtract(r *http.Request) *routeContext {
	return r.Context().Value(ctxRouteKey).(*routeContext)
}

func (ps *params) Set(key, value string) {
	*ps = append(*ps, param{
		key:   key,
		value: value,
	})
}

func (ps params) Get(key string) string {
	for _, v := range ps {
		if v.key == key {
			return v.value
		}
	}

	return ""
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
		return ctx.params.Get(key)
	}

	return ""
}

func (mx *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != byteSlash {
		panic("There is no leading slash")
	}

	if isParamPattern(pattern) {
		dirs, dirIndex := dirSplit(pattern)
		rt := newRouteParam(method, dirIndex)
		mx.node.param.route[rt] = append(mx.node.param.route[rt], routeParamPattern{
			pattern:     pattern,
			handlerFunc: handlerFunc,
			dirs:        dirs,
			dirIndex:    dirIndex,
		})

		return
	}

	mx.node.static[newRouteStatic(method, pattern)] = handlerFunc
}

func dirSplit(dir string) ([]string, int) {
	dirs := strings.FieldsFunc(dir, func(r rune) bool {
		return r == charSlash
	})

	if len(dirs) > 0 {
		return dirs, len(dirs) - 1
	}

	return dirs, 0
}

func (n nodeStatic) lookup(r *http.Request) http.HandlerFunc {
	if fn, ok := n[newRouteStatic(r.Method, r.URL.Path)]; ok {
		return fn
	}

	return nil
}

func (n nodeParam) lookup(r *http.Request) http.HandlerFunc {
	rDirs, rDirIndex := dirSplit(r.URL.Path)
	ctx := routeContextExtract(r)

	for _, v := range n.route[newRouteParam(r.Method, rDirIndex)] {
		for i, vv := range v.dirs {
			if vv[0] == byteWildCard {
				if i == v.dirIndex {
					return v.handlerFunc
				}

				continue
			}

			if vv[0] == byteColon {
				ctx.params.Set(vv[1:], rDirs[i])

				if i == v.dirIndex {
					return v.handlerFunc
				}

				continue
			}

			if vv == rDirs[i] {
				if i == v.dirIndex {
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
