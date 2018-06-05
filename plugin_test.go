package sugar

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

type mockTransporter struct {
	count    int
	response *http.Response
	error    error
}

func (t *mockTransporter) Do(req *http.Request) (*http.Response, error) {
	t.count++
	return t.response, t.error
}

func TestRetryer_Retry_If_Transporter_Returns_An_Error_Of_NetError(t *testing.T) {
	const attempts = 3
	p := PluginFunc(Retryer(attempts, time.Duration(1)*time.Second, 1.5, time.Duration(3)*time.Second))
	m := &mockTransporter{response: nil, error: &net.OpError{}}
	c := Context{
		transporter: m,
		plugins:     []Plugin{p},
	}

	c.Next()

	assert.Equal(t, attempts, m.count)
}

func TestRetryer_Not_Retry_If_Transporter_Returns_An_Error_Which_Is_Not_Instance_Of_NetError(t *testing.T) {
	const attempts = 1
	p := PluginFunc(Retryer(attempts, time.Duration(1)*time.Second, 1.5, time.Duration(3)*time.Second))
	m := &mockTransporter{response: nil, error: EncoderNotFound}
	c := Context{
		transporter: m,
		plugins:     []Plugin{p},
	}

	c.Next()

	assert.Equal(t, attempts, m.count)
}

func TestRetryer_Not_Retry_If_Request_Is_Handled_Correctly(t *testing.T) {
	const attempts = 1
	p := PluginFunc(Retryer(attempts, time.Duration(1)*time.Second, 1.5, time.Duration(3)*time.Second))
	m := &mockTransporter{response: &http.Response{StatusCode: http.StatusOK}, error: nil}
	c := Context{
		transporter: m,
		plugins:     []Plugin{p},
	}

	c.Next()

	assert.Equal(t, attempts, m.count)
}
