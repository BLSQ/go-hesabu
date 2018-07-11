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
{
  "a": 10,
  "b": 20,
  "c": 210,
}

```

# build the cli

go build -ldflags="-s -w" -o hesabucli hesabu.go; mv hesabucli ../hesabu/bin

# Dependencies

  [govaluate](https://github.com/Knetic/govaluate)
	[toposort](https://github.com/otaviokr/topological-sort)
  
