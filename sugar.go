package sugar

import (
	"net/http"
	"net/http/httputil"
	"reflect"
	"errors"
	"os"
	"fmt"
)

var (
	resolvers = make(map[reflect.Type]Resolver)
)

type Client struct {
	HttpClient *http.Client
	Log        func(string)
	presets    []interface{}
}

var (
	DefaultClient = NewClient()
	DefaultLog    = func(s string) {
		os.Stdout.WriteString(fmt.Sprintf("%s\n", s))
	}
)

func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{},
		Log:        DefaultLog,
	}
}

func (c *Client) Get(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodGet, rawUrl, params...)
}

func (c *Client) Post(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPost, rawUrl, params...)
}

func (c *Client) Put(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPut, rawUrl, params...)
}

func (c *Client) Patch(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPatch, rawUrl, params...)
}

func (c *Client) Delete(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodDelete, rawUrl, params...)
}

func (c *Client) Do(method, rawUrl string, params ...interface{}) (*Response, error) {
	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return nil, err
	}

	params = append(c.presets, params...)
	for i, param := range params {
		t := reflect.TypeOf(param)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if r, ok := resolvers[t]; ok {
			if err := r.resolve(req, params, param, i); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("No resolvers found for " + t.String())
		}
	}

	if c.Log != nil {
		b, _ := httputil.DumpRequest(req, true)
		c.Log(string(b))
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	r := Response(*resp)
	return &r, err
}

func (c *Client) Apply(v ...interface{}) {
	c.presets = append(c.presets, v...)
}

func (c *Client) Reset(v ...interface{}) {
	c.presets = nil
}

func Get(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Get(rawUrl, params...)
}

func Post(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Post(rawUrl, params...)
}

func Put(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Put(rawUrl, params...)
}

func Patch(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Patch(rawUrl, params...)
}

func Delete(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Delete(rawUrl, params...)
}

func Do(method, rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Do(method, rawUrl, params...)
}

func Apply(v ...interface{}) {
	DefaultClient.Apply(v...)
}

func Reset() {
	DefaultClient.Reset()
}

func Resolvers() map[reflect.Type]Resolver {
	return resolvers
}

func Register(v interface{}, resolver Resolver) {
	resolvers[reflect.TypeOf(v)] = resolver
}
