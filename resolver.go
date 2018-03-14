package sugar

import (
	"net/http"
	"strconv"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strings"
	"reflect"
	"net/url"
	"mime/multipart"
	"os"
	"io"
	"encoding/xml"
)

const (
	ContentType          = "Content-Type"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeJson      = "application/json; charset=UTF-8"
	ContentTypeXml       = "application/xml; charset=UTF-8"
	ContentTypePlainText = "text/plain"
)

var (
	Encode        = ToString
	resolverGroup ResolverGroup
)

type List []interface{}

type L = List

type Map map[string]interface{}

type M = Map

type Header Map

type H = Header

type Path Map

type P = Path

type Query Map

type Q = Query

type Form Map

type F = Form

type Json struct {
	Payload interface{}
}

type J = Json

type Cookie Map

type C = Cookie

type User struct {
	Name, Password string
}

type U = User

type MultiPart Map

type MP = MultiPart

type Xml struct {
	Payload interface{}
}

type Mapper struct {
	mapper func(*http.Request)
}

type RequestContext struct {
	Request *http.Request
	Params  []interface{}
	Param   interface{}
	Index   int
}

type Resolver interface {
	Resolve(context *RequestContext, chain *ResolverChain) error
}

type ResolverChain struct {
	resolver  Resolver
	successor *ResolverChain
}

func (r *ResolverChain) Next(context *RequestContext) error {
	if r != nil {
		return r.resolver.Resolve(context, r.successor)
	}
	return ResolverNotFound
}

type PathResolver struct {
}

func (r *PathResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	pathParams, ok := context.Param.(Path)
	if !ok {
		return chain.Next(context)
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

			key := req.URL.Path[i+1: j]
			value := pathParams[key]
			req.URL.Path = strings.Replace(req.URL.Path, req.URL.Path[i:j], Encode(value), -1)
		}
	}
	return nil
}

type QueryResolver struct {
}

func (r *QueryResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	queryParams, ok := context.Param.(Query)
	if !ok {
		return chain.Next(context)
	}

	req := context.Request
	q := req.URL.Query()
	for k, v := range queryParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			foreach(v, func(i interface{}) {
				q.Add(k, Encode(i))
			})
		default:
			q.Add(k, Encode(v))
		}
	}
	req.URL.RawQuery = q.Encode()
	return nil
}

type HeaderResolver struct {
}

func (r *HeaderResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	headerParams, ok := context.Param.(Header)
	if !ok {
		return chain.Next(context)
	}

	for k, v := range headerParams {
		context.Request.Header.Add(k, Encode(v))
	}
	return nil
}

type FormResolver struct {
}

func (r *FormResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	formParams, ok := context.Param.(Form)
	if !ok {
		return chain.Next(context)
	}

	form := url.Values{}
	for k, v := range formParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			foreach(v, func(i interface{}) {
				form.Add(k, Encode(i))
			})
		default:
			form.Add(k, Encode(v))
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

type JsonResolver struct {
}

func (r *JsonResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	jsonParams, ok := context.Param.(Json)
	if !ok {
		return chain.Next(context)
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
		req.Header.Set(ContentType, ContentTypeJson)
	}
	return nil
}

type CookieResolver struct {
}

func (r *CookieResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	cookieParams, ok := context.Param.(Cookie)
	if !ok {
		return chain.Next(context)
	}

	for k, v := range cookieParams {
		context.Request.AddCookie(&http.Cookie{Name: k, Value: Encode(v)})
	}
	return nil
}

type BasicAuthResolver struct {
}

func (r *BasicAuthResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	authParams, ok := context.Param.(User)
	if !ok {
		return chain.Next(context)
	}

	context.Request.SetBasicAuth(authParams.Name, authParams.Password)
	return nil
}

type MultiPartResolver struct {
}

func (r *MultiPartResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	multiPartParams, ok := context.Param.(MultiPart)
	if !ok {
		return chain.Next(context)
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
			w.WriteField(k, Encode(v))
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

type PlainTextResolver struct {
}

func (r *PlainTextResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	textParams, ok := context.Param.(string)
	if !ok {
		return chain.Next(context)
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

type XmlResolver struct {
}

func (r *XmlResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	xmlParams, ok := context.Param.(Xml)
	if !ok {
		return chain.Next(context)
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
		req.Header.Set(ContentType, ContentTypeXml)
	}
	return nil
}

type MapperResolver struct {
}

func (r *MapperResolver) Resolve(context *RequestContext, chain *ResolverChain) error {
	mapperParams, ok := context.Param.(Mapper)
	if !ok {
		return chain.Next(context)
	}

	mapperParams.mapper(context.Request)
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

type ResolverGroup struct {
	chain []*ResolverChain
}

func (g *ResolverGroup) Add(resolvers ... Resolver) {
	for _, resolver := range resolvers {
		c := &ResolverChain{resolver: resolver}
		g.chain = append(g.chain, c)
	}
	g.refresh()
}

func (g *ResolverGroup) refresh() {
	for i := 0; i < len(g.chain)-1; i++ {
		g.chain[i].successor = g.chain[i+1]
	}
}

func (g *ResolverGroup) Resolve(context *RequestContext) error {
	if len(g.chain) < 1 {
		return nil
	}
	return g.chain[0].Next(context)
}

func RegisterResolvers(resolvers ... Resolver) {
	resolverGroup.Add(resolvers...)
}

func init() {
	RegisterResolvers(
		&XmlResolver{},
		&PathResolver{},
		&JsonResolver{},
		&FormResolver{},
		&QueryResolver{},
		&HeaderResolver{},
		&MapperResolver{},
		&CookieResolver{},
		&BasicAuthResolver{},
		&MultiPartResolver{},
		&PlainTextResolver{},
	)
}
