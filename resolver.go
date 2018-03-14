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

type D = MultiPart

type Xml struct {
	Payload interface{}
}

type Mapper struct {
	mapper func(*http.Request)
}

type Resolver interface {
	Resolve(req *http.Request, params []interface{}, param interface{}, index int) error
}

type ChainedResolver struct {
	Resolver
	successor Resolver
}

func (r *ChainedResolver) Next(req *http.Request, params []interface{}, param interface{}, index int) error {
	if r.successor != nil {
		return r.successor.Resolve(req, params, param, index)
	}
	return nil
}

type PathResolver struct {
	ChainedResolver
}

func (r *PathResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	pathParams, ok := param.(Path)
	if !ok {
		return r.Next(req, params, param, index)
	}

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
			req.URL.Path = strings.Replace(req.URL.Path, req.URL.Path[i:j], Encode(value), -1)
		}
	}
	return nil
}

type QueryResolver struct {
	ChainedResolver
}

func (r *QueryResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	queryParams, ok := param.(Query)
	if !ok {
		return r.Next(req, params, param, index)
	}

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
	ChainedResolver
}

func (r *HeaderResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	headerParams, ok := param.(Header)
	if !ok {
		return r.Next(req, params, param, index)
	}

	for k, v := range headerParams {
		req.Header.Add(k, Encode(v))
	}
	return nil
}

type FormResolver struct {
	ChainedResolver
}

func (r *FormResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	formParams, ok := param.(Form)
	if !ok {
		return r.Next(req, params, param, index)
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
	ChainedResolver
}

func (r *JsonResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	jsonParams, ok := param.(Json)
	if !ok {
		return r.Next(req, params, param, index)
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
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeJson)
	}
	return nil
}

type CookieResolver struct {
	ChainedResolver
}

func (r *CookieResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	cookieParams, ok := param.(Cookie)
	if !ok {
		return r.Next(req, params, param, index)
	}

	for k, v := range cookieParams {
		req.AddCookie(&http.Cookie{Name: k, Value: Encode(v)})
	}
	return nil
}

type BasicAuthResolver struct {
	ChainedResolver
}

func (r *BasicAuthResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	authParams, ok := param.(User)
	if !ok {
		return r.Next(req, params, param, index)
	}

	req.SetBasicAuth(authParams.Name, authParams.Password)
	return nil
}

type MultiPartResolver struct {
	ChainedResolver
}

func (r *MultiPartResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	multiPartParams, ok := param.(MultiPart)
	if !ok {
		return r.Next(req, params, param, index)
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
	ChainedResolver
}

func (r *PlainTextResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	textParams, ok := param.(string)
	if !ok {
		return r.Next(req, params, param, index)
	}

	b := &bytes.Buffer{}
	b.WriteString(textParams)
	req.Body = ioutil.NopCloser(b)

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypePlainText)
	}
	return nil
}

type XmlResolver struct {
	ChainedResolver
}

func (r *XmlResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	xmlParams, ok := param.(Xml)
	if !ok {
		return r.Next(req, params, param, index)
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
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeXml)
	}
	return nil
}

type MapperResolver struct {
	ChainedResolver
}

func (r *MapperResolver) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	mapperParams, ok := param.(Mapper)
	if !ok {
		return r.Next(req, params, param, index)
	}

	mapperParams.mapper(req)
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
	resolvers []*ChainedResolver
}

func (g *ResolverGroup) Add(resolvers ... Resolver) {
	for _, resolver := range resolvers {
		c := &ChainedResolver{Resolver: resolver}

		g.resolvers = append(g.resolvers, c)
	}
	g.refresh()
}

func (g *ResolverGroup) refresh() {
	for i := 0; i < len(g.resolvers)-1; i++ {
		g.resolvers[i].successor = g.resolvers[i+1]
	}
}

func (g *ResolverGroup) Resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	if len(g.resolvers) < 1 {
		return nil
	}
	return g.resolvers[0].Resolve(req, params, param, index)
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
