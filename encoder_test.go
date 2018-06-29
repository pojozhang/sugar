package sugar

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestToString(t *testing.T) {
	assert.Equal(t, "8", ToString(uint(8)))
	assert.Equal(t, "8", ToString(uint8(8)))
	assert.Equal(t, "8", ToString(uint16(8)))
	assert.Equal(t, "8", ToString(uint32(8)))
	assert.Equal(t, "8", ToString(uint64(8)))
	assert.Equal(t, "8", ToString(int(8)))
	assert.Equal(t, "8", ToString(int8(8)))
	assert.Equal(t, "8", ToString(int16(8)))
	assert.Equal(t, "8", ToString(int32(8)))
	assert.Equal(t, "8", ToString(int64(8)))
	assert.Equal(t, "8.001", ToString(float32(8.001)))
	assert.Equal(t, "8.00000000001", ToString(float64(8.00000000001)))
	assert.Equal(t, "8", ToString("8"))
	assert.Equal(t, "true", ToString(true))
	assert.Equal(t, "", ToString(struct{}{}))
}

func TestResolvePath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://github.com/:id/", nil)

	new(PathEncoder).Encode(&RequestContext{Request: req, Params: L{P{"id": "golang"}}, Param: P{"id": "golang"}, ParamIndex: 0}, nil)

	assert.Equal(t, "http://github.com/golang/", req.URL.String())
}

func TestResolveJsonBytes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://github.com", nil)
	m := map[string]string{"k": "v"}
	b, _ := json.Marshal(m)

	new(JsonEncoder).Encode(&RequestContext{Request: req, Params: L{J{b}}, Param: J{b}, ParamIndex: 0}, nil)

	b, _ = ioutil.ReadAll(req.Body)
	var n map[string]*json.RawMessage
	json.Unmarshal(b, n)

	assert.Equal(t, "v", m["k"])
}

func TestResolveJsonList(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://github.com", nil)

	new(JsonEncoder).Encode(&RequestContext{Request: req, Params: L{J{L{1, 2, 3}}}, Param: J{L{1, 2, 3}}, ParamIndex: 0}, nil)

	b, _ := ioutil.ReadAll(req.Body)
	n := make([]int, 0)
	json.Unmarshal(b, &n)

	assert.Equal(t, 1, n[0])
	assert.Equal(t, 2, n[1])
	assert.Equal(t, 3, n[2])
}

type mockReader struct {
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock reader error")
}

func TestWriteFile_Returns_Error_When_File_Reader_Returns_Error(t *testing.T) {
	err := writeFile(multipart.NewWriter(os.Stdout), "file", "file", &mockReader{})

	assert.NotNil(t, err)
}

type mockWriter struct {
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return -1, os.ErrClosed
}

func TestWriteFile_Returns_Error_When_Fails_To_Create_File_Writer(t *testing.T) {
	err := writeFile(multipart.NewWriter(&mockWriter{}), "file", "file", nil)

	assert.NotNil(t, err)
}
