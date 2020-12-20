package requestid

import (
	"hash"
	"net/http"
	"net/textproto"
	"net/url"
)

// RequestIDGenerator generate RequestID from *http.Request
type RequestIDGenerator struct {
	HashFunc        func() hash.Hash
	requestCleaners []func(*Request) *Request
}

// GenerateID generate RequestID
func (gen *RequestIDGenerator) GenerateID(r *http.Request) (RequestID, bool) {
	cleanRequest := gen.cleanupRequest(r)
	return cleanRequest.GenerateID(gen.HashFunc)
}

// CleanupRequest は httpd.Request を整理します
func (gen *RequestIDGenerator) cleanupRequest(r *http.Request) *Request {
	r2 := &Request{
		Expected: false,
		Method:   r.Method,
		Path:     r.URL.Path,
		Query:    r.URL.Query(),
		Header:   r.Header.Clone(),
		Cookies:  url.Values{},
	}
	if r2.Method == "" {
		r2.Method = "GET"
	}
	for _, c := range r.Cookies() {
		r2.Cookies.Add(c.Name, c.Value)
	}
	for _, f := range gen.requestCleaners {
		r2 = f(r2)
		if r2.Expected {
			return r2
		}
	}
	if !r2.HeaderEnabled {
		r2.Header = nil
	}
	if !r2.QueryEnabled {
		r2.Query = nil
	}
	if !r2.CookieEnabled {
		r2.Cookies = nil
	}
	return r2
}

// AddRequestCleaner Request を集約する関数を登録する
func (gen *RequestIDGenerator) AddRequestCleaner(f func(*Request) *Request) *RequestIDGenerator {
	gen.requestCleaners = append(gen.requestCleaners, f)
	return gen
}

// MethodRestrict メソッドを制限する
func (gen *RequestIDGenerator) MethodRestrict(methods ...string) *RequestIDGenerator {
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.Expected = true
		for _, v := range methods {
			if r.Method == v {
				r.Expected = false
				break
			}
		}
		return r
	})
}

// PathRestrict パスを制限する
func (gen *RequestIDGenerator) PathRestrict(paths ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(paths).NotAnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		if m(r.Path) {
			r.Expected = true
		}
		return r
	})
}

// HeaderAccept ヘッダを制限する
func (gen *RequestIDGenerator) HeaderAccept(keys ...WildCard) *RequestIDGenerator {
	for i, v := range keys {
		keys[i] = WildCard(textproto.CanonicalMIMEHeaderKey(string(v)))
	}
	m := WildCardSlice(keys).NotAnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.HeaderEnabled = true
		for k := range r.Header {
			if m(k) {
				r.Header.Del(k)
			}
		}
		return r
	})
}

// HeaderAcceptAll ヘッダを全て利用する
func (gen *RequestIDGenerator) HeaderAcceptAll() *RequestIDGenerator {
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.HeaderEnabled = true
		return r
	})
}

// HeaderDrop ヘッダを一部削除する
func (gen *RequestIDGenerator) HeaderDrop(keys ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(keys).AnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.HeaderEnabled = true
		for k := range r.Header {
			if m(k) {
				r.Header.Del(k)
			}
		}
		return r
	})
}

// QueryAccept クエリパラメータを制限する
func (gen *RequestIDGenerator) QueryAccept(keys ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(keys).NotAnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.QueryEnabled = true
		for k := range r.Query {
			if m(k) {
				r.Query.Del(k)
			}
		}
		return r
	})
}

// QueryAcceptAll クエリパラメータを全て利用する
func (gen *RequestIDGenerator) QueryAcceptAll() *RequestIDGenerator {
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.QueryEnabled = true
		return r
	})
}

// QueryDrop クエリパラメータを一部削除する
func (gen *RequestIDGenerator) QueryDrop(keys ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(keys).AnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.QueryEnabled = true
		for k := range r.Query {
			if m(k) {
				r.Query.Del(k)
			}
		}
		return r
	})
}

// QueryDropTracking トラッキングパラメータを削除する
func (gen *RequestIDGenerator) QueryDropTracking() *RequestIDGenerator {
	return gen.QueryDrop("utm_*", "gclid", "fbclid")
}

// CookieAccept クッキーを制限する
func (gen *RequestIDGenerator) CookieAccept(keys ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(keys).NotAnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.CookieEnabled = true
		for k := range r.Cookies {
			if m(k) {
				r.Cookies.Del(k)
			}
		}
		return r
	})
}

// CookieAcceptAll クッキーを全て利用する
func (gen *RequestIDGenerator) CookieAcceptAll() *RequestIDGenerator {
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.CookieEnabled = true
		return r
	})
}

// CookieDrop クッキーを一部削除する
func (gen *RequestIDGenerator) CookieDrop(keys ...WildCard) *RequestIDGenerator {
	m := WildCardSlice(keys).AnyMatcher()
	return gen.AddRequestCleaner(func(r *Request) *Request {
		r.CookieEnabled = true
		for k := range r.Cookies {
			if m(k) {
				r.Cookies.Del(k)
			}
		}
		return r
	})
}
