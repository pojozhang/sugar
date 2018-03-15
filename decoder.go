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
		if strings.Contains(strings.ToLower(contentType), "application/json") {
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
		if strings.Contains(strings.ToLower(contentType), "application/xml") {
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
	)
}
