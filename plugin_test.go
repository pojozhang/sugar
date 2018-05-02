package sugar

import (
	"net"
	"net/http"
	"testing"
	"time"
)

type mockTransporter struct {
}

func (t *mockTransporter) Do(req *http.Request) (*http.Response, error) {
	return nil, &net.OpError{}
}

func TestRetryer(t *testing.T) {
	p := PluginFunc(Retryer(3, time.Duration(1)*time.Second, 1.5, time.Duration(3)*time.Second))
	c := Context{
		transporter: &mockTransporter{},
		plugins:     []Plugin{p},
	}
	c.Next()
}
