package sugar

import (
	"net/http/httputil"
	"log"
)

type Plugin func(c *Context) error

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
