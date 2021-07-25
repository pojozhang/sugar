<img align="middle" height="200px" src="logo.png">

![GitHub (pre-)release](https://img.shields.io/github/release/pojozhang/sugar/all.svg)
[![Go](https://github.com/pojozhang/sugar/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/pojozhang/sugar/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/pojozhang/sugar/branch/master/graph/badge.svg)](https://codecov.io/gh/pojozhang/sugar) [![Go Report Card](https://goreportcard.com/badge/github.com/pojozhang/sugar)](https://goreportcard.com/report/github.com/pojozhang/sugar) ![go](https://img.shields.io/badge/golang-1.13+-blue.svg) [![GoDoc](https://godoc.org/github.com/pojozhang/sugar?status.svg)](https://godoc.org/github.com/pojozhang/sugar) 
![license](https://img.shields.io/github/license/pojozhang/sugar.svg)

Sugar is a **DECLARATIVE** http client providing elegant APIs for Golang.

### [🇨🇳 中文文档](README.zh-cn.md)

## 🌈 Features
- Elegant APIs
- Plugins
- Chained invocations
- Highly extensible

## 🚀 Download
```bash
go get -add github.com/pojozhang/sugar
```

## 📙 Usage
Firstly you need to import the package.
```go
import . "github.com/pojozhang/sugar"
```
And now you are ready to easily send any request to any corner on this blue planet.

### Request
#### Plain Text
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: text/plain
Post(ctx, "http://api.example.com/books", "bookA")
```

#### Path
```go
// GET /books/123 HTTP/1.1
// Host: api.example.com
Get(ctx, "http://api.example.com/books/:id", Path{"id": 123})
Get(ctx, "http://api.example.com/books/:id", P{"id": 123})
```

#### Query
```go
// GET /books?name=bookA HTTP/1.1
// Host: api.example.com
Get(ctx, "http://api.example.com/books", Query{"name": "bookA"})
Get(ctx, "http://api.example.com/books", Q{"name": "bookA"})

// list
// GET /books?name=bookA&name=bookB HTTP/1.1
// Host: api.example.com
Get(ctx, "http://api.example.com/books", Query{"name": List{"bookA", "bookB"}})
Get(ctx, "http://api.example.com/books", Q{"name": L{"bookA", "bookB"}})
```

#### Cookie
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Cookie: name=bookA
Get(ctx, "http://api.example.com/books", Cookie{"name": "bookA"})
Get(ctx, "http://api.example.com/books", C{"name": "bookA"})
```

#### Header
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Name: bookA
Get(ctx, "http://api.example.com/books", Header{"name": "bookA"})
Get(ctx, "http://api.example.com/books", H{"name": "bookA"})
```

#### Json
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/json;charset=UTF-8
// {"name":"bookA"}
Post(ctx, "http://api.example.com/books", Json{`{"name":"bookA"}`})
Post(ctx, "http://api.example.com/books", J{`{"name":"bookA"}`})

// map
Post(ctx, "http://api.example.com/books", Json{Map{"name": "bookA"}})
Post(ctx, "http://api.example.com/books", J{M{"name": "bookA"}})

// list
Post(ctx, "http://api.example.com/books", Json{List{Map{"name": "bookA"}}})
Post(ctx, "http://api.example.com/books", J{L{M{"name": "bookA"}}})
```

#### Xml
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Authorization: Basic dXNlcjpwYXNzd29yZA==
// Content-Type: application/xml; charset=UTF-8
// <book name="bookA"></book>
Post(ctx, "http://api.example.com/books", Xml{`<book name="bookA"></book>`})
Post(ctx, "http://api.example.com/books", X{`<book name="bookA"></book>`})
```

#### Form
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/x-www-form-urlencoded
Post(ctx, "http://api.example.com/books", Form{"name": "bookA"})
Post(ctx, "http://api.example.com/books", F{"name": "bookA"})

// list
Post(ctx, "http://api.example.com/books", Form{"name": List{"bookA", "bookB"}})
Post(ctx, "http://api.example.com/books", F{"name": L{"bookA", "bookB"}})
```

#### Basic Auth
```go
// DELETE /books HTTP/1.1
// Host: api.example.com
// Authorization: Basic dXNlcjpwYXNzd29yZA==
Delete(ctx, "http://api.example.com/books", User{"user", "password"})
Delete(ctx, "http://api.example.com/books", U{"user", "password"})
```

#### Multipart
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
f, _ := os.Open("text")
Post(ctx, "http://api.example.com/books", MultiPart{"name": "bookA", "file": f})
Post(ctx, "http://api.example.com/books", MP{"name": "bookA", "file": f})
```

#### Mix
Due to Sugar's flexible design, different types of parameters can be freely combined.
```go
Patch(ctx, "http://api.example.com/books/:id", Path{"id": 123}, Json{`{"name":"bookA"}`}, User{"user", "password"})
```

#### Apply
You can use Apply() to preset some values which will be attached to every following request. Call Reset() to clean preset values.
```go
Apply(User{"user", "password"})
Get(ctx, "http://api.example.com/books")
Get(ctx, "http://api.example.com/books")
Reset()
Get(ctx, "http://api.example.com/books")
```
```go
Get(ctx, "http://api.example.com/books", User{"user", "password"})
Get(ctx, "http://api.example.com/books", User{"user", "password"})
Get(ctx, "http://api.example.com/books")
```
The latter is equal to the former.


### Response
A request API always returns a value of type `*Response` which also provides some nice APIs.

#### Raw
Raw() returns a value of type `*http.Response` and an `error` which is similar to standard go API.
```go
resp, err := Post(ctx, "http://api.example.com/books", "bookA").Raw()
...
```

#### ReadBytes
ReadBytes() is another syntax sugar to read bytes from response body.
Notice that this method will close body after reading.
```go
bytes, resp, err := Get(ctx, "http://api.example.com/books").ReadBytes()
...
```

#### Read
Read() reads different types of response via decoder API.
The following two examples read response body as plain text/JSON according to different content types.
```go
// plain text
var text = new(string)
resp, err := Get(ctx, "http://api.example.com/text").Read(text)

// json
var books []book
resp, err := Get(ctx, "http://api.example.com/json").Read(&books)
```

#### Download files
You can also use Read() to download files.
```go
f,_ := os.Create("tmp.png")
defer f.Close()
resp, err := Get(ctx, "http://api.example.com/logo.png").Read(f)
```

## 🔌 Extension
There are three major components in Sugar: **Encoder**, **Decoder** and **Plugin**.
- An encoder is used to encode your parameters and assemble requests.
- A decoder is used to decode the data from response body.
- A plugin is designed to work as an interceptor.

### Encoder
You can register your custom encoder which should implement `Encoder` interface.
```go
type MyEncoder struct {
}

func (r *MyEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
    myParams, ok := context.Param.(MyParam)
    if !ok {
	return chain.Next()
    }
    ...
    req := context.Request
    ...
    return nil
}

Encoders.Add(&MyEncoder{})

Get(ctx, "http://api.example.com/books", MyParam{})
```

### Decoder
You can implement `Decoder` interface so that you can convert a response body to a specific struct.
It is very convenient to get converted value via `func (r *Response) Read(v interface{}) (*http.Response, error)` API.
```go
type MyDecoder struct {
}

func (d *MyDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
    // decode data from body if a content type named `my-content-type` is set in header
    for _, contentType := range context.Response.Header[ContentType] {
	if strings.Contains(strings.ToLower(contentType), "my-content-type") {
	    body, err := ioutil.ReadAll(context.Response.Body)
	    if err != nil {
		return err
	    }
	    json.Unmarshal(body, context.Out)
	    ...
	    return nil
	}
    }
    return chain.Next()
}

Decoders.Add(&MyDecoder{})
```

### Plugin
Plugin is a new feature since V2. You can do anything as you like before the request is sent or after the response is received.
```go
// Implementation of builtin Logger plugin
Use(func(c *Context) error {
    b, _ := httputil.DumpRequest(c.Request, true)
    log.Println(string(b))
    defer func() {
        if c.Response != nil {
	    b, _ := httputil.DumpResponse(c.Response, true)
	    log.Println(string(b))
	}
    }()
    return c.Next()
})
```

#### Logger
You can use Logger plugin to log any request you send or any response you get.
```go
Use(Logger)
```

#### Retryer
You can use Retryer plugin to retry a request when the server returns 500 or when you get a net error.
```go
Use(Retryer(3, time.Second, 1, time.Second))
```

#### Custom error handling
Sometimes you may get an custom API error when a request is invalid. The following example shows you how to handle this situation via a plugin: 
```go
// your error struct
type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e apiError) Error() string {
	return e.Message
}

Use(func(c *Context) error {
		if err := c.Next(); err != nil {
			return err
		}

		if c.Response != nil && c.Response.StatusCode >= http.StatusBadRequest {
			defer func() { c.Response.Body.Close() }()
			body, err := ioutil.ReadAll(c.Response.Body)
			if err != nil {
				return err
			}

			e := apiError{}
			if err = json.Unmarshal(body, &e); err != nil {
				return err
			}
			return e
		}

		return nil
	})

// send request
_, err := client.Get(ctx, "some url").Read(&json{})
// type switch
switch e := err.(type) {
	case apiError:
		// your code
	}
// your code
```
