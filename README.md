# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           530 ns/op         384 B/op          5 allocs/op
Param5           1000000          1124 ns/op         832 B/op          8 allocs/op
Param20           500000          2727 ns/op        2368 B/op         10 allocs/op
ParamWrite       2000000           648 ns/op         400 B/op          6 allocs/op
GithubStatic    30000000          48.1 ns/op           0 B/op          0 allocs/op
GithubParam      2000000           761 ns/op         448 B/op          6 allocs/op
GithubAll          10000        141791 ns/op       78272 B/op       1005 allocs/op
GPlusStatic     50000000          32.3 ns/op           0 B/op          0 allocs/op
GPlusParam       3000000           538 ns/op         384 B/op          5 allocs/op
GPlus2Params     2000000           730 ns/op         448 B/op          6 allocs/op
GPlusAll          200000          7269 ns/op        4544 B/op         60 allocs/op
ParseStatic     30000000          40.4 ns/op           0 B/op          0 allocs/op
ParseParam       3000000           558 ns/op         384 B/op          5 allocs/op
Parse2Params     2000000           716 ns/op         448 B/op          6 allocs/op
ParseAll          200000         10719 ns/op        6336 B/op         83 allocs/op
StaticAll         200000          8460 ns/op           0 B/op          0 allocs/op
```