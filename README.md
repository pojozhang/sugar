# Sugar  [![Build Status](https://travis-ci.org/pojozhang/sugar.svg?branch=master)](https://travis-ci.org/pojozhang/sugar)

A simple http client with elegant APIs for Golang.

Now you can send requests in just one line.


## Set Up
```bash
dep ensure -add github.com/pojozhang/sugar
```

## Usage

### Build a get request with path variables
```go
sugar.Get("http://api.example.com/books/:id", Path{"id": 123})
```

### Post JSON
```go
sugar.Post("http://api.example.com/books", Json(`{"Name":"bookA"}`)
```
