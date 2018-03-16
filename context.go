package sugar

import "net/http"

type Context struct {
	Request     *http.Request
	Response    *http.Response
	Params      []interface{}
	Plugins     []*Plugin
	pluginIndex uint8
}

func (c *Context) Next() {
	c.pluginIndex++
}
