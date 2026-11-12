[![Go](https://github.com/BLSQ/go-hesabu/actions/workflows/go.yml/badge.svg)](https://github.com/BLSQ/go-hesabu/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/BLSQ/go-hesabu)](https://goreportcard.com/report/github.com/BLSQ/go-hesabu)

# go-hesabu
go-hesabu is golang equivalent to hesabu


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

## Development setup

requires the go version of `go.mod` (1.27)

```
./script/test && go tool cover -func=coverage.out && go tool cover  -html=coverage.out -o coverage.html
go run hesabu.go ./test/small.json
```

For more info see workspace structure in https://golang.org/doc/code.html

## build the cli

### if you don't have golang setup

you can use docker

```
script/docker-build
```

Which will generate both a Mac version and a Linux version in the bin folder.

### if you have golang

```
script/build
```

Which will generate both a Mac version and a Linux version in the bin folder.

## Profiling

You can run the binary with two additional flags, to enable cpu and memory profiling.

`bin/hesabucli -cpuprofile cpu.prof -memprofile mem.prof test/large_set_of_equations.json`

Now you should have both a cpu.prof and a mem.prof in your directory.

If you now run:

      go tool pprof cpu.prof

You'll be dropped into an interactive profile viewer.

Even more interesting is that you can open an interactive version in you browser with the following:

      pprof -http=localhost:1234 bin/hesabucli cpu.prof

## releasing

Push a tag, the `Release` workflow tests, builds `hesabucli` (linux, static) and
`hesabucli-mac` with the go version of `go.mod`, and publishes them with their
checksums as a github release:

```
git tag -a v0.1.0 -m "go 1.27"
git push origin v0.1.0
```

check the release page https://github.com/BLSQ/go-hesabu/releases/ then update
the binaries of the hesabu gem from it (see the hesabu README).

Local build (same flags, reproducible): `script/build`, or without go installed
`script/docker-build`.


## Dependencies

relies on
 - [govaluate](https://github.com/Knetic/govaluate)
 - [toposort](https://github.com/otaviokr/topological-sort)

## License

The code is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
