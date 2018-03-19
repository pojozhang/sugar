package sugar

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type ResponseContext struct {
	Request  *http.Request
	Response *http.Response
	Param    interface{}
}

type Decoder interface {
	Decode(context *ResponseContext, chain *DecoderChain) error
}

type DecoderChain struct {
	context  *ResponseContext
	decoders []Decoder
	index    int
}

func (c *DecoderChain) Next() error {
	if c.index < len(c.decoders) {
		c.index++
		return c.decoders[c.index-1].Decode(c.context, c)
	}
	return DecoderNotFound
}

func (c *DecoderChain) reset() *DecoderChain {
	c.index = 0
	return c
}

func (c *DecoderChain) Add(decoders ...Decoder) *DecoderChain {
	for _, decoder := range decoders {
		c.decoders = append(c.decoders, decoder)
	}
	return c
}

func NewDecoderChain(context *ResponseContext, decoders ...Decoder) *DecoderChain {
	chain := &DecoderChain{context: context, index: 0}
	chain.Add(decoders...)
	return chain
}

type JsonDecoder struct {
}

func (d *JsonDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	for _, contentType := range context.Response.Header[ContentType] {
		if strings.Contains(strings.ToLower(contentType), ContentTypeJson) {
			body, err := ioutil.ReadAll(context.Response.Body)
			if err != nil {
				return err
			}

			err = json.Unmarshal(body, context.Param)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return chain.Next()
}

type XmlDecoder struct {
}

func (d *XmlDecoder) Decode(context *ResponseContext, chain *DecoderChain) error {
	for _, contentType := range context.Response.Header[ContentType] {
		if strings.Contains(strings.ToLower(contentType), ContentTypeXml) {
			body, err := ioutil.ReadAll(context.Response.Body)
			if err != nil {
				return err
			}

			err = xml.Unmarshal(body, context.Param)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return chain.Next()
}

type PlainTextDecoder struct {
}

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

	*(context.Param.(*string)) = string(body)
	return nil
}
