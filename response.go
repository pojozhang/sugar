package sugar

import (
	"io/ioutil"
	"net/http"
	"unsafe"
)

// Response represents a context of response
type Response struct {
	http.Response
	Error    error
	request  *http.Request
	decoders []Decoder
}

// Raw returns a raw response and en error.
func (r *Response) Raw() (*http.Response, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return (*http.Response)(unsafe.Pointer(r)), nil
}

// Read decodes response data via decoders.
func (r *Response) Read(v interface{}) (*http.Response, error) {
	defer r.Close()

	resp, err := r.Raw()
	if err != nil {
		return resp, err
	}

	err = NewDecoderChain(&ResponseContext{Request: r.request, Response: resp, Out: v}, r.decoders...).Next()
	if err != nil {
		return resp, err
	}

	return resp, err
}

// ReadBytes reads response body into a byte slice.
func (r *Response) ReadBytes() ([]byte, *http.Response, error) {
	defer r.Close()

	resp, err := r.Raw()
	if err != nil {
		return nil, resp, err
	}

	bytes, err := ioutil.ReadAll(r.Body)
	return bytes, resp, err
}

// Close closes response body.
func (r *Response) Close() {
	if r != nil && r.Body != nil {
		r.Body.Close()
	}
}
