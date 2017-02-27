package mux

import (
	"context"
	"net/http"
)

const (
	GET        = "GET"
	POST       = "POST"
	PUT        = "PUT"
	DELETE     = "DELETE"
	HEAD       = "HEAD"
	OPTIONS    = "OPTIONS"
	PATCH      = "PATCH"
	CONNECT    = "CONNECT"
	TRACE      = "TRACE"
	ContextKey = "mux"
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
		tree          []*node
		NotFound      http.HandlerFunc
		maxPramNumber int
	}

	node struct {
		number      int
		child       map[string]*node
		handlerFunc http.HandlerFunc
		param       string
	}

	Context struct {
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

func (n *node) newChild(child *node, edge string) *node {
	if len(n.child) == 0 {
		n.child = make(map[string]*node)
	}

	n.child[edge] = child
	return child
}

func newContext(cap int) *Context {
	return &Context{
		params: make([]param, 0, cap),
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
	if ctx := r.Context().Value(ContextKey); ctx != nil {
		if ctx, ok := ctx.(*Context); ok {
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

func minDirChoose(path string, pn int) int {
	var n int

	for i := 0; i < len(path); i++ {
		if path[i] == bSlash {
			n++
			i++

			if n >= pn {
				return pn
			}
		}
	}

	return n
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

	parent := m.tree[0].child[method]

	if parent == nil {
		parent = m.tree[0]
		child := newNode(len(m.tree))
		m.tree = append(m.tree, child)
		parent = parent.newChild(child, method)
	}

	if isStaticPattern(pattern) {
		if _, ok := parent.child[pattern]; !ok {
			child := newNode(len(m.tree))
			child.handlerFunc = handlerFunc
			m.tree = append(m.tree, child)
			parent.newChild(child, pattern)
		}

		return
	}

	var number, si, ei, pi int

	for i := 0; i < len(pattern); i++ {
		if pattern[i] == bSlash {
			i++
		}

		si = i
		ei = i

		for ; i < len(pattern); i++ {
			if si < ei {
				if pattern[i] == bColon || pattern[i] == bWildcard {
					panic("Parameter are not first")
				}
			}

			if pattern[i] == bSlash {
				break
			}

			ei++
		}

		edge := pattern[si:ei]
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
			pi++
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
			number++
		}

		child.number = number
		m.tree = append(m.tree, child)
		parent = parent.newChild(child, edge)
	}

	if pi > m.maxPramNumber {
		m.maxPramNumber = pi
	}
}

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *Context) {
	if len(m.tree) < 2 {
		return nil, nil
	}

	s := r.URL.Path
	var parent, child *node

	if parent = m.tree[0].child[r.Method]; parent == nil {
		return nil, nil
	}

	//STATIC
	if child = parent.child[s]; child != nil {
		return child.handlerFunc, nil
	}

	var si, ei, bsi int
	var ctx *Context

	for i := 0; i < len(s); i++ {
		if s[i] == bSlash {
			i++
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

		if child = parent.child[edge]; child == nil {
			if child = parent.child[colon]; child != nil {
				if ctx == nil {
					ctx = newContext(minDirChoose(s, m.maxPramNumber))
				}
				ctx.params.Set(child.param, edge)

			} else if child = parent.child[wildcard]; child == nil {
				//BACKTRACK
				if child = m.tree[parent.number].child[colon]; child != nil {
					if ctx == nil {
						ctx = newContext(minDirChoose(s, m.maxPramNumber))
					}
					ctx.params.Set(child.param, s[bsi:si-1])
					si = bsi

				} else if child = m.tree[parent.number].child[wildcard]; child != nil {
					si = bsi
				}
			}
		}

		if child != nil {
			if i >= len(s)-1 && child.handlerFunc != nil {
				return child.handlerFunc, ctx
			}

			bsi = si
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
				r.Context(), ContextKey, ctx),
			))
			return
		}

		fn.ServeHTTP(w, r)
		return
	}

	m.NotFound.ServeHTTP(w, r)
}
