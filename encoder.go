package sugar

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	Stringify = ToString
)

type List []interface{}

// L is an alias for List.
type L = List

type Map map[string]interface{}

// M is an alias for Map.
type M = Map

type Header Map

// H is an alias for Header.
type H = Header

type Cookie Map

// C is an alias for Cookie.
type C = Cookie

type Path Map

// P is an alias for Path.
type P = Path

type Query Map

// Q is an alias for Query.
type Q = Query

type Form Map

// F is an alias for Form.
type F = Form

type Json struct {
	Payload interface{}
}

// J is an alias for Json.
type J = Json

type Xml struct {
	Payload interface{}
}

// X is an alias for Xml.
type X = Xml

type User struct {
	Name, Password string
}

// U is an alias for User.
type U = User

type MultiPart Map

// MP is an alias for MultiPart.
type MP = MultiPart

// RequestContext keeps values for an encoder.
type RequestContext struct {
	Request    *http.Request
	Response   *http.Response
	Params     []interface{}
	Param      interface{}
	ParamIndex int
}

// Encoder converts a request context into request params.
// It returns an error if any error occurs during encoding.
// Call chain.Next() to propagate context.
type Encoder interface {
	Encode(context *RequestContext, chain *EncoderChain) error
}

// EncoderChain keeps a set of encoders.
type EncoderChain struct {
	context  *RequestContext
	encoders []Encoder
	index    int
}

// Next propagates context to next encoder.
// It returns EncoderNotFound if current encoder is the last one.
func (c *EncoderChain) Next() error {
	if c.index < len(c.encoders) {
		c.index++
		return c.encoders[c.index-1].Encode(c.context, c)
	}
	return EncoderNotFound
}

func (c *EncoderChain) reset() *EncoderChain {
	c.index = 0
	return c
}

// Add adds encoders to a encoder chain.
func (c *EncoderChain) Add(Encoders ...Encoder) *EncoderChain {
	for _, Encoder := range Encoders {
		c.encoders = append(c.encoders, Encoder)
	}
	return c
}

// NewEncoderChain initializes a new encoder chain with given request context and encoders.
func NewEncoderChain(context *RequestContext, encoders ...Encoder) *EncoderChain {
	chain := &EncoderChain{context: context, index: 0}
	chain.reset().Add(encoders...)
	return chain
}

// PathEncoder encodes Path{} params.
type PathEncoder struct {
}

// Encode encodes Path{} params.
func (r *PathEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	pathParams, ok := context.Param.(Path)
	if !ok {
		return chain.Next()
	}

	req := context.Request
	for i := 0; i < len(req.URL.Path); i++ {
		if string(req.URL.Path[i]) == ":" {
			j := i + 1
			for ; j < len(req.URL.Path); j++ {
				s := string(req.URL.Path[j])
				if s == "/" {
					break
				}
			}

			key := req.URL.Path[i+1 : j]
			value := pathParams[key]
			req.URL.Path = strings.Replace(req.URL.Path, req.URL.Path[i:j], Stringify(value), -1)
		}
	}
	return nil
}

// QueryEncoder encodes Query{} params.
type QueryEncoder struct {
}

// Encode encodes Query{} params.
func (r *QueryEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	queryParams, ok := context.Param.(Query)
	if !ok {
		return chain.Next()
	}

	req := context.Request
	q := req.URL.Query()
	for k, v := range queryParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			foreach(v, func(i interface{}) {
				q.Add(k, Stringify(i))
			})
		default:
			q.Add(k, Stringify(v))
		}
	}
	req.URL.RawQuery = strings.Replace(q.Encode(), "+", "%20", -1)
	return nil
}

// HeaderEncoder encodes Header{} params.
type HeaderEncoder struct {
}

// Encode encodes Header{} params.
func (r *HeaderEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	headerParams, ok := context.Param.(Header)
	if !ok {
		return chain.Next()
	}

	for k, v := range headerParams {
		context.Request.Header.Add(k, Stringify(v))
	}
	return nil
}

// FormEncoder encodes Form{} params.
type FormEncoder struct {
}

// Encode encodes Form{} params.
func (r *FormEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	formParams, ok := context.Param.(Form)
	if !ok {
		return chain.Next()
	}

	form := url.Values{}
	for k, v := range formParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			foreach(v, func(i interface{}) {
				form.Add(k, Stringify(i))
			})
		default:
			form.Add(k, Stringify(v))
		}
	}

	req := context.Request
	req.PostForm = form
	err := req.ParseForm()
	if err != nil {
		return err
	}

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeForm)
	}

	return nil
}

