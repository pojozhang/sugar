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
sugar.Get("http://api.example.com/books", Cookie{"name": "bookA"})
sugar.Get("http://api.example.com/books", C{"name": "bookA"})
```

### Header
```go
sugar.Get("http://api.example.com/books", Header{"name": "bookA"})
sugar.Get("http://api.example.com/books", H{"name": "bookA"})
```

### Json
```go
//will automatically add 'Content-Type=application/json;charset=UTF-8' to header if 'Content-Type' not exists
sugar.Post("http://api.example.com/books", Json(`{"name":"bookA"}`))
sugar.Post("http://api.example.com/books", J(`{"name":"bookA"}`))
```

### Form
```go
sugar.Post("http://api.example.com/books", Form{"name": "bookA"})
sugar.Post("http://api.example.com/books", F{"name": "bookA"})
```
