package mux

import (
	"context"
	"net/http"

	"github.com/k0kubun/pp"
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

		parent.child = make(map[byte]*node)

		//pp.Println(number, string(edge))

		number += 1
		child.number = number
		m.tree = append(m.tree, child)
		parent.child[edge] = child
		parent = child
	}
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
	var treeIndex int

	for i := 0; i < len(s); i++ {
		if treeIndex > len(m.tree)-1 {
			break
		}

		edge := s[i]

		if n := m.findNode(treeIndex, edge); n != nil {
			if i == len(s)-1 {
				if n.handlerFunc != nil {
					return n.handlerFunc, ctx
				}

				// BACKTRACK
				if treeIndex > 0 {
					if n := m.findNode(treeIndex, colon); n != nil {
						ctx.params.Set(n.param, string(s[i]))

						if n.handlerFunc != nil {
							return n.handlerFunc, ctx
						}
					}

					if n := m.findNode(treeIndex, wildcard); n != nil {
						if n.handlerFunc != nil {
							return n.handlerFunc, ctx
						}
					}
				}
			}

			treeIndex += 1
			continue
		}

		// PARAM
		if n := m.findNodeBack(treeIndex, colon); n != nil {
			p := []byte{}

			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}

				p = append(p, s[i])
			}

			ctx.params.Set(n.param, string(p))

			if i >= len(s)-1 {
				if n.handlerFunc != nil {
					return n.handlerFunc, ctx
				}
			}

			treeIndex += 1
			continue
		}

		// WILDCARD
		if n := m.findNodeBack(treeIndex, wildcard); n != nil {
			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}
			}

			if i >= len(s)-1 {
				if n.handlerFunc != nil {
					return n.handlerFunc, ctx
				}
			}

			treeIndex += 1
			continue
		}

		break
	}

	return nil, nil
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pp.Println(m.tree)

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
