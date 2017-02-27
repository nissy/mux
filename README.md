# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           539 ns/op         384 B/op          5 allocs/op
Param5           2000000           727 ns/op         512 B/op          5 allocs/op
Param20          1000000          1518 ns/op         992 B/op          5 allocs/op
ParamWrite       2000000           664 ns/op         400 B/op          6 allocs/op
GithubStatic    30000000          49.4 ns/op           0 B/op          0 allocs/op
GithubParam      2000000           751 ns/op         480 B/op          5 allocs/op
GithubAll          10000        141625 ns/op       78560 B/op        835 allocs/op
GPlusStatic     50000000          33.1 ns/op           0 B/op          0 allocs/op
GPlusParam       2000000           618 ns/op         416 B/op          5 allocs/op
GPlus2Params     2000000           693 ns/op         416 B/op          5 allocs/op
GPlusAll          200000          7684 ns/op        4576 B/op         55 allocs/op
ParseStatic     50000000          37.9 ns/op           0 B/op          0 allocs/op
ParseParam       2000000           607 ns/op         416 B/op          5 allocs/op
Parse2Params     2000000           670 ns/op         416 B/op          5 allocs/op
ParseAll          100000         11569 ns/op        6656 B/op         80 allocs/op
StaticAll         200000          8664 ns/op           0 B/op          0 allocs/op
```
