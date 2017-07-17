# ways2go

2 Way SQL implementation for golang

[![Build Status](https://travis-ci.org/mattn/ways2go.png?branch=master)](https://travis-ci.org/mattn/ways2go)
[![GoDoc](https://godoc.org/github.com/mattn/ways2go?status.svg)](http://godoc.org/github.com/mattn/ways2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattn/ways2go)](https://goreportcard.com/report/github.com/mattn/ways2go)

# Usage

```go
got, err := Eval(`
select * from foo where id = /*id*/5 /* IF enabled */and bar = /*bar*/ /*END*/
`, map[string]{"enabled": true}, ways2go.Question)
```

# License

MIT

# Author

Yasuhiro Matsumoto (a.k.a. mattn)
