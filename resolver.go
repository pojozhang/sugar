package rum

import (
	"net/http"
	"strconv"
)

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

type Resolver interface {
	resolve(req *http.Request, params []interface{}, param interface{}, index int)
}

type QueryResolver struct {
}

func (r *QueryResolver) resolve(req *http.Request, params []interface{}, param interface{}, index int) {
	q := param.(Query)
	query := req.URL.Query()
	for k, v := range q {
		query.Add(k, ToString(v))
	}
	req.URL.RawQuery = query.Encode()
}

func ToString(v interface{}) string {
	var s string
	switch x := v.(type) {
	case uint:
		s = strconv.FormatUint(uint64(x), 10)
	case uint8:
		s = strconv.FormatUint(uint64(x), 10)
	case uint16:
		s = strconv.FormatUint(uint64(x), 10)
	case uint32:
		s = strconv.FormatUint(uint64(x), 10)
	case uint64:
		s = strconv.FormatUint(uint64(x), 10)
	case int:
		s = strconv.FormatInt(int64(x), 10)
	case int8:
		s = strconv.FormatInt(int64(x), 10)
	case int16:
		s = strconv.FormatInt(int64(x), 10)
	case int32:
		s = strconv.FormatInt(int64(x), 10)
	case int64:
		s = strconv.FormatInt(int64(x), 10)
	case float32:
		s = strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(float64(x), 'f', -1, 64)
	case string:
		s = v.(string)
	}

	return s
}
