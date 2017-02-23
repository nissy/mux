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
Param            2000000           854 ns/op         384 B/op          5 allocs/op
Param5           1000000          1765 ns/op         832 B/op          8 allocs/op
Param20           300000          4992 ns/op        2368 B/op         10 allocs/op
ParamWrite       1000000          1034 ns/op         400 B/op          6 allocs/op
GithubStatic    30000000          44.8 ns/op           0 B/op          0 allocs/op
GithubParam      1000000          1664 ns/op         448 B/op          6 allocs/op
GithubAll           5000        312488 ns/op       78272 B/op       1005 allocs/op
GPlusStatic     100000000         25.3 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1010 ns/op         384 B/op          5 allocs/op
GPlus2Params     1000000          1653 ns/op         448 B/op          6 allocs/op
GPlusAll          100000         16790 ns/op        4544 B/op         60 allocs/op
ParseStatic     50000000          30.6 ns/op           0 B/op          0 allocs/op
ParseParam       1000000          1149 ns/op         384 B/op          5 allocs/op
Parse2Params     1000000          1497 ns/op         448 B/op          6 allocs/op
ParseAll          100000         23220 ns/op        6336 B/op         83 allocs/op
StaticAll         200000          8053 ns/op           0 B/op          0 allocs/op
```