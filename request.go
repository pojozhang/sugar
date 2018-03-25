package sugar

import (
	"net/http"
)

// Client is a entrance to Sugar.
// It keeps important components for building requests and parsing responses.
type Client struct {
	HttpClient *http.Client
	Encoders   []Encoder
	Decoders   []Decoder
	Plugins    []Plugin
	Presets    []interface{}
}

var (
	defaultClient = &Client{
		HttpClient: &http.Client{},
	}
	Get        = defaultClient.Get
	Post       = defaultClient.Post
	Put        = defaultClient.Put
	Patch      = defaultClient.Patch
	Delete     = defaultClient.Delete
	Do         = defaultClient.Do
	Apply      = defaultClient.Apply
	Reset      = defaultClient.Reset
	Use        = defaultClient.Use
	NewRequest = defaultClient.NewRequest
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

// RegisterEncoders registers global encoders.
func RegisterEncoders(encoders ...Encoder) {
	defaultClient.Encoders = append(defaultClient.Encoders, encoders...)
}

// RegisterDecoders registers global decoders.
func RegisterDecoders(decoders ...Decoder) {
	defaultClient.Decoders = append(defaultClient.Decoders, decoders...)
}

func init() {
	RegisterEncoders(
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

	RegisterDecoders(
		&JsonDecoder{},
		&XmlDecoder{},
		&PlainTextDecoder{},
		&FileDecoder{},
	)
}
