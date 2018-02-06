package sugar

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	assert.Equal(t, "8", ToString(uint(8)))
	assert.Equal(t, "8", ToString(uint8(8)))
	assert.Equal(t, "8", ToString(uint16(8)))
	assert.Equal(t, "8", ToString(uint32(8)))
	assert.Equal(t, "8", ToString(uint64(8)))
	assert.Equal(t, "8", ToString(int(8)))
	assert.Equal(t, "8", ToString(int8(8)))
	assert.Equal(t, "8", ToString(int16(8)))
	assert.Equal(t, "8", ToString(int32(8)))
	assert.Equal(t, "8", ToString(int64(8)))
	assert.Equal(t, "8.001", ToString(float32(8.001)))
	assert.Equal(t, "8.00000000001", ToString(float64(8.00000000001)))
	assert.Equal(t, "8", ToString("8"))
	assert.Equal(t, "true", ToString(true))
	assert.Equal(t, "", ToString(struct{}{}))
}
