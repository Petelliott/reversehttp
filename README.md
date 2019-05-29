# reversehttp

reversehttp is a reverse http server/client library based on
[this rfc draft](https://tools.ietf.org/html/draft-lentczner-rhttp-00)

[![Build Status](https://travis-ci.com/Petelliott/reversehttp.svg?branch=master)](https://travis-ci.com/Petelliott/reversehttp)
[![Coverage Status](https://coveralls.io/repos/github/Petelliott/reversehttp/badge.svg?branch=master)](https://coveralls.io/github/Petelliott/reversehttp?branch=master)
[![GoDoc](https://godoc.org/github.com/Petelliott/reversehttp?status.svg)](https://godoc.org/github.com/Petelliott/reversehttp)
[![Go Report Card](https://goreportcard.com/badge/github.com/petelliott/reversehttp)](https://goreportcard.com/report/github.com/petelliott/reversehttp)

## installation

```
go get github.com/Petelliott/reversehttp
```

## example

In actual use you should check errors, but those checks have been removed for brevity

### server

```go
http.HandleFunc("/ptth", func(w http.ResponseWriter, r *http.Request) {
	client, _ := ReverseRequest(w, r)
	resp, _ := c.Get("http://example.com/path")
	// do whatever you want with the response
})
```

### client

```go
err := ReverseFunc("http://example.com/ptth", func(w http.ResponseWriter, r *http.Request) {
	// this could be any http.Handler
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("hello world\n"))
})
```

## notes

for clarity,

server shall refer to the http server, but reverse http client

client shall refer to the http client, but reverse http server
