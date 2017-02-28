# Patricia trie mux
Go http request multiplexer

### Benchmark

https://github.com/ngc224/go-http-routing-benchmark

```
Param            3000000           423 ns/op         320 B/op          3 allocs/op
Param5           3000000           505 ns/op         320 B/op          3 allocs/op
Param20          2000000           899 ns/op         320 B/op          3 allocs/op
ParamWrite       3000000           522 ns/op         336 B/op          4 allocs/op
GithubStatic    30000000          48.1 ns/op           0 B/op          0 allocs/op
GithubParam      3000000           525 ns/op         320 B/op          3 allocs/op
GithubAll          20000         99943 ns/op       53443 B/op        501 allocs/op
GPlusStatic     50000000          32.5 ns/op           0 B/op          0 allocs/op
GPlusParam       3000000           438 ns/op         320 B/op          3 allocs/op
GPlus2Params     3000000           513 ns/op         320 B/op          3 allocs/op
GPlusAll          300000          6319 ns/op        3520 B/op         33 allocs/op
ParseStatic     30000000          40.5 ns/op           0 B/op          0 allocs/op
ParseParam       3000000           460 ns/op         320 B/op          3 allocs/op
Parse2Params     3000000           479 ns/op         320 B/op          3 allocs/op
ParseAll          200000          8105 ns/op        5120 B/op         48 allocs/op
StaticAll         200000          9120 ns/op           0 B/op          0 allocs/op
```
