package sugar

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDecoderChain_Next_Returns_Error_If_No_Decoders(t *testing.T) {
	c := &DecoderChain{}
	assert.Equal(t, DecoderNotFound, c.Next())
}

func TestJsonDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeJson}}, Body: errorBody{}}}
	decoder := &JsonDecoder{}
	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestXmlDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeXml}}, Body: errorBody{}}}
	decoder := &XmlDecoder{}
	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestPlainTextDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypePlainText}}, Body: errorBody{}}}
	decoder := &PlainTextDecoder{}
	assert.NotNil(t, decoder.Decode(context, nil))
}

type errorBody struct {
}

func (b errorBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (b errorBody) Close() error {
	return nil
}
