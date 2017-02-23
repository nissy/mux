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
Param            1000000          1159 ns/op         512 B/op         10 allocs/op
Param5            500000          2721 ns/op        1152 B/op         22 allocs/op
Param20           200000          8320 ns/op        3696 B/op         56 allocs/op
ParamWrite       1000000          1284 ns/op         528 B/op         11 allocs/op
GithubStatic    30000000          39.3 ns/op           0 B/op          0 allocs/op
GithubParam       500000          2565 ns/op        1024 B/op         17 allocs/op
GithubAll           3000        509263 ns/op      155584 B/op       2479 allocs/op
GPlusStatic     50000000          28.6 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1727 ns/op         712 B/op         13 allocs/op
GPlus2Params      500000          3143 ns/op        1072 B/op         18 allocs/op
GPlusAll           50000         27732 ns/op        9448 B/op        152 allocs/op
ParseStatic     50000000          31.8 ns/op           0 B/op          0 allocs/op
ParseParam       1000000          1660 ns/op         640 B/op         11 allocs/op
Parse2Params     1000000          2208 ns/op         744 B/op         15 allocs/op
ParseAll           50000         33771 ns/op       11704 B/op        207 allocs/op
StaticAll         200000          7837 ns/op           0 B/op          0 allocs/op
```