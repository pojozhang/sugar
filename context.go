package sugar

import "net/http"

type Context struct {
	Request       *http.Request
	Response      *http.Response
	method        string
	rawUrl        string
	params        []interface{}
	plugins       []Plugin
	index         int
	resolverChain *ResolverChain
	httpClient    *http.Client
}

func (c *Context) Next() error {
	if c.index == 0 {
		if err := c.prepareRequest(); err != nil {
			return err
		}
	}

	if c.index < len(c.plugins) {
		c.plugins[c.index](c)
		c.index++
		return c.Next()
	}

	if err := c.resolveRequest(); err != nil {
		return err
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

func (c *Context) resolveRequest() error {
	for i, param := range c.params {
		context := &RequestContext{Request: c.Request, Params: c.params, Param: param, ParamIndex: i}
		if err := c.resolverChain.Next(context); err != nil {
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
