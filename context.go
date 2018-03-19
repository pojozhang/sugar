package sugar

import "net/http"

type Context struct {
	Request    *http.Request
	Response   *http.Response
	method     string
	rawUrl     string
	params     []interface{}
	plugins    []Plugin
	index      int
	encoders   []Encoder
	httpClient *http.Client
}

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

	if err := c.doRequest(); err != nil {
		return err
	}

	return nil
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
	if err != nil {
		return err
	}
	return nil
}
