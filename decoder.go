package sugar

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
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
	decoders []Decoder
	index    int
}

func (c *DecoderChain) Next(context *ResponseContext) error {
	if c.index < len(c.decoders) {
		defer func() { c.index++ }()
		return c.decoders[c.index].Decode(context, c)
	}
	return DecoderNotFound
}

func (c *DecoderChain) Reset() *DecoderChain {
	c.decoders = []Decoder{}
	c.index = 0
	return c
}

func (c *DecoderChain) Add(decoders ... Decoder) *DecoderChain {
	for _, decoder := range decoders {
		c.decoders = append(c.decoders, decoder)
	}
	return c
}

func NewDecoderChain(decoders ... Decoder) *DecoderChain {
	chain := &DecoderChain{}
	chain.Reset()
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

	return chain.Next(context)
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

	return chain.Next(context)
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

		return chain.Next(context)
	}

DECODE:
	body, err := ioutil.ReadAll(context.Response.Body)
	if err != nil {
		return err
	}

	*(context.Param.(*string)) = string(body)
	return nil
}
