package mux

import (
	"context"
	"net/http"
	"sync"
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
		tree     *node
		pool     sync.Pool
		maxParam int
		NotFound http.HandlerFunc
	}

	node struct {
		parent      *node
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

	m.pool = sync.Pool{
		New: func() interface{} {
			return newContext(m.maxParam)
		},
	}

	m.tree = newNode()
	return m
}

func newNode() *node {
	return &node{
		child: make(map[string]*node),
	}
}

func (n *node) newChild(child *node, edge string) *node {
	if len(n.child) == 0 {
		n.child = make(map[string]*node)
	}

	child.parent = n
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

	parent := m.tree.child[method]

	if parent == nil {
		parent = m.tree.newChild(newNode(), method)
	}

	if isStaticPattern(pattern) {
		if _, ok := parent.child[pattern]; !ok {
			child := newNode()
			child.handlerFunc = handlerFunc
			parent.newChild(child, pattern)
		}

		return
	}

	var si, ei, pi int

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

		child, exist := parent.child[edge]

		if !exist {
			child = newNode()
		}

		if len(param) > 0 {
			child.param = param
			pi++
		}

		if i >= len(pattern)-1 {
			child.handlerFunc = handlerFunc
		}

		if exist {
			parent = child
			continue
		}

		parent = parent.newChild(child, edge)
	}

	if pi > m.maxParam {
		m.maxParam = pi
	}
}

func (m *Mux) lookup(r *http.Request) (http.HandlerFunc, *Context) {
	s := r.URL.Path
	var parent, child *node

	if parent = m.tree.child[r.Method]; parent == nil {
		return nil, nil
	}

	//STATIC PATH
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
					ctx = m.pool.Get().(*Context)
				}
				ctx.params.Set(child.param, edge)

			} else if child = parent.child[wildcard]; child == nil {
				//BACKTRACK
				if child = parent.parent.child[colon]; child != nil {
					if ctx == nil {
						ctx = m.pool.Get().(*Context)
					}
					ctx.params.Set(child.param, s[bsi:si-1])
					si = bsi

				} else if child = parent.parent.child[wildcard]; child != nil {
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

			ctx.params = ctx.params[:0]
			m.pool.Put(ctx)
			return
		}

		fn.ServeHTTP(w, r)
		return
	}

	m.NotFound.ServeHTTP(w, r)
}
