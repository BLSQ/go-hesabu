[![Maintainability](https://api.codeclimate.com/v1/badges/521737120ca70381247d/maintainability)](https://codeclimate.com/github/BLSQ/go-hesabu/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/BLSQ/go-hesabu)](https://goreportcard.com/report/github.com/BLSQ/go-hesabu)
[![Build Status](https://travis-ci.org/BLSQ/go-hesabu.svg?branch=master)](https://travis-ci.org/BLSQ/go-hesabu)

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

require a go 1.9 and dep

```
mkdir -p $GOPATH/src/github.com/BLSQ
cd $GOPATH/src/github.com/BLSQ
git clone git@github.com:BLSQ/go-hesabu.git
cd go-hesabu
dep ensure
go test github.com/BLSQ/go-hesabu/hesabu -coverprofile=coverage.out
go tool cover -html=coverage.out
go run hesabu.go ./test/small.json
```

For more info see workspace structure in https://golang.org/doc/code.html

## build the cli

```
go build -ldflags="-s -w" -o hesabucli hesabu.go; mv hesabucli ../hesabu/bin
```

for mac

```
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o hesabucli hesabu.go
```

## releasing

have goreleaser https://goreleaser.com/install/

on ubuntu : `snap install goreleaser --classic`

for github token, use scope repo and see : https://github.com/settings/tokens/new

```
export GITHUB_TOKEN=...
rm -rf ./dist/
git tag -a v0.0.3 -m "First release with goreleaser"
git push origin v0.0.3
goreleaser release
```

check the release page https://github.com/BLSQ/go-hesabu/releases/


## Dependencies

relies on
 - [govaluate](https://github.com/Knetic/govaluate)
 - [toposort](https://github.com/otaviokr/topological-sort)

## License

The code is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
