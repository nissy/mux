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
	charColon    = ':'
	charWildcard = '*'
	charSlash    = '/'
	byteColon    = byte(charColon)
	byteWildcard = byte(charWildcard)
	byteSlash    = byte(charSlash)
)

type (
	Mux struct {
		tree     []*node
		static   map[string]http.HandlerFunc
		NotFound http.HandlerFunc
	}

	node struct {
		child       map[byte]*node
		edge        byte
		handlerFunc http.HandlerFunc //End is handlerFunc != nil
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
	m := &Mux{
		static:   make(map[string]http.HandlerFunc),
		NotFound: http.NotFound,
	}

	m.tree = append(m.tree, &node{
		child: make(map[byte]*node),
	})

	return m
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
		if pattern[i] == byteColon || pattern[i] == byteWildcard {
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
	if pattern[0] != byteSlash {
		panic("There is no leading slash")
	}

	s := method + pattern

	if isStaticPattern(pattern) {
		m.static[s] = handlerFunc
		return
	}

	for i := 0; i < len(s); i++ {
		edge := s[i]

		if i > len(m.tree)-1 {
			m.tree = append(m.tree, &node{
				child: make(map[byte]*node),
			})
		}

		if _, ok := m.tree[i].child[edge]; ok {
			continue
		}

		n := &node{}

		if i == len(s)-1 {
			n.handlerFunc = handlerFunc
		}

		if edge == byteColon {
			for ii := i; ii < len(s); ii++ {
				if s[ii] == byteSlash {
					n.param = s[i+1 : ii]
					m.tree[i].child[edge] = n
					i += ii - 1

					break
				}
				if ii == len(s)-1 {
					n.handlerFunc = handlerFunc
					n.param = s[i+1 : ii+1]
					m.tree[i].child[edge] = n
					i += ii - 1

					break
				}
			}

			continue
		}

		if edge == byteWildcard {
			for ii := i; ii < len(s); ii++ {
				if s[ii] == byteSlash || ii == len(s)-1 {
					if ii == len(s)-1 {
						n.handlerFunc = handlerFunc
					}

					m.tree[i].child[edge] = n
					i += ii - 1

					break
				}
			}

			continue
		}

		m.tree[i].child[edge] = n
	}
}

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *rContext) {
	s := r.Method + r.URL.Path

	if fn, ok := m.static[s]; ok {
		return fn, nil
	}

	ctx := &rContext{}

	for i := 0; i < len(s); i++ {
		edge := s[i]

		if _, ok := m.tree[i].child[edge]; ok {
			if m.tree[i].child[edge].handlerFunc != nil {
				return m.tree[i].child[edge].handlerFunc, ctx
			}

			continue
		}

		if n, ok := m.tree[i].child[byteColon]; ok {
			for ii := i; ii < len(s); ii++ {
				if s[ii] == byteSlash {
					ctx.params.Set(n.param, s[i:ii])
					i += ii
					break
				}
				if ii == len(s)-1 {
					ctx.params.Set(n.param, s[i:ii+1])
					i += ii
					break
				}
			}
			if n.handlerFunc != nil {
				return n.handlerFunc, ctx
			}

			continue
		}

		if n, ok := m.tree[i].child[byteWildcard]; ok {
			for ii := i; ii < len(s); ii++ {
				if s[ii] == byteSlash || ii+i == len(s)-1 {
					i += ii
					break
				}
			}

			if n.handlerFunc != nil {
				return n.handlerFunc, ctx
			}

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
