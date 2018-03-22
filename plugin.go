package sugar

import (
	"log"
	"net/http/httputil"
)

// Plugin works as an interceptor.
// You can get an assembled request from the context.
// If you want to get the response, you should do that in a defer block.
// Must not forget to call c.Next() to propagate the context.
type Plugin func(c *Context) error

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
