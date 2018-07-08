# go-hesabu
go-hesabu is golang equivalent to hesabu

```
go get -u github.com/otaviokr/topological-sort
go get -u github.com/Knetic/govaluate

# to get ide debugger
go get -u github.com/derekparker/delve/cmd/dlv
```

# From the command line
Taking equations

```json
{
  "c": "a + 10 * b",
  "b": "10+a",
  "a": "10"
}

```

logs the solution
```
go run hesabu.go ./test/small.json
2018/07/07 00:07:58 sum1 = 10.54 (SUM(10,.54))
2018/07/07 00:07:58 sum2 = 10.54 (sum(10,0.54))
2018/07/07 00:07:58 max1 = 100 (max(10,100))
2018/07/07 00:07:58 min2 = 10 (MIN(10,100))
2018/07/07 00:07:58 min1 = 10 (min(10,100))
2018/07/07 00:07:58 a = 10 (10)
2018/07/07 00:07:58 b = 20 (10+a)
2018/07/07 00:07:58 c = 210 (a + 10 * b)
2018/07/07 00:07:58 min3 = 10 (Min(10,100))
2018/07/07 00:07:58 sum3 = 64 (Sum(10,54))
2018/07/07 00:07:58 max3 = 100 (Max(10,100))
2018/07/07 00:07:58 max2 = 100 (MAX(10,100))
{
  "a": 10,
  "b": 20,
  "c": 210,
  "max1": 100,
  "max2": 100,
  "max3": 100,
  "min1": 10,
  "min2": 10,
  "min3": 10,
  "sum1": 10.54,
  "sum2": 10.54,
  "sum3": 64
}

```

# As a rest service

start the server

```
go run hesabu.go s
```

and post problem and get solutions

```
time curl -H "Content-Type: application/json"  --data @test/small.json http://127.0.0.1:8080
{
  "a": 10,
  "avg1": 2.5,
  "b": 20,
  "c": 210,
  "max1": 100,
  "max2": 100,
  "max3": 100,
  "min1": 10,
  "min2": 10,
  "min3": 10,
  "sum1": 10.54,
  "sum2": 10.54,
  "sum3": 64
}
real	0m0.009s
user	0m0.004s
sys	0m0.000s

```