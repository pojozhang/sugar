package sugar

import (
	"testing"
	"encoding/json"
	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"io/ioutil"
)

type book struct {
	Name string
}

func TestGet(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		Reply(200).
		JSON(`[{"name":"bookA"},{"name":"bookB"}]`)

	resp, err := Get("http://api.example.com/books")

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestGetWithQueryVar(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		println(request.URL.String())
		return request.URL.Query()["name"][0] == "bookA", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200).
		JSON(`[{"name":"bookA"}]`)

	resp, err := Get("http://api.example.com/books", Query{"name": "bookA"})

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestGetWithPathVar(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return strings.Contains(request.URL.Path, "123"), nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200).
		JSON(`[{"name":"bookA"}]`)

	resp, err := Get("http://api.example.com/books/:id", Path{"id": 123})

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestDelete(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return strings.Contains(request.URL.Path, "123"), nil
	})
	gock.New("http://api.example.com").
		Delete("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Delete("http://api.example.com/books/:id", Path{"id": 123})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostJSON(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		var book book
		json.Unmarshal(b, &book)
		return request.Header[ContentType][0] == ContentTypeJson && book.Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(201)

	resp, err := Post("http://api.example.com/books", Json(`{"Name":"bookA"}`))

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestPostForm(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Header[ContentType][0] == ContentTypeForm && request.Form["name"][0] == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Post("http://api.example.com/books", Form{"name": "bookA"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestWriteCookies(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		c, _ := request.Cookie("name")
		return c.Value == "bookA", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Get("http://api.example.com/books", Cookie{"name": "bookA"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestWriteHeaders(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Header.Get("name") == "bookA", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Get("http://api.example.com/books", Header{"name": "bookA"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
