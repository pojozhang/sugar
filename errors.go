package sugar

import "errors"

var (
	ResolverNotFound = errors.New("resolver not found")
	DecoderNotFound  = errors.New("decoder not found")
)
