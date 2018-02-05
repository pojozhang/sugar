package sugar

import (
	"net/http"
	"strconv"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strings"
)

type MapParams map[string]interface{}

type Header MapParams

type H = Header

type Path MapParams

type P = Path

type Query MapParams

type Q = Query

type Form MapParams

type F = Form

type Json MapParams

type J = Json

type Cookie MapParams

type C = Cookie

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

			key := req.URL.Path[i+1 : j]
			value := param.(Path)[key]
			req.URL.Path = strings.Replace(req.URL.Path, req.URL.Path[i:j], ToString(value), -1)
		}
	}
	return nil
}

type QueryResolver struct {
}

func (r *QueryResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	q := param.(Query)
	query := req.URL.Query()
	for k, v := range q {
		query.Add(k, ToString(v))
	}
	req.URL.RawQuery = query.Encode()
	return nil
}

type HeaderResolver struct {
}

func (r *HeaderResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	h := param.(Header)
	header := req.Header
	for k, v := range h {
		header.Add(k, ToString(v))
	}
	return nil
}

type FormResolver struct {
}

func (r *FormResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	return nil
}

type JsonResolver struct {
}

func (r *JsonResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	return nil
}

type CookieResolver struct {
}

func (r *CookieResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) error {
	c := param.(Cookie)
	for k, v := range c {
		req.AddCookie(&http.Cookie{Name: k, Value: ToString(v)})
	}
	return nil
}

func ToString(v interface{}) string {
	var s string
	switch x := v.(type) {
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
