package sugar

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
)

var (
	decoderGroup DecoderGroup
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
	decoder   Decoder
	successor *DecoderChain
}

func (c *DecoderChain) Next(context *ResponseContext) error {
	if c != nil {
		return c.decoder.Decode(context, c.successor)
	}
	return DecoderNotFound
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

type DecoderGroup struct {
	chain []*DecoderChain
}

func (g *DecoderGroup) Add(decoders ... Decoder) {
	for _, decoder := range decoders {
		c := &DecoderChain{decoder: decoder}
		g.chain = append(g.chain, c)
	}
	g.refresh()
}

func (g *DecoderGroup) refresh() {
	for i := 0; i < len(g.chain)-1; i++ {
		g.chain[i].successor = g.chain[i+1]
	}
}

func (g *DecoderGroup) Decode(context *ResponseContext) error {
	if len(g.chain) < 1 {
		return nil
	}
	return g.chain[0].Next(context)
}

func RegisterDecoders(decoders ... Decoder) {
	decoderGroup.Add(decoders...)
}

func init() {
	RegisterDecoders(
		&JsonDecoder{},
		&XmlDecoder{},
		&PlainTextDecoder{},
	)
}
