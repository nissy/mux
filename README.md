# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           393 ns/op         320 B/op          3 allocs/op
Param5           3000000           537 ns/op         320 B/op          3 allocs/op
Param20          2000000           976 ns/op         320 B/op          3 allocs/op
ParamWrite       3000000           506 ns/op         336 B/op          4 allocs/op
GithubStatic    30000000          49.2 ns/op           0 B/op          0 allocs/op
GithubParam      3000000           503 ns/op         320 B/op          3 allocs/op
GithubAll          20000         95186 ns/op       53443 B/op        501 allocs/op
GPlusStatic     50000000          32.9 ns/op           0 B/op          0 allocs/op
GPlusParam       3000000           412 ns/op         320 B/op          3 allocs/op
GPlus2Params     3000000           471 ns/op         320 B/op          3 allocs/op
GPlusAll          300000          5548 ns/op        3520 B/op         33 allocs/op
ParseStatic     50000000          40.2 ns/op           0 B/op          0 allocs/op
ParseParam       3000000           467 ns/op         320 B/op          3 allocs/op
Parse2Params     3000000           506 ns/op         320 B/op          3 allocs/op
ParseAll          200000          8760 ns/op        5120 B/op         48 allocs/op
StaticAll         200000          9026 ns/op           0 B/op          0 allocs/op
```
