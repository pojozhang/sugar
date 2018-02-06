# Sugar  [![Build Status](https://travis-ci.org/pojozhang/sugar.svg?branch=master)](https://travis-ci.org/pojozhang/sugar)

## Set Up
```bash
dep ensure -add github.com/pojozhang/sugar
```

## Usage

### Build a get request with path variables
```go
sugar.Get("http://api.example.com/books/:id", Path{"id": 123})
```
