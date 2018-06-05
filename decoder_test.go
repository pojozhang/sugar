package sugar

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestDecoderChain_Next_Returns_Error_If_No_Decoders(t *testing.T) {
	c := &DecoderChain{}

	assert.Equal(t, DecoderNotFound, c.Next())
}

func TestJsonDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeJson}}, Body: readErrorBody{}}}
	decoder := &JsonDecoder{}

	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestJsonDecoder_Decode_Returns_Error_If_Fail_To_Unmarshal(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeJson}}, Body: unmarshalErrorBody{}}}
	decoder := &JsonDecoder{}

	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestXmlDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeXml}}, Body: readErrorBody{}}}
	decoder := &XmlDecoder{}

	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestXmlDecoder_Decode_Returns_Error_If_Fail_To_Unmarshal(t *testing.T) {
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeXml}}, Body: unmarshalErrorBody{}}}
	decoder := &XmlDecoder{}

	assert.NotNil(t, decoder.Decode(context, nil))
}

func TestPlainTextDecoder_Decode_Returns_Error_If_Fail_To_Read_Body(t *testing.T) {
	var out string
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypePlainText}}, Body: readErrorBody{}}, Out: &out}
	decoder := &PlainTextDecoder{}

	assert.NotNil(t, decoder.Decode(context, nil))
}

type mockDecoder struct {
	Called bool
}

func (d *mockDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	d.Called = true
	return chain.Next()
}

func TestPlainTextDecoder_Will_Propagate_If_Out_Is_Not_Type_Of_String_Pointer(t *testing.T) {
	var out int
	context := &ResponseContext{Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypePlainText}}, Body: readErrorBody{}}, Out: &out}
	decoder := &PlainTextDecoder{}
	nextDecoder := &mockDecoder{Called: false}

	decoder.Decode(context, NewDecoderChain(&ResponseContext{Response: nil, Out: nil}, nextDecoder))

	assert.True(t, nextDecoder.Called)
}

func TestFileDecoder_Decode_Returns_Error_If_Out_Is_Not_Ptr_Of_OsFile(t *testing.T) {
	context := &ResponseContext{Out: "string"}
	chain := &DecoderChain{}
	decoder := &FileDecoder{}

	assert.Equal(t, DecoderNotFound, decoder.Decode(context, chain))
}

type readErrorBody struct {
}

func (b readErrorBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (b readErrorBody) Close() error {
	return nil
}

type unmarshalErrorBody struct {
}

func (b unmarshalErrorBody) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (b unmarshalErrorBody) Close() error {
	return nil
}

func TestPlainTextDecoder_Decode_Returns_Error_If_Content_Type_Is_Wrong(t *testing.T) {
	s := new(string)
	context := &ResponseContext{Out: s, Response: &http.Response{Header: http.Header{ContentType: []string{ContentTypeForm}}}}
	chain := &DecoderChain{}
	decoder := &PlainTextDecoder{}

	assert.Equal(t, DecoderNotFound, decoder.Decode(context, chain))
}
