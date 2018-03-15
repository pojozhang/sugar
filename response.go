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

func (r *Response) Read(v interface{}) (error) {
	defer r.Close()

	if r.Error != nil {
		return r.Error
	}

	response, err := r.Raw()
	if err != nil {
		return err
	}

	err = decoderGroup.Decode(&ResponseContext{Request: r.request, Response: response, Param: v})
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) ReadBytes() ([]byte, error) {
	defer r.Close()

	if r.Error != nil {
		return nil, r.Error
	}

	return ioutil.ReadAll(r.Body)
}

func (r *Response) Close() {
	if r != nil {
		r.Body.Close()
	}
}
