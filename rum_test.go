package rum

import (
	"testing"
)

func TestGet(t *testing.T) {
	Get("https://baidu.com",Q{"qq":"ww"})
}