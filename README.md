# Sugar  [![Build Status](https://travis-ci.org/pojozhang/sugar.svg?branch=master)](https://travis-ci.org/pojozhang/sugar) [![codecov](https://codecov.io/gh/pojozhang/sugar/branch/master/graph/badge.svg)](https://codecov.io/gh/pojozhang/sugar)

Sugar is a **DECLARATIVE** http client providing elegant APIs for Golang.

Now you can send requests in just one line.


## Set Up
```bash
dep ensure -add github.com/pojozhang/sugar
```

## Usage

### Path
```go
// GET /books/123 HTTP/1.1
// Host: api.example.com
sugar.Get("http://api.example.com/books/:id", Path{"id": 123})
sugar.Get("http://api.example.com/books/:id", P{"id": 123})
```

### Query
```go
// GET /books?name=bookA HTTP/1.1
// Host: api.example.com
sugar.Get("http://api.example.com/books", Query{"name": "bookA"})
sugar.Get("http://api.example.com/books", Q{"name": "bookA"})

// list
// GET /books?name=bookA&name=bookB HTTP/1.1
// Host: api.example.com
sugar.Get("http://api.example.com/books", Query{"name": List{"bookA", "bookB"}})
sugar.Get("http://api.example.com/books", Q{"name": L{"bookA", "bookB"}})
```

### Cookie
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Cookie: name=bookA
sugar.Get("http://api.example.com/books", Cookie{"name": "bookA"})
sugar.Get("http://api.example.com/books", C{"name": "bookA"})
```

### Header
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Name: bookA
sugar.Get("http://api.example.com/books", Header{"name": "bookA"})
sugar.Get("http://api.example.com/books", H{"name": "bookA"})
```

### Json
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/json;charset=UTF-8
// {"name":"bookA"}
// automatically set 'Content-Type=application/json;charset=UTF-8' if 'Content-Type' not exists
sugar.Post("http://api.example.com/books", Json(`{"name":"bookA"}`))
sugar.Post("http://api.example.com/books", J(`{"name":"bookA"}`))

// map
sugar.Post("http://api.example.com/books", Json(Map{"name": "bookA"}))
sugar.Post("http://api.example.com/books", J(M{"name": "bookA"}))

// list
sugar.Post("http://api.example.com/books", Json(List{Map{"name": "bookA"}}))
sugar.Post("http://api.example.com/books", J(L{M{"name": "bookA"}}))
```

### Form
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/x-www-form-urlencoded
sugar.Post("http://api.example.com/books", Form{"name": "bookA"})
sugar.Post("http://api.example.com/books", F{"name": "bookA"})

// list
sugar.Post("http://api.example.com/books", Form{"name": List{"bookA", "bookB"}})
sugar.Post("http://api.example.com/books", F{"name": L{"bookA", "bookB"}})
```

### Basic Auth
```go
// DELETE /books HTTP/1.1
// Host: api.example.com
// Authorization: Basic dXNlcjpwYXNzd29yZA==
sugar.Delete("http://api.example.com/books", User{"user", "password"})
sugar.Delete("http://api.example.com/books", U{"user", "password"})
```

### Multipart
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: multipart/form-data; boundary=19b8acc2469f1914a24fc6e0152aac72f1f92b6f5104b57477262816ab0f
//
// --19b8acc2469f1914a24fc6e0152aac72f1f92b6f5104b57477262816ab0f
// Content-Disposition: form-data; name="name"
//
// bookA
// --19b8acc2469f1914a24fc6e0152aac72f1f92b6f5104b57477262816ab0f
// Content-Disposition: form-data; name="file"; filename="text"
// Content-Type: application/octet-stream
//
// hello sugar!
// --19b8acc2469f1914a24fc6e0152aac72f1f92b6f5104b57477262816ab0f--
sugar.Post("http://api.example.com/books", MultiPart{"name": "bookA", "file": File("text")})
sugar.Post("http://api.example.com/books", D{"name": "bookA", "file": File("text")})

// we also support *os.File
f, _ := os.Open("text")
sugar.Post("http://api.example.com/books", MultiPart{"name": "bookA", "file": f})
```

### Mix
Due to Sugar's flexible design, different types of parameters can be freely combined.

```go
sugar.Patch("http://api.example.com/books/:id", Path{"id": 123}, Json(`{"name":"bookA"}`), User{"user", "password"})
```

### Apply
You can use Apply() to preset some values which will be attached to every following request.

```go
sugar.Apply(User{"user", "password"})
sugar.Get("http://api.example.com/books")
sugar.Get("http://api.example.com/books")
```
```go
sugar.Get("http://api.example.com/books", User{"user", "password"})
sugar.Get("http://api.example.com/books", User{"user", "password"})
```
The latter is equal to the former.

### Extension
You can register your custom resolver which should implement interface `Resolver` and bind it to the target type.  
```go
sugar.Register(Custom{}, &CustomResolver{})
sugar.Get("http://api.example.com/books", Custom{})
```