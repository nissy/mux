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
	rCtxKey = "mux"
)

var (
	colon    = byte(':')
	slash    = byte('/')
	wildcard = byte('*')
)

type (
	Mux struct {
		tree     []*node
		static   map[string]http.HandlerFunc
		NotFound http.HandlerFunc
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
		child:  make(map[byte]*node),
	}
}

func (m *Mux) findNode(index int, edge byte) *node {
	if len(m.tree) >= index {
		if n, ok := m.tree[index].child[edge]; ok {
			return n
		}
	}

	return nil
}

func (m *Mux) findNodeBack(index int, edge byte) *node {
	if n := m.findNode(index, edge); n != nil {
		return n
	}

	if index > 0 {
		if n := m.findNode(index-1, edge); n != nil {
			return n
		}
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

func (m *Mux) Entry(method, pattern string, handlerFunc http.HandlerFunc) {
	if pattern[0] != slash {
		panic("There is no leading slash")
	}

	s := method + pattern

	if isStaticPattern(pattern) {
		m.static[s] = handlerFunc
		return
	}

	number := 1
	parent := m.tree[0]

	for i := 0; i < len(s); i++ {
		edge := s[i]

		child := &node{
			number: number,
		}

		if n, _ := parent.child[edge]; n != nil {
			child = n
		}

		if edge == colon {
			p := []byte{}
			i += 1

			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}

				p = append(p, s[i])
			}

			child.param = string(p)
		}

		if edge == wildcard {
			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}
			}
		}

		if i >= len(s)-1 {
			child.handlerFunc = handlerFunc
		}

		if _, ok := parent.child[edge]; ok {
			number += 1
			parent = child
			continue
		}

		if number < len(m.tree) {
			number = len(m.tree)
		}

		// 兄弟がいない
		if len(parent.child) == 0 {
			parent.child = make(map[byte]*node)
		}

		number += 1
		child.number = number
		m.tree = append(m.tree, child)
		parent.child[edge] = child
		parent = child
	}
}

func (n *node) findChild(edge byte) *node {
	if n, ok := n.child[edge]; ok {
		return n
	}

	return nil
}

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *rCtx) {
	s := r.Method + r.URL.Path

	if fn, ok := m.static[s]; ok {
		return fn, nil
	}

	if len(m.tree) == 0 {
		return nil, nil
	}

	ctx := &rCtx{}
	//var treeIndex int

	parent := m.tree[0]

	var route []*node
	route = append(route, parent)

	for i := 0; i < len(s); i++ {
		edge := s[i]

		child := parent.findChild(edge)

		if child != nil {
			if i == len(s)-1 {
				if child.handlerFunc != nil {
					return child.handlerFunc, ctx
				}
			}

			route = append(route, child)
			parent = child
			continue
		}

		// PARAM
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
		}

		// WILDCARD
		if n := parent.findChild(wildcard); n != nil {
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

			route = append(route, child)
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
