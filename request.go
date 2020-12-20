package requestid

import (
	"hash"
	"net/http"
	"net/url"
)

// RequestID 重複判定用のハッシュ値
type RequestID string

// Request http.Request のサブセット的なやつ
type Request struct {
	Expected      bool
	Method        string
	Path          string
	Header        http.Header
	HeaderEnabled bool
	Query         url.Values
	QueryEnabled  bool
	Cookies       url.Values
	CookieEnabled bool
}

// GenerateID generate RequestID
func (r *Request) GenerateID(newHashFunc func() hash.Hash) (RequestID, bool) {
	if r.Expected {
		return RequestID(""), false
	}
	hash := newHashFunc()
	hash.Write([]byte("method:"))
	hash.Write([]byte(r.Method))
	hash.Write([]byte("\n"))
	hash.Write([]byte("path:"))
	hash.Write([]byte(url.PathEscape(r.Path)))
	hash.Write([]byte("\n"))
	hash.Write([]byte("header:"))
	hash.Write([]byte(url.Values(r.Header).Encode()))
	hash.Write([]byte("\n"))
	hash.Write([]byte("query:"))
	hash.Write([]byte(r.Query.Encode()))
	hash.Write([]byte("\n"))
	hash.Write([]byte("cookie:"))
	hash.Write([]byte(r.Cookies.Encode()))
	hash.Write([]byte("\n"))
	return RequestID(hash.Sum(nil)), true
}

// Clone Request を複製する
func (r *Request) Clone() *Request {
	r2 := &Request{
		Expected:      r.Expected,
		Method:        r.Method,
		Path:          r.Path,
		Header:        r.Header.Clone(),
		HeaderEnabled: r.HeaderEnabled,
		Query:         CloneValues(r.Query),
		QueryEnabled:  r.QueryEnabled,
		Cookies:       CloneValues(r.Cookies),
		CookieEnabled: r.CookieEnabled,
	}
	return r2
}
