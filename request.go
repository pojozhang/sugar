package sugar

import (
	"net/http"
)

type Encoders []Encoder

func (e *Encoders) Add(encoders ...Encoder) {
	*e = append(*e, encoders...)
}

type Decoders []Decoder

func (d *Decoders) Add(decoders ...Decoder) {
	*d = append(*d, decoders...)
}

// Client is a entrance to Sugar.
// It keeps important components for building requests and parsing responses.
type Client struct {
	HttpClient *http.Client
	Encoders   Encoders
	Decoders   Decoders
	Plugins    []Plugin
	Presets    []interface{}
}

var (
	defaultClient = &Client{
		HttpClient: &http.Client{},
	}
	Default    = defaultClient
	Get        = Default.Get
	Post       = Default.Post
	Put        = Default.Put
	Patch      = Default.Patch
	Delete     = Default.Delete
	Do         = Default.Do
	Apply      = Default.Apply
	Reset      = Default.Reset
	Use        = Default.Use
	NewRequest = Default.NewRequest
)

// NewClient returns a new Client given a http client, encoders and decoders.
func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{},
		Encoders:   defaultClient.Encoders,
		Decoders:   defaultClient.Decoders,
	}
}

// Get is a shortcut for client.Do("Get", url, params).
func (c *Client) Get(rawUrl string, params ...interface{}) *Response {
	return c.Do(http.MethodGet, rawUrl, params...)
}

// Post is a shortcut for client.Do("Post", url, params).
func (c *Client) Post(rawUrl string, params ...interface{}) *Response {
	return c.Do(http.MethodPost, rawUrl, params...)
}

// Put is a shortcut for client.Do("Put", url, params).
func (c *Client) Put(rawUrl string, params ...interface{}) *Response {
	return c.Do(http.MethodPut, rawUrl, params...)
}

// Patch is a shortcut for client.Do("Patch", url, params).
func (c *Client) Patch(rawUrl string, params ...interface{}) *Response {
	return c.Do(http.MethodPatch, rawUrl, params...)
}

// Delete is a shortcut for client.Do("Delete", url, params).
func (c *Client) Delete(rawUrl string, params ...interface{}) *Response {
	return c.Do(http.MethodDelete, rawUrl, params...)
}

// Do builds a context and then sends a request via the context.
func (c *Client) Do(method, rawUrl string, params ...interface{}) *Response {
	context := &Context{
		method:     method,
		rawUrl:     rawUrl,
		params:     append(c.Presets, params...),
		encoders:   c.Encoders,
		plugins:    c.Plugins,
		httpClient: c.HttpClient,
	}

	req, err := context.BuildRequest()
	if err != nil {
		return &Response{Error: err, request: req, decoders: c.Decoders}
	}

	context.Request = req
	context.reset()
	if err := context.Next(); err != nil {
		return &Response{Error: err, request: req, decoders: c.Decoders}
	}

	return &Response{Response: *context.Response, Error: nil, request: context.Request, decoders: c.Decoders}
}

// NewRequest builds a request via context.
func (c *Client) NewRequest(method, rawUrl string, params ...interface{}) (*http.Request, error) {
	context := &Context{
		method:   method,
		rawUrl:   rawUrl,
		params:   append(c.Presets, params...),
		encoders: c.Encoders,
	}

	return context.BuildRequest()
}

// Apply attaches params to every following request.
// Call Reset() to clean.
func (c *Client) Apply(v ...interface{}) {
	c.Presets = append(c.Presets, v...)
}

// Reset cleans all params added by Apply().
func (c *Client) Reset(v ...interface{}) {
	c.Presets = nil
}

// Use applies plugins.
func (c *Client) Use(plugins ...Plugin) {
	c.Plugins = append(c.Plugins, plugins...)
}

func init() {
	Default.Encoders.Add(
		&XmlEncoder{},
		&PathEncoder{},
		&JsonEncoder{},
		&FormEncoder{},
		&QueryEncoder{},
		&HeaderEncoder{},
		&CookieEncoder{},
		&BasicAuthEncoder{},
		&MultiPartEncoder{},
		&PlainTextEncoder{},
	)

	Default.Decoders.Add(
		&JsonDecoder{},
		&XmlDecoder{},
		&PlainTextDecoder{},
		&FileDecoder{},
	)
}
