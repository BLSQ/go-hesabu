# go-hesabu
go-hesabu is golang equivalent to hesabu

# Development setup

require a go 1.9 and dep

```
cd ~/go/src/github.com/BLSQ
git clone git@github.com:BLSQ/go-hesabu.git
cd go-hesabu
dep ensure
go run hesabu.go ./test/small.json
```

For more info see workspace structure in https://golang.org/doc/code.html
# Usage
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
# Development
## build the cli

go build -ldflags="-s -w" -o hesabucli hesabu.go; mv hesabucli ../hesabu/bin

## Dependencies

relies on
 - [govaluate](https://github.com/Knetic/govaluate)
 - [toposort](https://github.com/otaviokr/topological-sort)

## License

The code is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
