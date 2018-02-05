package sugar

import (
	"testing"
	"encoding/json"
	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
)

type book struct {
	Name string
}

func TestGetBooks(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		Reply(200).
		JSON(`[{"name":"bookA"},{"name":"bookB"}]`)

	resp, err := Get("http://api.example.com/books")

	if err != nil {
		t.Fatalf("%v", err)
	}

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestFindBooksByName(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		AddMatcher(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.URL.Query()["name"][0] == "bookA", nil
	}).Reply(200).JSON(`[{"name":"bookA"}]`)

	resp, err := Get("http://api.example.com/books", Query{"name": "bookA"})

	if err != nil {
		t.Fatalf("%v", err)
	}

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestFindBookById(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		AddMatcher(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return strings.Contains(request.URL.Path, "123"), nil
	}).Reply(200).JSON(`[{"name":"bookA"}]`)

	resp, err := Get("http://api.example.com/books/:id", Path{"id": 123})

	if err != nil {
		t.Fatalf("%v", err)
	}

	var books []book
	json.Unmarshal(resp.ReadBytes(), &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestDeleteBookById(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Delete("/books").
		AddMatcher(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return strings.Contains(request.URL.Path, "123"), nil
	}).Reply(200)

	resp, err := Delete("http://api.example.com/books/:id", Path{"id": 123})

	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestSendCookies(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		AddMatcher(func(request *http.Request, request2 *gock.Request) (bool, error) {
		c, _ := request.Cookie("name")
		return c.Value == "sugar", nil
	}).Reply(200)

	resp, err := Get("http://api.example.com/books", Cookie{"name": "sugar"})

	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}