// JsonEncoder encodes Json{} params.
type JsonEncoder struct {
}

// Encode encodes Json{} params.
func (r *JsonEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	jsonParams, ok := context.Param.(Json)
	if !ok {
		return chain.Next()
	}

	var b []byte
	var err error
	switch x := jsonParams.Payload.(type) {
	case []byte:
		b, err = json.RawMessage(x).MarshalJSON()
	case string:
		b, err = json.RawMessage([]byte(x)).MarshalJSON()
	default:
		b, err = json.Marshal(x)
	}

	if err != nil {
		return err
	}

	req := context.Request
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeJsonUtf8)
	}
	return nil
}

// CookieEncoder encodes Cookie{} params.
type CookieEncoder struct {
}

// Encode encodes Cookie{} params.
func (r *CookieEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	cookieParams, ok := context.Param.(Cookie)
	if !ok {
		return chain.Next()
	}

	for k, v := range cookieParams {
		context.Request.AddCookie(&http.Cookie{Name: k, Value: Stringify(v)})
	}
	return nil
}

// BasicAuthEncoder encodes User{} params.
type BasicAuthEncoder struct {
}

// Encode encodes User{} params.
func (r *BasicAuthEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	authParams, ok := context.Param.(User)
	if !ok {
		return chain.Next()
	}

	context.Request.SetBasicAuth(authParams.Name, authParams.Password)
	return nil
}

// MultiPartEncoder encodes MultiPart{} params.
type MultiPartEncoder struct {
}

// Encode encodes MultiPart{} params.
func (r *MultiPartEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	multiPartParams, ok := context.Param.(MultiPart)
	if !ok {
		return chain.Next()
	}

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	defer w.Close()
	for k, v := range multiPartParams {
		switch x := v.(type) {
		case *os.File:
			if err := writeFile(w, k, x.Name(), x); err != nil {
				return err
			}
		default:
			w.WriteField(k, Stringify(v))
		}
	}

	req := context.Request
	req.Body = ioutil.NopCloser(b)

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, w.FormDataContentType())
	}
	return nil
}

func writeFile(w *multipart.Writer, fieldName, fileName string, file io.Reader) error {
	fileWriter, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
	}

	return nil
}

// PlainTextEncoder encodes string params.
type PlainTextEncoder struct {
}

// Encode encodes string params.
func (r *PlainTextEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	textParams, ok := context.Param.(string)
	if !ok {
		return chain.Next()
	}

	b := &bytes.Buffer{}
	b.WriteString(textParams)
	req := context.Request
	req.Body = ioutil.NopCloser(b)

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypePlainText)
	}
	return nil
}

// XmlEncoder encodes Xml{} params.
type XmlEncoder struct {
}

// Encode encodes Xml{} params.
func (r *XmlEncoder) Encode(context *RequestContext, chain *EncoderChain) error {
	xmlParams, ok := context.Param.(Xml)
	if !ok {
		return chain.Next()
	}

	var b []byte
	var err error
	switch x := xmlParams.Payload.(type) {
	case string:
		b = []byte(x)
	default:
		b, err = xml.Marshal(x)
	}

	if err != nil {
		return err
	}

	req := context.Request
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeXmlUtf8)
	}
	return nil
}

func ToString(v interface{}) string {
	var s string
	switch x := v.(type) {
	case bool:
		s = strconv.FormatBool(x)
	case uint:
		s = strconv.FormatUint(uint64(x), 10)
	case uint8:
		s = strconv.FormatUint(uint64(x), 10)
	case uint16:
		s = strconv.FormatUint(uint64(x), 10)
	case uint32:
		s = strconv.FormatUint(uint64(x), 10)
	case uint64:
		s = strconv.FormatUint(uint64(x), 10)
	case int:
		s = strconv.FormatInt(int64(x), 10)
	case int8:
		s = strconv.FormatInt(int64(x), 10)
	case int16:
		s = strconv.FormatInt(int64(x), 10)
	case int32:
		s = strconv.FormatInt(int64(x), 10)
	case int64:
		s = strconv.FormatInt(int64(x), 10)
	case float32:
		s = strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(float64(x), 'f', -1, 64)
	case string:
		s = v.(string)
	}
	return s
}

func foreach(v interface{}, f func(interface{})) {
	a := reflect.ValueOf(v)
	for i := 0; i < a.Len(); i++ {
		f(a.Index(i).Elem().Interface())
	}
}
