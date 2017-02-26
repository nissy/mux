# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           565 ns/op         416 B/op          5 allocs/op
Param5           2000000           731 ns/op         512 B/op          5 allocs/op
Param20          1000000          1467 ns/op         992 B/op          5 allocs/op
ParamWrite       2000000           694 ns/op         432 B/op          6 allocs/op
GithubStatic    30000000          46.6 ns/op           0 B/op          0 allocs/op
GithubParam      2000000           746 ns/op         480 B/op          5 allocs/op
GithubAll          10000        143790 ns/op       81408 B/op        835 allocs/op
GPlusStatic     50000000          32.5 ns/op           0 B/op          0 allocs/op
GPlusParam       2000000           589 ns/op         416 B/op          5 allocs/op
GPlus2Params     2000000           707 ns/op         480 B/op          5 allocs/op
GPlusAll          200000          7565 ns/op        4960 B/op         55 allocs/op
ParseStatic     50000000          40.5 ns/op           0 B/op          0 allocs/op
ParseParam       2000000           622 ns/op         448 B/op          5 allocs/op
Parse2Params     2000000           672 ns/op         480 B/op          5 allocs/op
ParseAll          100000         11185 ns/op        7264 B/op         80 allocs/op
StaticAll         200000          8560 ns/op           0 B/op          0 allocs/op
```