package mux

import (
	"context"
	"net/http"
)

const (
	GET         = "GET"
	POST        = "POST"
	PUT         = "PUT"
	DELETE      = "DELETE"
	HEAD        = "HEAD"
	OPTIONS     = "OPTIONS"
	PATCH       = "PATCH"
	ctxRouteKey = "mux"
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
		child       map[byte]*node
		handlerFunc http.HandlerFunc //end is handlerFunc != nil
		param       string
	}

	rContext struct {
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
		static:   make(map[string]http.HandlerFunc),
		NotFound: http.NotFound,
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
	if ctx := r.Context().Value(ctxRouteKey); ctx != nil {
		if ctx, ok := ctx.(*rContext); ok {
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

	var treeIndex int

	for i := 0; i < len(s); i++ {
		edge := s[i]

		if treeIndex == len(m.tree) {
			m.tree = append(m.tree, &node{
				child: make(map[byte]*node),
			})
		}

		if _, ok := m.tree[treeIndex].child[edge]; ok {
			if edge == colon || edge == wildcard {
				for ; i < len(s); i++ {
					if s[i] == slash {
						i -= 1
						break
					}
				}
			}

			treeIndex += 1
			continue
		}

		n := &node{}

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

			n.param = string(p)
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
			n.handlerFunc = handlerFunc
		}

		m.tree[treeIndex].child[edge] = n
		treeIndex += 1
	}
}

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *rContext) {
	s := r.Method + r.URL.Path

	if fn, ok := m.static[s]; ok {
		return fn, nil
	}

	if len(m.tree) == 0 {
		return nil, nil
	}

	ctx := &rContext{}
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

		// PARAM NODE
		var pNode *node

		if pNode = m.findNode(treeIndex, colon); pNode == nil {
			// BACKTRACK
			if treeIndex > 0 {
				pNode = m.findNode(treeIndex-1, colon)
			}
		}

		if pNode != nil {
			p := []byte{}

			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}

				p = append(p, s[i])
			}

			ctx.params.Set(pNode.param, string(p))

			if i >= len(s)-1 {
				if pNode.handlerFunc != nil {
					return pNode.handlerFunc, ctx
				}
			}

			treeIndex += 1
			continue
		}

		// WILDCARD NODE
		var wNode *node

		if wNode = m.findNode(treeIndex, wildcard); wNode == nil {
			// BACKTRACK
			if treeIndex > 0 {
				wNode = m.findNode(treeIndex-1, wildcard)
			}
		}

		if wNode != nil {
			for ; i < len(s); i++ {
				if s[i] == slash {
					i -= 1
					break
				}
			}

			if i >= len(s)-1 {
				if wNode.handlerFunc != nil {
					return wNode.handlerFunc, ctx
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
	if fn, ctx := m.lookup(r); fn != nil {
		if ctx != nil {
			fn.ServeHTTP(w, r.WithContext(context.WithValue(
				r.Context(), ctxRouteKey, ctx),
			))
			return
		}

		fn.ServeHTTP(w, r)
		return
	}

	m.NotFound.ServeHTTP(w, r)
}
