package sugar

import (
	"net/http"
	"net/http/httputil"
	"os"
	"fmt"
)

type Client struct {
	HttpClient *http.Client
	Log        func(string)
	Resolvers  []Resolver
	Plugins    []Plugin
	Presets    []interface{}
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
	Use           = DefaultClient.Use
	DefaultLog    = func(s string) {
		os.Stdout.WriteString(fmt.Sprintf("%s\n", s))
	}
)

func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{},
		Log:        DefaultLog,
		Resolvers:  DefaultResolvers,
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
	context := &Context{
		method:        method,
		rawUrl:        rawUrl,
		params:        append(c.Presets, params...),
		plugins:       nil,
		resolverChain: NewResolverChain(),
		httpClient:    c.HttpClient,
	}
	if err := context.Next(); err != nil {
		return &Response{Error: err, request: context.Request}
	}

	if c.Log != nil {
		b, _ := httputil.DumpRequest(context.Request, true)
		c.Log(string(b))
	}

	return &Response{Response: *context.Response, Error: nil, request: context.Request}
}

func (c *Client) Apply(v ...interface{}) {
	c.Presets = append(c.Presets, v...)
}

func (c *Client) Reset(v ...interface{}) {
	c.Presets = nil
}

func (c *Client) Use(plugins ... Plugin) {
	c.Plugins = append(c.Plugins, plugins...)
}
