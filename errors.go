package sugar

import "errors"

var (
	EncoderNotFound = errors.New("encoder not found")
	DecoderNotFound = errors.New("decoder not found")
)
