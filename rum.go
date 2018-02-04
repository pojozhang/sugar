package rum

import (
	"net/http"
	"net/url"
	"reflect"
	"log"
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

func (c *Client) Get(rawUrl string, params ...interface{}) (*http.Response, error) {
	return c.Do(http.MethodGet, rawUrl, params)
}

func (c *Client) Post(rawUrl string, params ...interface{}) (*http.Response, error) {
	return c.Do(http.MethodPost, rawUrl, params)
}

func (c *Client) Put(rawUrl string, params ...interface{}) (*http.Response, error) {
	return c.Do(http.MethodPut, rawUrl, params)
}

func (c *Client) Patch(rawUrl string, params ...interface{}) (*http.Response, error) {
	return c.Do(http.MethodPatch, rawUrl, params)
}

func (c *Client) Delete(rawUrl string, params ...interface{}) (*http.Response, error) {
	return c.Do(http.MethodDelete, rawUrl, params)
}

func (c *Client) Do(method, rawUrl string, params ...interface{}) (*http.Response, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	req := &http.Request{URL: u}
	for i, param := range params {
		if r, ok := resolvers[reflect.TypeOf(param)]; ok {
			r.resolve(req, params, param, i)
		}
	}
	req.URL.RawQuery = req.URL.Query().Encode()

	log.Printf("%s %s/n", method, req.URL.String())
	return c.HttpClient.Do(req)
}

func Get(rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Get(rawUrl, params)
}

func Post(rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Post(rawUrl, params)
}

func Put(rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Put(rawUrl, params)
}

func Patch(rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Patch(rawUrl, params)
}

func Delete(rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Delete(rawUrl, params)
}

func Do(method, rawUrl string, params ...interface{}) (*http.Response, error) {
	return DefaultClient.Do(method, rawUrl, params)
}

func GetResolvers() map[reflect.Type]Resolver {
	return resolvers
}

func init() {
	resolvers[reflect.TypeOf(Query{})] = &QueryResolver{}
	resolvers[reflect.TypeOf(Header{})] = &QueryResolver{}
	resolvers[reflect.TypeOf(Json{})] = &QueryResolver{}
	resolvers[reflect.TypeOf(Form{})] = &QueryResolver{}
}
