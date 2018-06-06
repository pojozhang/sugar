package sugar

import "net/http"

// Context keeps all necessary params to build a request,
// and it allows us to pass params between plugins and encoders.
type Context struct {
	Request     *http.Request
	Response    *http.Response
	Method      string
	RawUrl      string
	params      []interface{}
	plugins     []Plugin
	index       int
	Encoders    EncoderGroup
	Decoders    DecoderGroup
	transporter Transporter
}

// BuildRequest initializes a new request and encodes params via encoders.
func (c *Context) BuildRequest() (*http.Request, error) {
	req, err := http.NewRequest(c.Method, c.RawUrl, nil)
	if err != nil {
		return nil, err
	}

	for i, param := range c.params {
		chain := NewEncoderChain(&RequestContext{Request: req, Params: c.params, Param: param, ParamIndex: i}, c.Encoders...)
		if err := chain.Next(); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (c *Context) reset() {
	c.index = 0
}

// Next invokes plugins and then sends the request via *http.Client.
func (c *Context) Next() error {
	if c.index < len(c.plugins) {
		c.index++
		return c.plugins[c.index-1].Handle(c)
	}

	resp, err := c.transporter.Do(c.Request)
	c.Response = resp
	return err
}
