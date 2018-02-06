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
sugar.Get("http://api.example.com/books/:id", P{"id": 123})
```

### Query
```go
sugar.Get("http://api.example.com/books", Query{"name": "bookA"})
sugar.Get("http://api.example.com/books", Q{"name": "bookA"})
```

### Cookie
```go
sugar.Get("http://api.example.com/books", Cookie{"name": "sugar"})
sugar.Get("http://api.example.com/books", C{"name": "sugar"})
```

### Header
```go
sugar.Get("http://api.example.com/books", Header{"name": "bookA"})
sugar.Get("http://api.example.com/books", H{"name": "bookA"})
```

### JSON
```go
//will automatically add 'Content-Type=application/json;charset=UTF-8' to header
sugar.Post("http://api.example.com/books", Json(`{"Name":"bookA"}`))
sugar.Post("http://api.example.com/books", J(`{"Name":"bookA"}`))
```
