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
Param            1000000          1393 ns/op         640 B/op         11 allocs/op
Param5            500000          2765 ns/op        1152 B/op         22 allocs/op
Param20           200000          8096 ns/op        3744 B/op         57 allocs/op
ParamWrite       1000000          1523 ns/op         656 B/op         12 allocs/op
GithubStatic    20000000          62.6 ns/op           0 B/op          0 allocs/op
GithubParam       500000          2650 ns/op        1072 B/op         18 allocs/op
GithubAll           3000        478909 ns/op      167024 B/op       2587 allocs/op
GPlusStatic     30000000          45.6 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1700 ns/op         712 B/op         13 allocs/op
GPlus2Params      500000          2994 ns/op        1136 B/op         19 allocs/op
GPlusAll           50000         26723 ns/op       10248 B/op        160 allocs/op
ParseStatic     20000000          59.3 ns/op         0 B/op            0 allocs/op
ParseParam       1000000          1634 ns/op         640 B/op         11 allocs/op
Parse2Params     1000000          2429 ns/op        1000 B/op         16 allocs/op
ParseAll           50000         34781 ns/op       13384 B/op        216 allocs/op
StaticAll         100000         14135 ns/op         528 B/op         11 allocs/op
```