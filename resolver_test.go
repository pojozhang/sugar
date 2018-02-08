package sugar

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
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
	new(PathResolver).resolve(req, L{P{"id": "golang"}}, P{"id": "golang"}, 0)
	assert.Equal(t, "http://github.com/golang/", req.URL.String())
}

func TestResolveJsonBytes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://github.com", nil)
	m := map[string]string{"k": "v"}
	b, _ := json.Marshal(m)
	new(JsonResolver).resolve(req, L{J(b)}, J(b), 0)
	b, _ = ioutil.ReadAll(req.Body)
	var n map[string]*json.RawMessage
	json.Unmarshal(b, n)
	assert.Equal(t, "v", m["k"])
}

func TestResolveJsonList(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://github.com", nil)
	new(JsonResolver).resolve(req, L{J(L{1, 2, 3})}, J(L{1, 2, 3}), 0)
	b, _ := ioutil.ReadAll(req.Body)
	n := make([]int, 0)
	json.Unmarshal(b, &n)
	assert.Equal(t, 1, n[0])
	assert.Equal(t, 2, n[1])
	assert.Equal(t, 3, n[2])
}

func TestReturnErrorIfWriterIsNull(t *testing.T) {
	f, _ := os.Open("file")
	err := writeFile(nil, "file", f)
	assert.NotNil(t, err)
}
