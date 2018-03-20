package sugar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse_Read_When_Response_Has_Error(t *testing.T) {
	var v string
	resp := &Response{Error: EncoderNotFound}
	_, err := resp.Read(&v)
	assert.NotNil(t, err)
}

func TestResponse_ReadBytes_When_Response_Has_Error(t *testing.T) {
	resp := &Response{Error: EncoderNotFound}
	_, _, err := resp.ReadBytes()
	assert.NotNil(t, err)
}
