# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           555 ns/op         384 B/op          5 allocs/op
Param5           1000000          1232 ns/op         832 B/op          8 allocs/op
Param20           500000          2885 ns/op        2368 B/op         10 allocs/op
ParamWrite       2000000           901 ns/op         400 B/op          6 allocs/op
GithubStatic    30000000          50.1 ns/op           0 B/op          0 allocs/op
GithubParam      1000000          1091 ns/op         448 B/op          6 allocs/op
GithubAll          10000        167884 ns/op       78272 B/op       1005 allocs/op
GPlusStatic     50000000          30.4 ns/op           0 B/op          0 allocs/op
GPlusParam       2000000           600 ns/op         384 B/op          5 allocs/op
GPlus2Params     2000000           794 ns/op         448 B/op          6 allocs/op
GPlusAll          200000          7587 ns/op        4544 B/op         60 allocs/op
ParseStatic     50000000          31.7 ns/op           0 B/op          0 allocs/op
ParseParam       2000000           586 ns/op         384 B/op          5 allocs/op
Parse2Params     2000000           763 ns/op         448 B/op          6 allocs/op
ParseAll          200000         11834 ns/op        6336 B/op         83 allocs/op
StaticAll         200000          8003 ns/op           0 B/op          0 allocs/op
```