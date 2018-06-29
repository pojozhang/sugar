<img align="middle" height="200px" src="logo.png">

![GitHub (pre-)release](https://img.shields.io/github/release/pojozhang/sugar/all.svg)
[![Build Status](https://travis-ci.org/pojozhang/sugar.svg?branch=master)](https://travis-ci.org/pojozhang/sugar) [![codecov](https://codecov.io/gh/pojozhang/sugar/branch/master/graph/badge.svg)](https://codecov.io/gh/pojozhang/sugar) [![Go Report Card](https://goreportcard.com/badge/github.com/pojozhang/sugar)](https://goreportcard.com/report/github.com/pojozhang/sugar) ![go](https://img.shields.io/badge/golang-1.9+-blue.svg) [![GoDoc](https://godoc.org/github.com/pojozhang/sugar?status.svg)](https://godoc.org/github.com/pojozhang/sugar) ![license](https://img.shields.io/github/license/pojozhang/sugar.svg)

Sugar是一个Go语言编写的声明式Http客户端，提供了一些优雅的接口，目的是减少冗余的拼装代码。

## 特性
- 优雅的接口设计
- 插件功能
- 链式调用
- 高度可定制

## 下载
```bash
dep ensure -add github.com/pojozhang/sugar
```

## 使用
首先导入包，为了看起来更简洁，此处用省略包名的方式导入。
```go
import . "github.com/pojozhang/sugar"
```
一切就绪！

### 发送请求
#### Plain Text
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: text/plain
Post("http://api.example.com/books", "bookA")
```

#### Path
```go
// GET /books/123 HTTP/1.1
// Host: api.example.com
Get("http://api.example.com/books/:id", Path{"id": 123})
Get("http://api.example.com/books/:id", P{"id": 123})
```

#### Query
```go
// GET /books?name=bookA HTTP/1.1
// Host: api.example.com
Get("http://api.example.com/books", Query{"name": "bookA"})
Get("http://api.example.com/books", Q{"name": "bookA"})

// list
// GET /books?name=bookA&name=bookB HTTP/1.1
// Host: api.example.com
Get("http://api.example.com/books", Query{"name": List{"bookA", "bookB"}})
Get("http://api.example.com/books", Q{"name": L{"bookA", "bookB"}})
```

#### Cookie
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Cookie: name=bookA
Get("http://api.example.com/books", Cookie{"name": "bookA"})
Get("http://api.example.com/books", C{"name": "bookA"})
```

#### Header
```go
// GET /books HTTP/1.1
// Host: api.example.com
// Name: bookA
Get("http://api.example.com/books", Header{"name": "bookA"})
Get("http://api.example.com/books", H{"name": "bookA"})
```

#### Json
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/json;charset=UTF-8
// {"name":"bookA"}
Post("http://api.example.com/books", Json{`{"name":"bookA"}`})
Post("http://api.example.com/books", J{`{"name":"bookA"}`})

// map
Post("http://api.example.com/books", Json{Map{"name": "bookA"}})
Post("http://api.example.com/books", J{M{"name": "bookA"}})

// list
Post("http://api.example.com/books", Json{List{Map{"name": "bookA"}}})
Post("http://api.example.com/books", J{L{M{"name": "bookA"}}})
```

#### Xml
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Authorization: Basic dXNlcjpwYXNzd29yZA==
// Content-Type: application/xml; charset=UTF-8
// <book name="bookA"></book>
Post("http://api.example.com/books", Xml{`<book name="bookA"></book>`})
Post("http://api.example.com/books", X{`<book name="bookA"></book>`})
```

#### Form
```go
// POST /books HTTP/1.1
// Host: api.example.com
// Content-Type: application/x-www-form-urlencoded
Post("http://api.example.com/books", Form{"name": "bookA"})
Post("http://api.example.com/books", F{"name": "bookA"})

// list
Post("http://api.example.com/books", Form{"name": List{"bookA", "bookB"}})
Post("http://api.example.com/books", F{"name": L{"bookA", "bookB"}})
```

#### Basic Auth
```go
// DELETE /books HTTP/1.1
// Host: api.example.com
// Authorization: Basic dXNlcjpwYXNzd29yZA==
Delete("http://api.example.com/books", User{"user", "password"})
Delete("http://api.example.com/books", U{"user", "password"})
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
Post("http://api.example.com/books", MultiPart{"name": "bookA", "file": f})
Post("http://api.example.com/books", MP{"name": "bookA", "file": f})
```

#### Mix
你可以任意组合参数。
```go
Patch("http://api.example.com/books/:id", Path{"id": 123}, Json{`{"name":"bookA"}`}, User{"user", "password"})
```

#### Apply
Apply方法传入的参数会被应用到之后所有的请求中，可以使用Reset()方法重置。
```go
Apply(User{"user", "password"})
Get("http://api.example.com/books")
Get("http://api.example.com/books")
Reset()
Get("http://api.example.com/books")
```
```go
Get("http://api.example.com/books", User{"user", "password"})
Get("http://api.example.com/books", User{"user", "password"})
Get("http://api.example.com/books")
```
以上两段代码是等价的。


### 解析响应
一个请求发送后会返回`*Response`类型的返回值，其中包含了一些有用的语法糖。

#### Raw
Raw()会返回一个`*http.Response`和一个`error`，就和Go自带的SDK一样（所以叫Raw）。
```go
resp, err := Post("http://api.example.com/books", "bookA").Raw()
...
```

#### ReadBytes
ReadBytes()可以直接从返回的`body`读取字节切片。需要注意的是，该方法返回前会自动释放`body`资源。
```go
bytes, resp, err := Get("http://api.example.com/books").ReadBytes()
...
```

#### Read
Read()方法通过注册在系统中的`Decoder`对返回值进行解析。
以下两个例子是在不同的情况下分别解析成字符串或者JSON，解析过程对调用者来说是透明的。
```go
// plain text
var text = new(string)
resp, err := Get("http://api.example.com/text").Read(text)

// json
var books []book
resp, err := Get("http://api.example.com/json").Read(&books)
```

#### 文件下载
我们也可以通过`Read()`方法下载文件。
```go
f,_ := os.Create("tmp.png")
defer f.Close()
resp, err := Get("http://api.example.com/logo.png").Read(f)
```

## 自定义
Sugar中有三大组件 **Encoder**, **Decoder** 和 **Plugin**.
- **Encoder**负责把调用者传入参数组装成一个请求体。
- **Decoder**负责把服务器返回的数据解析成一个结构体。
- **Plugin**起到拦截器的作用。

### Encoder
你可以通过实现`Encoder`接口来实现自己的编码器。
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

Get("http://api.example.com/books", MyParam{})
```

### Decoder
你可以实现`Decoder`接口来实现自己的解码器。`Read()`方法会使用解码器去解析返回值。
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
插件是一个特殊的组件，你可以在请求发送前或收到响应后进行一些额外的处理。
```go
// 内置Logger插件的实现
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
Logger插件用来记录发送出去的请求数据以及接收到的响应数据。
```go
Use(Logger)
```

#### Retryer
Retryer插件用来在请求遇到错误时自动进行重试。
```go
Use(Retryer(3, time.Second, 1, time.Second))
```

#### 自定义接口异常处理
通过插件机制，我们可以定制一个异常处理器来处理接口返回的错误描述。下面这个例子展示了当服务器返回错误码和错误信息时如何用插件进行处理：
```go
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

// 发送请求
_, err := client.Get("some url").Read(&json{})
// 类型判断
switch e := err.(type) {
	case apiError:
		// ...
	}
// ...
```
