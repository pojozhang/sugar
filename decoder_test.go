package sugar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecoderChain_Next_Returns_Error_If_No_Decoders(t *testing.T) {
	c := &DecoderChain{}
	assert.Equal(t, DecoderNotFound, c.Next())
}
