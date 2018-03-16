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

func (r *Response) Read(v interface{}) (*Response, error) {
	defer r.Close()

	response, err := r.Raw()
	if err != nil {
		return r, err
	}

	err = decoderGroup.Decode(&ResponseContext{Request: r.request, Response: response, Param: v})
	if err != nil {
		return r, err
	}

	return r, err
}

func (r *Response) ReadBytes() ([]byte, *Response, error) {
	defer r.Close()

	_, err := r.Raw()
	if err != nil {
		return nil, r, err
	}

	bytes, err := ioutil.ReadAll(r.Body)
	return bytes, r, err
}

func (r *Response) Close() {
	if r != nil {
		r.Body.Close()
	}
}
