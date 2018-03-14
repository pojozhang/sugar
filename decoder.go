package sugar

import (
	"net/http"
)

var (
	decoderGroup DecoderGroup
)

type Decoder interface {
	Decode(req *http.Request, resp *http.Response, param interface{}) error
}

type ChainedDecoder struct {
	Decoder
	successor Decoder
}

func (d *ChainedDecoder) Next(req *http.Request, resp *http.Response, param interface{}) error {
	if d.successor != nil {
		return d.successor.Decode(req, resp, param)
	}
	return nil
}

type JsonDecoder struct {
	ChainedDecoder
}

func (d *JsonDecoder) Decode(req *http.Request, resp *http.Response, param interface{}) error {
	return nil
}

type DecoderGroup struct {
	decoders []*ChainedDecoder
}

func (g *DecoderGroup) Add(decoders ... Decoder) {
	for _, decoder := range decoders {
		g.decoders = append(g.decoders, &ChainedDecoder{Decoder: decoder})
	}
	g.refresh()
}

func (g *DecoderGroup) refresh() {
	for i := 0; i < len(g.decoders)-1; i++ {
		g.decoders[i].successor = g.decoders[i+1]
	}
}

func (g *DecoderGroup) Decode(req *http.Request, resp *http.Response, param interface{}) error {
	if len(g.decoders) < 1 {
		return nil
	}
	return g.decoders[0].Decode(req, resp, param)
}

func RegisterDecoders(decoders ... Decoder) {
	decoderGroup.Add(decoders...)
}

func init() {
	RegisterDecoders(&JsonDecoder{})
}
