package sugar

import "net/http"

// Transporter is an important interface for different http clients.
type Transporter interface {
	// Do issues the request.
	Do(req *http.Request) (*http.Response, error)
}
