package sugar

import (
	"io/ioutil"
	"net/http"
	"unsafe"
)

type Response struct {
	http.Response
	Error    error
	request  *http.Request
	decoders []Decoder
}

func (r *Response) Raw() (*http.Response, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return (*http.Response)(unsafe.Pointer(r)), nil
}

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

func (r *Response) ReadBytes() ([]byte, *http.Response, error) {
	defer r.Close()

	resp, err := r.Raw()
	if err != nil {
		return nil, resp, err
	}

	bytes, err := ioutil.ReadAll(r.Body)
	return bytes, resp, err
}

func (r *Response) Close() {
	if r != nil {
		r.Body.Close()
	}
}
