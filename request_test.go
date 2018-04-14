package sugar

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type book struct {
	Name string
}

func init() {
	Use(Logger)
}

func TestGetText(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/echo").
		Reply(http.StatusOK).
		BodyString("sugar")

	var text = new(string)
	_, err := Get("http://api.example.com/echo").Read(text)

	assert.Nil(t, err)
	assert.Equal(t, "sugar", *text)
}

func TestGetPlainText(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/echo").
		Reply(http.StatusOK).
		AddHeader(ContentType, ContentTypePlainText).
		BodyString("sugar")

	var text = new(string)
	resp, err := Get("http://api.example.com/echo").Read(text)

	assert.Nil(t, err)
	assert.Contains(t, resp.Header[ContentType], ContentTypePlainText)
	assert.Equal(t, "sugar", *text)
}

func TestGetJson(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		Reply(http.StatusOK).
		JSON(`[{"name":"bookA"},{"name":"bookB"}]`)

	var books []book
	_, err := Get("http://api.example.com/books").Read(&books)

	assert.Nil(t, err)
	assert.Equal(t, "bookA", books[0].Name)
}

func TestGetXml(t *testing.T) {
	type book struct {
		XMLName xml.Name `xml:"book"`
		Name    string   `xml:"name,attr"`
	}

	defer gock.Off()
	gock.New("http://api.example.com").
		Get("/books").
		Reply(http.StatusOK).
		XML(`<book name="bookA"></book>`)

	var b book
	_, err := Get("http://api.example.com/books").Read(&b)

	assert.Nil(t, err)
	assert.Equal(t, "bookA", b.Name)
}

func TestGetWithQueryPair(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.URL.Query()["name"][0] == "bookA", nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(http.StatusOK).
		JSON(`[{"name":"bookA"}]`)

	bytes, _, err := Get("http://api.example.com/books", Query{"name": "bookA"}).ReadBytes()

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
		Reply(http.StatusOK).
		JSON(`[{"name":"bookA"},{"name":"bookB"}]`)

	bytes, _, err := Get("http://api.example.com/books", Query{"name": List{"bookA", "bookB"}}).ReadBytes()

	assert.Nil(t, err)

	var books []book
	json.Unmarshal(bytes, &books)
	assert.Equal(t, "bookA", books[0].Name)
	assert.Equal(t, "bookB", books[1].Name)
}

