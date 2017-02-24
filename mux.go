package mux

import (
	"context"
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	CONNECT = "CONNECT"
	TRACE   = "TRACE"
	rCtxKey = "mux"
)

var (
	colon     = ":"
	wildcard  = "*"
	bColon    = byte(':')
	bSlash    = byte('/')
	bWildcard = byte('*')
)

type (
	Router struct {
		mux      [9]*mux
		NotFound http.HandlerFunc
	}

	mux struct {
		tree   []*node
		static map[string]http.HandlerFunc
	}

	node struct {
		number      int
		child       map[string]*node
		handlerFunc http.HandlerFunc
		param       string
	}

	rCtx struct {
		params params
	}

	params []param

	param struct {
		key   string
		value string
	}
)

func New() *Router {
	r := &Router{
		NotFound: http.NotFound,
	}

	for i := 0; i < 9; i++ {
		r.mux[i] = newMux()
	}

	return r
}

func newMux() *mux {
	m := &mux{
		static: make(map[string]http.HandlerFunc),
	}

	m.tree = append(m.tree, newNode(0))
	return m
}

func (r *Router) enter(method string) *mux {
	switch method {
	case GET:
		return r.mux[0]
	case POST:
		return r.mux[1]
	case PUT:
		return r.mux[2]
	case DELETE:
		return r.mux[3]
	case HEAD:
		return r.mux[4]
	case OPTIONS:
		return r.mux[5]
	case PATCH:
		return r.mux[6]
	case CONNECT:
		return r.mux[7]
	case TRACE:
		return r.mux[8]
	}

	return nil
}

func newNode(number int) *node {
	return &node{
		number: number,
		child:  make(map[string]*node),
	}
}

func (n *node) findChild(edge string) *node {
	if n, ok := n.child[edge]; ok {
		return n
	}

	return nil
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

func URLParam(r *http.Request, key string) string {
	if ctx := r.Context().Value(rCtxKey); ctx != nil {
		if ctx, ok := ctx.(*rCtx); ok {
			return ctx.params.Get(key)
		}
	}

	return ""
}

func isStaticPattern(pattern string) bool {
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == bColon || pattern[i] == bWildcard {
			return false
		}
	}

	return true
}

func (r *Router) Get(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(GET).handle(pattern, handlerFunc)
}

func (r *Router) Post(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(POST).handle(pattern, handlerFunc)
}

func (r *Router) Put(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(PUT).handle(pattern, handlerFunc)
}

func (r *Router) Delete(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(DELETE).handle(pattern, handlerFunc)
}

func (r *Router) Head(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(HEAD).handle(pattern, handlerFunc)
}

func (r *Router) Options(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(OPTIONS).handle(pattern, handlerFunc)
}

func (r *Router) Patch(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(PATCH).handle(pattern, handlerFunc)
}

func (r *Router) Connect(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(CONNECT).handle(pattern, handlerFunc)
}

func (r *Router) Trace(pattern string, handlerFunc http.HandlerFunc) {
	r.enter(TRACE).handle(pattern, handlerFunc)
}

func (m *mux) handle(pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != bSlash {
		panic("There is no leading bSlash")
	}

	if isStaticPattern(pattern) {
		m.static[pattern] = handlerFunc
		return
	}

	var number, si, ei int
	parent := m.tree[0]

	for i := 0; i < len(pattern); i++ {
		if pattern[i] == bSlash {
			i += 1
		}

		si = i
		ei = i

		for ; i < len(pattern); i++ {
			if pattern[i] == bSlash {
				break
			}

			ei++
		}

		edge := pattern[si:ei]
		var param string

		if edge[0] == bColon {
			param = edge[1:]
			edge = colon
		}

		child := &node{
			number: number,
		}

		if n := parent.findChild(edge); n != nil {
			child = n
		}

		if len(param) > 0 {
			child.param = param
		}

		if i >= len(pattern)-1 {
			child.handlerFunc = handlerFunc
		}

		if _, ok := parent.child[edge]; ok {
			parent = child
			continue
		}

		if number < len(m.tree)-1 {
			number = len(m.tree)
		} else {
			number += 1
		}

		// Not have brother
		if len(parent.child) == 0 {
			parent.child = make(map[string]*node)
		}

		child.number = number
		m.tree = append(m.tree, child)
		parent.child[edge] = child
		parent = child
	}
}

func (m *mux) lookup(r *http.Request) (http.HandlerFunc, *rCtx) {
	s := r.URL.Path

	if fn, ok := m.static[s]; ok {
		return fn, nil
	}

	if len(m.tree) == 0 {
		return nil, nil
	}

	var si, ei, bsi, route int

	parent := m.tree[0]
	ctx := &rCtx{}

	for i := 0; i < len(s); i++ {
		if s[i] == bSlash {
			i += 1
		}

		si = i
		ei = i

		for ; i < len(s); i++ {
			if s[i] == bSlash {
				break
			}

			ei++
		}

		edge := s[si:ei]
		child := parent.findChild(edge)

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			bsi = si
			route = child.number
			parent = child
			continue
		}

		if n := parent.findChild(colon); n != nil {
			ctx.params.Set(n.param, edge)
			child = n
		} else if n := parent.findChild(wildcard); n != nil {
			child = n
		}

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			bsi = si
			route = child.number
			parent = child
			continue
		}

		//BACKTRACK
		if n := m.tree[route].findChild(colon); n != nil {
			ctx.params.Set(n.param, s[bsi:si-1])
			child = n
		} else if n := m.tree[route].findChild(wildcard); n != nil {
			child = n
		}

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			route = child.number
			parent = child
			continue
		}

		break
	}

	return nil, nil
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fn, ctx := rt.enter(r.Method).lookup(r); fn != nil {
		if ctx != nil {
			fn.ServeHTTP(w, r.WithContext(context.WithValue(
				r.Context(), rCtxKey, ctx),
			))
			return
		}

		fn.ServeHTTP(w, r)
		return
	}

	rt.NotFound.ServeHTTP(w, r)
}
