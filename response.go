package sugar

import (
	"net/http"
	"io/ioutil"
)

type Response http.Response

func (r *Response) ReadBytes() []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}
