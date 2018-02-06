package sugar

import (
	"net/http"
	"net/http/httputil"
	"reflect"
	"log"
	"errors"
)

var (
	resolvers = make(map[reflect.Type]Resolver)
)

type Client struct {
	HttpClient *http.Client
}

var DefaultClient = Client{
	HttpClient: http.DefaultClient,
}

func (c *Client) Get(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodGet, rawUrl, params...)
}

func (c *Client) Post(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPost, rawUrl, params...)
}

func (c *Client) Put(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPut, rawUrl, params...)
}

func (c *Client) Patch(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodPatch, rawUrl, params...)
}

func (c *Client) Delete(rawUrl string, params ...interface{}) (*Response, error) {
	return c.Do(http.MethodDelete, rawUrl, params...)
}

func (c *Client) Do(method, rawUrl string, params ...interface{}) (*Response, error) {
	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return nil, err
	}

	for i, param := range params {
		t := reflect.TypeOf(param)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if r, ok := resolvers[t]; ok {
			if err := r.resolve(req, params, param, i); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("No resolvers found for " + t.String())
		}
	}

	d, _ := httputil.DumpRequest(req, true)
	log.Printf("%s\n", d)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	r := Response(*resp)
	return &r, err
}

func Get(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Get(rawUrl, params...)
}

func Post(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Post(rawUrl, params...)
}

func Put(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Put(rawUrl, params...)
}

func Patch(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Patch(rawUrl, params...)
}

func Delete(rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Delete(rawUrl, params...)
}

func Do(method, rawUrl string, params ...interface{}) (*Response, error) {
	return DefaultClient.Do(method, rawUrl, params...)
}

func GetResolvers() map[reflect.Type]Resolver {
	return resolvers
}
