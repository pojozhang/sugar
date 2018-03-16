package sugar

import (
	"net/http"
	"io/ioutil"
	"unsafe"
)

type Response struct {
	http.Response
	Error   error
	request *http.Request
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

	err = decoderGroup.Decode(&ResponseContext{Request: r.request, Response: resp, Param: v})
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
