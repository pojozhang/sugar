package rum

import "net/http/httptest"

type MapParams map[string]interface{}

type Header MapParams

type H = Header

type Path MapParams

type P = Path

type Query MapParams

type Q = Query

type Form MapParams

type F = Form

type Json MapParams

type J = Json

func Get(url string, params ...interface{}) {

}

func Post(url string, params ...interface{}) {

}

func Put(url string, params ...interface{}) {

}

func Patch(url string, params ...interface{}) {

}

func Delete(url string, params ...interface{}) {

}

func Do(method, url string, params ...interface{}) {

	_ := httptest.NewRequest(method, url, nil)
}
