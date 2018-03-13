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
	Get           = DefaultClient.Get
	Post          = DefaultClient.Post
	Put           = DefaultClient.Put
	Patch         = DefaultClient.Patch
	Delete        = DefaultClient.Delete
	Do            = DefaultClient.Do
	Apply         = DefaultClient.Apply
	Reset         = DefaultClient.Reset
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

func (c *Client) Get(rawUrl string, params ...interface{}) (*Response) {
	return c.Do(http.MethodGet, rawUrl, params...)
}

func (c *Client) Post(rawUrl string, params ...interface{}) (*Response) {
	return c.Do(http.MethodPost, rawUrl, params...)
}

func (c *Client) Put(rawUrl string, params ...interface{}) (*Response) {
	return c.Do(http.MethodPut, rawUrl, params...)
}

func (c *Client) Patch(rawUrl string, params ...interface{}) (*Response) {
	return c.Do(http.MethodPatch, rawUrl, params...)
}

func (c *Client) Delete(rawUrl string, params ...interface{}) (*Response) {
	return c.Do(http.MethodDelete, rawUrl, params...)
}

func (c *Client) Do(method, rawUrl string, params ...interface{}) (*Response) {
	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return &Response{Error: err}
	}

	params = append(c.presets, params...)
	for i, param := range params {
		t := reflect.TypeOf(param)
		if r, ok := resolvers[t]; ok {
			if err := r.Resolve(req, params, param, i); err != nil {
				return &Response{Error: err}
			}
		} else {
			return &Response{Error: errors.New("No resolvers found for " + t.String())}
		}
	}

	if c.Log != nil {
		b, _ := httputil.DumpRequest(req, true)
		c.Log(string(b))
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{Response: *resp, Error: nil}
}

func (c *Client) Apply(v ...interface{}) {
	c.presets = append(c.presets, v...)
}

func (c *Client) Reset(v ...interface{}) {
	c.presets = nil
}

func Resolvers() map[reflect.Type]Resolver {
	return resolvers
}

func RegisterResolver(v interface{}, resolver Resolver) {
	resolvers[reflect.TypeOf(v)] = resolver
}
