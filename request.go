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
	Encoders []Encoder
	Decoders []Decoder

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
)

// NewClient returns a new Client given a http client, encoders and decoders.
func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{},
		Encoders:   Encoders,
		Decoders:   Decoders,
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

func (c *Client) Do(method, rawUrl string, params ...interface{}) *Response {
	context := &Context{
		method:     method,
		rawUrl:     rawUrl,
		params:     append(c.Presets, params...),
		encoders:   c.Encoders,
		plugins:    c.Plugins,
		httpClient: c.HttpClient,
	}
	if err := context.Next(); err != nil {
		return &Response{Error: err, request: context.Request, decoders: c.Decoders}
	}

	return &Response{Response: *context.Response, Error: nil, request: context.Request, decoders: c.Decoders}
}

func (c *Client) Apply(v ...interface{}) {
	c.Presets = append(c.Presets, v...)
}

func (c *Client) Reset(v ...interface{}) {
	c.Presets = nil
}

func (c *Client) Use(plugins ...Plugin) {
	c.Plugins = append(c.Plugins, plugins...)
}

// RegisterEncoders registers global encoders.
func RegisterEncoders(encoders ...Encoder) {
	Encoders = append(Encoders, encoders...)
}

// RegisterDecoders registers global decoders.
func RegisterDecoders(decoders ...Decoder) {
	Decoders = append(Decoders, decoders...)
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
	)

	DefaultClient.Encoders = Encoders
	DefaultClient.Decoders = Decoders
}
