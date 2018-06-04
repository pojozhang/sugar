package sugar

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

type mockTransporter struct {
	count int
}

func (t *mockTransporter) Do(req *http.Request) (*http.Response, error) {
	t.count++
	return nil, &net.OpError{}
}

func TestRetryerShouldRetryIfTransporterReturnsAnErrorOfNetError(t *testing.T) {
	const attempts = 3
	p := PluginFunc(Retryer(attempts, time.Duration(1)*time.Second, 1.5, time.Duration(3)*time.Second))
	m := &mockTransporter{}
	c := Context{
		transporter: m,
		plugins:     []Plugin{p},
	}

	c.Next()

	assert.Equal(t, attempts, m.count)
}
