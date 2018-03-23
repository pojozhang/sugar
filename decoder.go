package sugar

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ResponseContext keeps values for an decoder.
// Out should be the pointer of the object you want to assign.
type ResponseContext struct {
	Request  *http.Request
	Response *http.Response
	Out      interface{}
}

// Decoder converts a response context into a struct.
// It returns an error if any error occurs during decoding.
// Call chain.Next() to propagate context.
type Decoder interface {
	Decode(context *ResponseContext, chain *DecoderChain) error
}

// DecoderChain keeps a set of decoders.
type DecoderChain struct {
	context  *ResponseContext
	decoders []Decoder
	index    int
}

// Next propagates context to next decoder.
// It returns DecoderNotFound if current decoder is the last one.
func (c *DecoderChain) Next() error {
	if c.index < len(c.decoders) {
		c.index++
		return c.decoders[c.index-1].Decode(c.context, c)
	}
	return DecoderNotFound
}

// Add adds decoders to a decoder chain.
func (c *DecoderChain) Add(decoders ...Decoder) *DecoderChain {
	for _, decoder := range decoders {
		c.decoders = append(c.decoders, decoder)
	}
	return c
}

// NewDecoderChain initializes a new decoder chain with given response context and decoders.
func NewDecoderChain(context *ResponseContext, decoders ...Decoder) *DecoderChain {
	chain := &DecoderChain{context: context, index: 0}
	chain.Add(decoders...)
	return chain
}

// JsonDecoder parses JSON-encoded data.
type JsonDecoder struct {
}

// Decode decodes response body via json.Unmarshal.
func (d *JsonDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	for _, contentType := range context.Response.Header[ContentType] {
		if strings.Contains(strings.ToLower(contentType), ContentTypeJson) {
			body, err := ioutil.ReadAll(context.Response.Body)
			if err != nil {
				return err
			}

			err = json.Unmarshal(body, context.Out)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return chain.Next()
}

// XmlDecoder parses XML-encoded data.
type XmlDecoder struct {
}

// Decode decodes response body via xml.Unmarshal.
func (d *XmlDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	for _, contentType := range context.Response.Header[ContentType] {
		if strings.Contains(strings.ToLower(contentType), ContentTypeXml) {
			body, err := ioutil.ReadAll(context.Response.Body)
			if err != nil {
				return err
			}

			err = xml.Unmarshal(body, context.Out)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return chain.Next()
}

// XmlDecoder parses plain text.
type PlainTextDecoder struct {
}

// Decode reads a byte slice from response body via ioutil.ReadAll and then converts it to a string.
func (d *PlainTextDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	if contentTypes, ok := context.Response.Header[ContentType]; ok {
		for _, contentType := range contentTypes {
			if strings.Contains(strings.ToLower(contentType), ContentTypePlainText) {
				goto DECODE
			}
		}

		return chain.Next()
	}

DECODE:
	body, err := ioutil.ReadAll(context.Response.Body)
	if err != nil {
		return err
	}

	*(context.Out.(*string)) = string(body)
	return nil
}

// FileDecoder parses binary-encoded data.
type FileDecoder struct {
}

// Decode decodes response body by writing data to a file.
func (d *FileDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	f, ok := context.Out.(*os.File)
	if !ok {
		return chain.Next()
	}

	targetMime := mime.TypeByExtension(filepath.Ext(f.Name()))
	if targetMime == "" {
		targetMime = ContentTypeOctetStream
	}

	for _, contentType := range context.Response.Header[ContentType] {
		if strings.Contains(strings.ToLower(contentType), targetMime) {
			_, err := io.Copy(f, context.Response.Body)
			return err
		}
	}

	return chain.Next()
}
