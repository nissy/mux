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
	m := mux.New()

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
Param            1000000          1033 ns/op         400 B/op          7 allocs/op
Param5           1000000          2352 ns/op         912 B/op         18 allocs/op
Param20           200000          7333 ns/op        2688 B/op         50 allocs/op
ParamWrite       1000000          1150 ns/op         416 B/op          8 allocs/op
GithubStatic    30000000          42.2 ns/op           0 B/op          0 allocs/op
GithubParam      1000000          1982 ns/op         528 B/op         12 allocs/op
GithubAll           5000        369240 ns/op       83888 B/op       1691 allocs/op
GPlusStatic     50000000          25.3 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1238 ns/op         472 B/op          9 allocs/op
GPlus2Params     1000000          2027 ns/op         576 B/op         13 allocs/op
GPlusAll          100000         19739 ns/op        5016 B/op        101 allocs/op
ParseStatic     50000000          31.0 ns/op           0 B/op          0 allocs/op
ParseParam       1000000          1135 ns/op         400 B/op          7 allocs/op
Parse2Params     1000000          1645 ns/op         504 B/op         11 allocs/op
ParseAll           50000         24304 ns/op        7096 B/op        140 allocs/op
StaticAll         200000          8317 ns/op           0 B/op          0 allocs/op
```