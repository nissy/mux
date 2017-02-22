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
Param            1000000          1506 ns/op         640 B/op         11 allocs/op
Param5            500000          3020 ns/op        1152 B/op         22 allocs/op
Param20           200000          8773 ns/op        3744 B/op         57 allocs/op
ParamWrite       1000000          1651 ns/op         656 B/op         12 allocs/op
GithubStatic    20000000          66.1 ns/op           0 B/op          0 allocs/op
GithubParam       500000          2846 ns/op        1072 B/op         18 allocs/op
GithubAll           3000        476908 ns/op      167024 B/op       2587 allocs/op
GPlusStatic     30000000          47.1 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1776 ns/op         712 B/op         13 allocs/op
GPlus2Params      500000          2996 ns/op        1136 B/op         19 allocs/op
GPlusAll           50000         27778 ns/op       10248 B/op        160 allocs/op
ParseStatic     20000000          62.3 ns/op           0 B/op          0 allocs/op
ParseParam       1000000          1653 ns/op         640 B/op         11 allocs/op
Parse2Params     1000000          2535 ns/op        1000 B/op         16 allocs/op
ParseAll           50000         35221 ns/op       13384 B/op        216 allocs/op
StaticAll         100000         16113 ns/op         528 B/op         11 allocs/op
```