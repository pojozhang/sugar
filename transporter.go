package sugar

import "net/http"

type Transporter interface {
	Do(req *http.Request) (*http.Response, error)
}
