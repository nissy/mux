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
	colon    = byte(':')
	slash    = byte('/')
	wildcard = byte('*')
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
		child       map[byte]*node
		handlerFunc http.HandlerFunc //end is handlerFunc != nil
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
		child:  make(map[byte]*node),
	}
}

func (n *node) findChild(edge byte) *node {
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
		if pattern[i] == colon || pattern[i] == wildcard {
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
	if pattern[0] != slash {
		panic("There is no leading slash")
	}

	if isStaticPattern(pattern) {
		m.static[pattern] = handlerFunc
		return
	}

	number := 0
	parent := m.tree[0]

	for i := 0; i < len(pattern); i++ {
		edge := pattern[i]

		child := &node{
			number: number,
		}

		if n := parent.findChild(edge); n != nil {
			child = n
		}

		if edge == colon {
			p := []byte{}
			i += 1

			for ; i < len(pattern); i++ {
				if pattern[i] == slash {
					i -= 1
					break
				}

				p = append(p, pattern[i])
			}

			child.param = string(p)
		}

		if edge == wildcard {
			for ; i < len(pattern); i++ {
				if pattern[i] == slash {
					i -= 1
					break
				}
			}
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
			parent.child = make(map[byte]*node)
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

	ctx := &rCtx{}
	parent := m.tree[0]

	var route [2]int

	for i := 0; i < len(s); i++ {
		edge := s[i]
		child := parent.findChild(edge)

		if child != nil {
			if i == len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			route[1] = route[0]
			route[0] = child.number
			parent = child
			continue
		}

		//PARAM
		if n := parent.findChild(colon); n != nil {
			p := []byte{}

			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}

				p = append(p, s[i])
			}

			ctx.params.Set(n.param, string(p))
			child = n

		} else if n := parent.findChild(wildcard); n != nil {
			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}
			}

			child = n
		}

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			route[1] = route[0]
			route[0] = child.number
			parent = child
			continue
		}

		//BACKTRACK PARAM
		if route[1] > 0 {
			if n := m.tree[route[1]].findChild(colon); n != nil {
				p := []byte{}
				i -= 1

				for ; i < len(s); i++ {
					if s[i] == slash {
						i -= 1
						break
					}

					p = append(p, s[i])
				}

				ctx.params.Set(n.param, string(p))
				child = n

			} else if n := m.tree[route[1]].findChild(wildcard); n != nil {
				for ; i < len(s); i++ {
					if s[i] == slash {
						i -= 1
						break
					}
				}

				child = n
			}

			if child != nil {
				if i >= len(s)-1 {
					if child.handlerFunc != nil {
						return child.handlerFunc, ctx
					}
				}

				route[1] = route[0]
				route[0] = child.number
				parent = child
				continue
			}
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
