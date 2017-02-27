# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           466 ns/op         384 B/op          5 allocs/op
Param5           2000000           648 ns/op         512 B/op          5 allocs/op
Param20          1000000          1357 ns/op         992 B/op          5 allocs/op
ParamWrite       3000000           584 ns/op         400 B/op          6 allocs/op
GithubStatic    50000000          44.6 ns/op           0 B/op          0 allocs/op
GithubParam      2000000           655 ns/op         480 B/op          5 allocs/op
GithubAll          10000        122888 ns/op       78560 B/op        835 allocs/op
GPlusStatic     50000000          32.7 ns/op           0 B/op          0 allocs/op
GPlusParam       3000000           526 ns/op         416 B/op          5 allocs/op
GPlus2Params     3000000           590 ns/op         416 B/op          5 allocs/op
GPlusAll          200000          6504 ns/op        4576 B/op         55 allocs/op
ParseStatic     30000000          39.0 ns/op           0 B/op          0 allocs/op
ParseParam       3000000           529 ns/op         416 B/op          5 allocs/op
Parse2Params     3000000           557 ns/op         416 B/op          5 allocs/op
ParseAll          200000          9705 ns/op        6656 B/op         80 allocs/op
StaticAll         200000          8821 ns/op           0 B/op          0 allocs/op
```
