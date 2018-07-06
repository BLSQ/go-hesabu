# go-hesabu
go-hesabu is golang equivalent to hesabu

```
go get -u github.com/philopon/go-toposort
go get -u github.com/Knetic/govaluate

# to get ide debugger
go get -u github.com/derekparker/delve/cmd/dlv
```

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
go run hello.go
2018/07/06 22:26:21 equations loaded: 3 
2018/07/06 22:26:21 vars  [a b]
2018/07/06 22:26:21 vars  [a]
2018/07/06 22:26:21 vars  []
[c b a]
[a b c]
2018/07/06 22:26:21 a = 10 (10)
2018/07/06 22:26:21 b = 20 (10+a)
2018/07/06 22:26:21 c = 210 (a + 10 * b)
```
