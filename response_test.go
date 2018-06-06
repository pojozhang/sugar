package sugar

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse_Read_Returns_Error_When_Response_Has_Error(t *testing.T) {
	var v string
	resp := &Response{Error: EncoderNotFound}
	_, err := resp.Read(&v)
	assert.NotNil(t, err)
}

func TestResponse_ReadBytes_Returns_Error_When_Response_Has_Error(t *testing.T) {
	resp := &Response{Error: EncoderNotFound}
	_, _, err := resp.ReadBytes()
	assert.NotNil(t, err)
}

func TestResponse_Read_Returns_Error_If_DecoderChain_Returns_Error(t *testing.T) {
	resp := &Response{decoders: DecoderGroup{&errorDecoder{}}}
	var v string
	_, err := resp.Read(&v)
	assert.NotNil(t, err)
}

type errorDecoder struct {
}

func (d *errorDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	return errors.New("error")
}
