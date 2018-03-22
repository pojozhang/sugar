package sugar

import "net/http"

// Context is the core of Sugar.
// It keeps all necessary params to build a request,
// and it allows us to pass params between plugins, encoders and decoders.
type Context struct {
	Request    *http.Request
	Response   *http.Response
	method     string
	rawUrl     string
	params     []interface{}
	plugins    []Plugin
	index      int
	encoders   []Encoder
	decoders   []Decoder
	httpClient *http.Client
}

// Execute is the entrance of Context.
func (c *Context) Execute() *Response {
	c.reset()

	if err := c.Next(); err != nil {
		return &Response{Error: err, request: c.Request, decoders: c.decoders}
	}

	return &Response{Response: *c.Response, Error: nil, request: c.Request, decoders: c.decoders}
}

func (c *Context) reset() {
	c.index = 0
}

// Next is the core method of Context.
// It encodes request params via encoders, invokes plugins
// and then send the request via *http.Client.
func (c *Context) Next() error {
	if c.index == 0 {
		if err := c.prepareRequest(); err != nil {
			return err
		}

		if err := c.encodeRequest(); err != nil {
			return err
		}
	}

	if c.index < len(c.plugins) {
		c.index++
		return c.plugins[c.index-1](c)
	}

	return c.doRequest()
}

func (c *Context) prepareRequest() error {
	req, err := http.NewRequest(c.method, c.rawUrl, nil)
	if err != nil {
		return err
	}
	c.Request = req
	return nil
}

func (c *Context) encodeRequest() error {
	for i, param := range c.params {
		chain := NewEncoderChain(&RequestContext{Request: c.Request, Params: c.params, Param: param, ParamIndex: i}, c.encoders...)
		if err := chain.Next(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) doRequest() error {
	resp, err := c.httpClient.Do(c.Request)
	c.Response = resp
	return err
}
