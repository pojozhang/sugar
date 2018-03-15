package sugar

import (
	"testing"
	"encoding/json"
	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/base64"
	"os"
	"encoding/xml"
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

	var books []book
	err := Get("http://api.example.com/books").Read(&books)

	assert.Nil(t, err)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestGetWithQueryVar(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.URL.Query()["name"][0] == "bookA", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200).
		JSON(`[{"name":"bookA"}]`)

	bytes, err := Get("http://api.example.com/books", Query{"name": "bookA"}).ReadBytes()

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(bytes, &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestGetWithQueryList(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		q := request.URL.Query()
		return q["name"][0] == "bookA" && q["name"][1] == "bookB", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200).
		JSON(`[{"name":"bookA"},{"name":"bookB"}]`)

	bytes, err := Get("http://api.example.com/books", Query{"name": List{"bookA", "bookB"}}).ReadBytes()

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(bytes, &books)
	assert.Equal(t, "bookA", books[0].Name)
	assert.Equal(t, "bookB", books[1].Name)
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

	bytes, err := Get("http://api.example.com/books/:id", Path{"id": 123}).ReadBytes()

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(bytes, &books)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestPostJsonString(t *testing.T) {
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

	resp, err := Post("http://api.example.com/books", Json{`{"name":"bookA"}`}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestPostJsonMap(t *testing.T) {
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

	resp, err := Post("http://api.example.com/books", Json{Map{"name": "bookA"}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestPostJsonList(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		var book []book
		json.Unmarshal(b, &book)
		return request.Header[ContentType][0] == ContentTypeJson && book[0].Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(201)

	resp, err := Post("http://api.example.com/books", Json{List{Map{"name": "bookA"}}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestPostBadJson(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Post("/books").
		Reply(200)

	badValue := make(chan int)
	_, err := Post("http://api.example.com/books", Json{badValue}).Raw()

	assert.NotNil(t, err)
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

	resp, err := Post("http://api.example.com/books", Form{"name": "bookA"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostFormList(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Header[ContentType][0] == ContentTypeForm && request.Form["name"][0] == "bookA" && request.Form["name"][1] == "bookB", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Post("http://api.example.com/books", Form{"name": List{"bookA", "bookB"}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostFormWithBadUrl(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Post("/books").
		Reply(200)

	_, err := Post("http://api.example.com/books?%%", Form{"name": "bookA"}).Raw()

	assert.NotNil(t, err)
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

	resp, err := Get("http://api.example.com/books", Cookie{"name": "bookA"}).Raw()

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

	resp, err := Get("http://api.example.com/books", Header{"name": "bookA"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
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

	resp, err := Delete("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPut(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Put("/books/123").
		Reply(204)

	resp, err := Put("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestPatch(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Path("/books/123").
		Reply(204)

	resp, err := Patch("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestBasicAuth(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		s := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))
		return strings.Contains(request.Header["Authorization"][0], s), nil
	})
	gock.New("http://api.example.com").
		Delete("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Delete("http://api.example.com/books", User{"user", "password"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDo(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Method == http.MethodTrace, nil
	})
	gock.New("http://api.example.com").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Do(http.MethodTrace, "http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestWrongUrl(t *testing.T) {
	_, err := Patch("wrong://wrong-url.com").Raw()
	assert.NotNil(t, err)
}

func TestWrongRequest(t *testing.T) {
	_, err := Do("?", "http://wrong-url").Raw()
	assert.NotNil(t, err)
}

func TestNoResolverFound(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Path("/books").
		Reply(200)

	resp, err := Get("http://api.example.com/books", struct{}{}).Raw()

	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func TestApply(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		s := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))
		return strings.Contains(request.Header["Authorization"][0], s), nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200)

	Apply(User{"user", "password"})
	defer Reset()
	resp, err := Get("http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReset(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return len(request.Header["Authorization"]) == 0, nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(200)

	Apply(User{"user", "password"})
	Reset()
	resp, err := Get("http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostMultiPartWithOsFile(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		request.ParseMultipartForm(32 << 20)
		file, _, _ := request.FormFile("file")
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		return string(b) == "hello sugar!" && request.FormValue("name") == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(200)

	f, _ := os.Open("text")
	defer f.Close()
	resp, err := Post("http://api.example.com/books", MultiPart{"name": "bookA", "file": f}).Raw()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostPlainText(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		return request.Header[ContentType][0] == ContentTypePlainText && string(b) == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Post("http://api.example.com/books", "bookA").Raw()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostXml(t *testing.T) {
	type book struct {
		XMLName xml.Name `xml:"book"`
		Name    string   `xml:"name,attr"`
	}

	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		var book book
		xml.Unmarshal(b, &book)
		return request.Header[ContentType][0] == ContentTypeXml && book.Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(200)

	resp, err := Post("http://api.example.com/books", Xml{`<book name="bookA"></book>`}).Raw()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostBadXml(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Post("/books").
		Reply(200)

	_, err := Post("http://api.example.com/books", Xml{make(chan int)}).Raw()
	assert.NotNil(t, err)
}

func TestMap(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		Reply(200)

	Apply(Mapper{func(req *http.Request) {
		if req.URL.Host == "book-service" {
			req.URL.Host = "api.example.com"
		}
	}})
	resp, err := Get("http://book-service/books").Raw()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