func TestGetWithPathVariable(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return strings.Contains(request.URL.Path, "123"), nil
	})
	gock.New("http://api.example.com").
		Get("/books").
		SetMatcher(matcher).
		Reply(http.StatusOK).
		JSON(`[{"name":"bookA"}]`)

	bytes, _, err := Get("http://api.example.com/books/:id", Path{"id": 123}).ReadBytes()

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
		return request.Header[ContentType][0] == ContentTypeJsonUtf8 && book.Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(http.StatusCreated)

	resp, err := Post("http://api.example.com/books", Json{`{"name":"bookA"}`}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestPostJsonPair(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		var book book
		json.Unmarshal(b, &book)
		return request.Header[ContentType][0] == ContentTypeJsonUtf8 && book.Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(http.StatusCreated)

	resp, err := Post("http://api.example.com/books", Json{Map{"name": "bookA"}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestPostJsonList(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		b, _ := ioutil.ReadAll(request.Body)
		var book []book
		json.Unmarshal(b, &book)
		return request.Header[ContentType][0] == ContentTypeJsonUtf8 && book[0].Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(http.StatusCreated)

	resp, err := Post("http://api.example.com/books", Json{List{Map{"name": "bookA"}}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestPostBadJson(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Post("/books").
		Reply(http.StatusOK)

	badValue := make(chan int)
	_, err := Post("http://api.example.com/books", Json{badValue}).Raw()

	assert.NotNil(t, err)
}

func TestPostFormPair(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Header[ContentType][0] == ContentTypeForm && request.Form["name"][0] == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(http.StatusOK)

	resp, err := Post("http://api.example.com/books", Form{"name": "bookA"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Post("http://api.example.com/books", Form{"name": List{"bookA", "bookB"}}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Get("http://api.example.com/books", Cookie{"name": "bookA"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Get("http://api.example.com/books", Header{"name": "bookA"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Delete("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPut(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Put("/books/123").
		Reply(http.StatusNoContent)

	resp, err := Put("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPatch(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Path("/books/123").
		Reply(http.StatusNoContent)

	resp, err := Patch("http://api.example.com/books/:id", Path{"id": 123}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Delete("http://api.example.com/books", User{"user", "password"}).Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDo(t *testing.T) {
	defer gock.Off()
	matcher := gock.NewBasicMatcher()
	matcher.Add(func(request *http.Request, request2 *gock.Request) (bool, error) {
		return request.Method == http.MethodTrace, nil
	})
	gock.New("http://api.example.com").
		SetMatcher(matcher).
		Reply(http.StatusOK)

	resp, err := Do(http.MethodTrace, "http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBadUrl(t *testing.T) {
	_, err := Patch("wrong://wrong-url.com").Raw()
	assert.NotNil(t, err)
}

func TestBadRequest(t *testing.T) {
	_, err := Do("?", "http://wrong-url").Raw()
	assert.NotNil(t, err)
}

func TestNoEncoderFound(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Path("/books").
		Reply(http.StatusOK)

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
		Reply(http.StatusOK)

	Apply(User{"user", "password"})
	defer Reset()
	resp, err := Get("http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	Apply(User{"user", "password"})
	Reset()
	resp, err := Get("http://api.example.com/books").Raw()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostMultiPart(t *testing.T) {
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
		Reply(http.StatusOK)

	f, _ := os.Open("text")
	defer f.Close()
	resp, err := Post("http://api.example.com/books", MultiPart{"name": "bookA", "file": f}).Raw()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		Reply(http.StatusOK)

	resp, err := Post("http://api.example.com/books", "bookA").Raw()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
		return request.Header[ContentType][0] == ContentTypeXmlUtf8 && book.Name == "bookA", nil
	})
	gock.New("http://api.example.com").
		Post("/books").
		SetMatcher(matcher).
		Reply(http.StatusOK)

	resp, err := Post("http://api.example.com/books", Xml{`<book name="bookA"></book>`}).Raw()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostBadXml(t *testing.T) {
	defer gock.Off()
	gock.New("http://api.example.com").
		Post("/books").
		Reply(http.StatusOK)

	_, err := Post("http://api.example.com/books", Xml{make(chan int)}).Raw()
	assert.NotNil(t, err)
}

func TestDownloadFile(t *testing.T) {
	defer gock.Off()
	fileName := "logo.png"
	gock.New("http://api.example.com").
		Get("/"+fileName).
		Reply(http.StatusOK).
		File(fileName).
		AddHeader(ContentType, mime.TypeByExtension(filepath.Ext(fileName)))

	f, _ := os.Create("tmp.png")
	defer func() {
		f.Close()
		os.Remove("tmp.png")
	}()

	_, err := Get("http://api.example.com/:file", P{"file": fileName}).Read(f)
	assert.Nil(t, err)

	actualStat, _ := f.Stat()
	assert.True(t, actualStat.Size() > 0)

	expected, _ := os.Open(fileName)
	expectedStat, _ := expected.Stat()
	assert.Equal(t, expectedStat.Size(), actualStat.Size())
}

func TestDownloadFileWithUnknownExt(t *testing.T) {
	defer gock.Off()
	fileName := "logo.png"
	gock.New("http://api.example.com").
		Get("/"+fileName).
		Reply(http.StatusOK).
		File(fileName).
		AddHeader(ContentType, ContentTypeOctetStream)

	f, _ := os.Create("tmp.image")
	defer func() {
		f.Close()
		os.Remove("tmp.image")
	}()

	_, err := Get("http://api.example.com/:file", P{"file": fileName}).Read(f)
	assert.Nil(t, err)

	actualStat, _ := f.Stat()
	assert.True(t, actualStat.Size() > 0)

	expected, _ := os.Open(fileName)
	expectedStat, _ := expected.Stat()
	assert.Equal(t, expectedStat.Size(), actualStat.Size())
}

func TestClient_NewRequest(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://api.example.com/books/:id", Path{"id": 1})
	assert.Nil(t, err)
	assert.Equal(t, "http://api.example.com/books/1", req.URL.String())
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	assert.NotNil(t, client.HttpClient)
	assert.Equal(t, *Encoders, client.Encoders)
	assert.Equal(t, *Decoders, client.Decoders)
}
