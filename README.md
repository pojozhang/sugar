# Sugar  [![Build Status](https://travis-ci.org/pojozhang/sugar.svg?branch=master)](https://travis-ci.org/pojozhang/sugar)

A simple http client with elegant APIs for Golang.

Now you can send requests in just one line.


## Set Up
```bash
dep ensure -add github.com/pojozhang/sugar
```

## Usage

### Path
```go
sugar.Get("http://api.example.com/books/:id", Path{"id": 123})
```

### Query
```go
sugar.Get("http://api.example.com/books", Query{"name": "bookA"})
```

### Cookie
```go
sugar.Get("http://api.example.com/books", Cookie{"name": "sugar"})
```

### JSON
```go
sugar.Post("http://api.example.com/books", Json(`{"Name":"bookA"}`))
```
