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
Param            2000000           751 ns/op         400 B/op          6 allocs/op
Param5           1000000          1435 ns/op         896 B/op          9 allocs/op
Param20           500000          2769 ns/op        2672 B/op         11 allocs/op
ParamWrite       2000000           853 ns/op         416 B/op          7 allocs/op
GithubStatic    20000000          68.4 ns/op           0 B/op          0 allocs/op
GithubParam      1000000          1926 ns/op         880 B/op          9 allocs/op
GithubAll           5000        350965 ns/op      213696 B/op       1444 allocs/op
GPlusStatic     20000000          72.6 ns/op           0 B/op          0 allocs/op
GPlusParam       1000000          1086 ns/op         400 B/op          6 allocs/op
GPlus2Params     1000000          1933 ns/op         624 B/op          8 allocs/op
GPlusAll          100000         16311 ns/op        5168 B/op         73 allocs/op
ParseStatic     20000000          71.9 ns/op           0 B/op          0 allocs/op
ParseParam       1000000          1050 ns/op         416 B/op          6 allocs/op
Parse2Params     1000000          1284 ns/op         496 B/op          7 allocs/op
ParseAll          100000         18281 ns/op        6896 B/op         99 allocs/op
StaticAll         100000         13640 ns/op           0 B/op          0 allocs/op
```