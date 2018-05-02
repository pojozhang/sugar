package sugar

import (
	"log"
	"net/http/httputil"
)

// Plugin works as an interceptor.
// You can get an assembled request from the context.
// If you want to do something after you get the response, you should do that in a defer block.
// Must not forget to call c.Next() to propagate the context.
type Plugin interface {
	// Handle includes custom logic.
	Handle(c *Context) error
}

// PluginFunc is a function adapter for Plugin.
type PluginFunc func(c *Context) error

// Handle includes custom logic.
func (f PluginFunc) Handle(c *Context) error {
	return f(c)
}

// Logger is a builtin plugin to dump requests and responses.
func Logger(c *Context) error {
	b, _ := httputil.DumpRequest(c.Request, true)
	log.Println(string(b))
	defer func() {
		if c.Response != nil {
			b, _ := httputil.DumpResponse(c.Response, true)
			log.Println(string(b))
		}
	}()
	return c.Next()
}
