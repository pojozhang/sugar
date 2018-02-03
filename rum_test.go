package rum

import "testing"

func TestGet(t *testing.T) {
	Get("https://",Q{"":[]string{"ss"}})
}