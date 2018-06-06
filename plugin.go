package sugar

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
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

// Retryer provides a common policy to retry the request.
func Retryer(attempts int, delay time.Duration, multiplier float32, maxDelay time.Duration) func(c *Context) error {
	return func(c *Context) (err error) {
		for d, i := delay, 0; i < attempts; i++ {
			err = c.Next()
			if c.Response != nil && c.Response.StatusCode < http.StatusInternalServerError {
				return
			}

			if _, ok := err.(net.Error); err == nil || !ok {
				return
			}

			if i < attempts-1 {
				time.Sleep(d)
				if t := d * time.Duration(multiplier); maxDelay < t {
					d = maxDelay
				} else {
					d = t
				}
			}
		}

		return
	}
}
