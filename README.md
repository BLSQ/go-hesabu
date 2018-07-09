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
```

or via the piped version

```
cat ./test/small.json | go run hesabu.go
```

you will get

```
go run hesabu.go ./test/small.json
cat ./test/small.json | go run hesabu.go
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

# build the cli

go build -ldflags="-s -w" -o hesabucli hesabu.go; mv hesabucli ../hesabu/bin
