package sugar

import (
	"net/http"
)

type DecoderContext struct {
	Decoders []Decoder
}

func (c *DecoderContext) Next() error {

}

func (c *DecoderContext) Add(decoder Decoder) {
	c.Decoders = append(c.Decoders, decoder)
}

type Decoder interface {
	Decode(req *http.Request, resp *http.Response, param interface{}, chain DecoderContext) error
}

func init() {

}
