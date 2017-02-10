# mux
Go http request multiplexer

### Example

```
package main

import (
	"net/http"

	"github.com/ngc224/mux"
)

func main() {
	m := mux.NewMux()

	m.Get("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Static"))
	})

	m.Get("/a/b/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Asterisk"))
	})

	m.Get("/a/b/c/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Param id is " + mux.URLParam(r, "id")))
	})

	http.ListenAndServe(":8080", m)
}
```

### Benchmark

```
Param                2000000               770 ns/op             416 B/op          7 allocs/op
Param5               1000000              1482 ns/op             912 B/op         10 allocs/op
Param20               500000              2795 ns/op            2688 B/op         12 allocs/op
ParamWrite           2000000               905 ns/op             432 B/op          8 allocs/op
GithubStatic         5000000               376 ns/op             336 B/op          4 allocs/op
GithubParam          1000000              1954 ns/op             896 B/op         10 allocs/op
GithubAll               5000            378866 ns/op          228464 B/op       1755 allocs/op
ParseStatic          5000000               390 ns/op             336 B/op          4 allocs/op
ParseParam           2000000               942 ns/op             432 B/op          7 allocs/op
Parse2Params         1000000              1242 ns/op             512 B/op          8 allocs/op
ParseAll               50000             25963 ns/op           10512 B/op        155 allocs/op
StaticAll              20000             83069 ns/op           52752 B/op        628 allocs/op
```