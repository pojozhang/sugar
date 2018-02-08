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
)

const (
	ContentType     = "Content-Type"
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJson = "application/json;charset=UTF-8"
)

var (
	Encode = ToString
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

type File string

type JSON struct {
	Data interface{}
}

type Cookie Map

type C = Cookie

type User struct {
	Name, Password string
}

type U = User

type MultiPart Map

type D = MultiPart

type Resolver interface {
	resolve(req *http.Request, params []interface{}, param interface{}, index int) error
}

type PathResolver struct {
}

func (r *PathResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
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
			value := param.(Path)[key]
			req.URL.Path = strings.Replace(req.URL.Path, req.URL.Path[i:j], Encode(value), -1)
		}
	}
	return nil
}

type QueryResolver struct {
}

func (r *QueryResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	q := req.URL.Query()
	for k, v := range param.(Query) {
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

func (r *HeaderResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	for k, v := range param.(Header) {
		req.Header.Add(k, Encode(v))
	}
	return nil
}

type FormResolver struct {
}

func (r *FormResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	form := url.Values{}
	for k, v := range param.(Form) {
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
}

func (r *JsonResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	v := param.(*JSON)
	var b []byte
	var err error
	switch x := v.Data.(type) {
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

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, ContentTypeJson)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	return nil
}

type CookieResolver struct {
}

func (r *CookieResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	c := param.(Cookie)
	for k, v := range c {
		req.AddCookie(&http.Cookie{Name: k, Value: Encode(v)})
	}
	return nil
}

type BasicAuthResolver struct {
}

func (r *BasicAuthResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	u := param.(User)
	req.SetBasicAuth(u.Name, u.Password)
	return nil
}

type MultiPartResolver struct {
}

func (r *MultiPartResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	m := param.(MultiPart)
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	defer w.Close()
	for k, v := range m {
		switch x := v.(type) {
		case *os.File:
			if err := writeFile(w, k, x); err != nil {
				return err
			}
		case File:
			if f, err := os.Open(string(x)); err != nil {
				return err
			} else {
				if err := writeFile(w, k, f); err != nil {
					f.Close()
					return err
				}
				f.Close()
			}
		default:
			w.WriteField(k, Encode(v))
		}
	}

	if _, ok := req.Header[ContentType]; !ok {
		req.Header.Set(ContentType, w.FormDataContentType())
	}
	req.Body = ioutil.NopCloser(b)
	return nil
}

func writeFile(w *multipart.Writer, fieldName string, file *os.File) error {
	f, err := file.Stat()
	if err != nil {
		return err
	}

	fileWriter, err := w.CreateFormFile(fieldName, f.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
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

func Json(v interface{}) *JSON {
	return J(v)
}

func J(v interface{}) *JSON {
	return &JSON{Data: v}
}

func foreach(v interface{}, f func(interface{})) {
	a := reflect.ValueOf(v)
	for i := 0; i < a.Len(); i++ {
		f(a.Index(i).Elem().Interface())
	}
}

func init() {
	Register(Path{}, &PathResolver{})
	Register(Query{}, &QueryResolver{})
	Register(Header{}, &HeaderResolver{})
	Register(JSON{}, &JsonResolver{})
	Register(Form{}, &FormResolver{})
	Register(Cookie{}, &CookieResolver{})
	Register(User{}, &BasicAuthResolver{})
	Register(MultiPart{}, &MultiPartResolver{})
}
