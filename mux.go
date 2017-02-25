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
	Mux struct {
		tree     []*node
		static   map[string]http.HandlerFunc
		NotFound http.HandlerFunc
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

func New() *Mux {
	return NewMux()
}

func NewMux() *Mux {
	m := &Mux{
		static:   make(map[string]http.HandlerFunc),
		NotFound: http.NotFound,
	}

	m.tree = append(m.tree, newNode(0))
	return m
}

func newNode(number int) *node {
	return &node{
		number: number,
		child:  make(map[string]*node),
	}
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

func (m *Mux) Get(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(GET, pattern, handlerFunc)
}

func (m *Mux) Post(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(POST, pattern, handlerFunc)
}

func (m *Mux) Put(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(PUT, pattern, handlerFunc)
}

func (m *Mux) Delete(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(DELETE, pattern, handlerFunc)
}

func (m *Mux) Head(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(HEAD, pattern, handlerFunc)
}

func (m *Mux) Options(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(OPTIONS, pattern, handlerFunc)
}

func (m *Mux) Patch(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(PATCH, pattern, handlerFunc)
}

func (m *Mux) Connect(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(CONNECT, pattern, handlerFunc)
}

func (m *Mux) Trace(pattern string, handlerFunc http.HandlerFunc) {
	m.Entry(TRACE, pattern, handlerFunc)
}

func (m *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != bSlash {
		panic("There is no leading slash")
	}

	s := method + pattern

	if isStaticPattern(s) {
		m.static[s] = handlerFunc
		return
	}

	var number, si, ei int
	parent := m.tree[0]

	for i := 0; i < len(s); i++ {
		if s[i] == bSlash {
			i += 1
		}

		si = i
		ei = i

		for ; i < len(s); i++ {
			if si < ei {
				if s[i] == bColon || s[i] == bWildcard {
					panic("Parameter are not first")
				}
			}

			if s[i] == bSlash {
				break
			}

			ei++
		}

		edge := s[si:ei]
		var param string

		switch edge[0] {
		case bColon:
			param = edge[1:]
			edge = colon
		case bWildcard:
			edge = wildcard
		}

		child := &node{
			number: number,
		}

		if n := parent.child[edge]; n != nil {
			child = n
		}

		if len(param) > 0 {
			child.param = param
		}

		if i >= len(s)-1 {
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

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *rCtx) {
	s := r.URL.Path

	//TODO allo
	if fn := m.static[r.Method+s]; fn != nil {
		return fn, nil
	}

	if len(m.tree) < 2 {
		return nil, nil
	}

	var parent, child *node

	if parent = m.tree[0].child[r.Method]; parent == nil {
		return nil, nil
	}

	var si, ei, bsi int
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

		//STATIC
		if child = parent.child[edge]; child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			bsi = si
			parent = child
			continue
		}

		//PARAM
		if child = parent.child[colon]; child != nil {
			ctx.params.Set(child.param, edge)
		} else {
			child = parent.child[wildcard]
		}

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			bsi = si
			parent = child
			continue
		}

		//BACKTRACK
		if child = m.tree[parent.number].child[colon]; child != nil {
			ctx.params.Set(child.param, s[bsi:si-1])
		} else {
			child = m.tree[parent.number].child[wildcard]
		}

		if child != nil {
			if i >= len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			parent = child
			continue
		}

		break
	}

	return nil, nil
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fn, ctx := m.lookup(r); fn != nil {
		if ctx != nil {
			fn.ServeHTTP(w, r.WithContext(context.WithValue(
				r.Context(), rCtxKey, ctx),
			))
			return
		}

		fn.ServeHTTP(w, r)
		return
	}

	m.NotFound.ServeHTTP(w, r)
}
