# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            5000000           412 ns/op         320 B/op          3 allocs/op
Param5           3000000           514 ns/op         320 B/op          3 allocs/op
Param20          1000000          1122 ns/op         320 B/op          3 allocs/op
ParamWrite       3000000           545 ns/op         336 B/op          4 allocs/op
GithubStatic    30000000          46.5 ns/op           0 B/op          0 allocs/op
GithubParam      3000000           515 ns/op         320 B/op          3 allocs/op
GithubAll          10000        104713 ns/op       53443 B/op        501 allocs/op
GPlusStatic     50000000          33.5 ns/op           0 B/op          0 allocs/op
GPlusParam       3000000           516 ns/op         320 B/op          3 allocs/op
GPlus2Params     3000000           510 ns/op         320 B/op          3 allocs/op
GPlusAll          300000          5283 ns/op        3520 B/op         33 allocs/op
ParseStatic     30000000          37.1 ns/op           0 B/op          0 allocs/op
ParseParam       3000000           426 ns/op         320 B/op          3 allocs/op
Parse2Params     3000000           469 ns/op         320 B/op          3 allocs/op
ParseAll          200000          8277 ns/op        5120 B/op         48 allocs/op
StaticAll         200000          9441 ns/op           0 B/op          0 allocs/op
```
