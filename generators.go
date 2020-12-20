package requestid

import (
	"crypto/sha256"
	"hash"
	"strings"
)

// RequestIDGeneratorConfig RequestIDGeneratorConfigを生成するための設定
type RequestIDGeneratorConfig struct {
	MethodRestrict *[]string
	PathRestrict   *[]WildCard
	PathExcept     *[]WildCard
	HeaderEnabled  bool
	HeaderAccept   *[]WildCard
	HeaderDrop     *[]WildCard
	QueryEnabled   bool
	QueryAccept    *[]WildCard
	QueryDrop      *[]WildCard
	CookieEnabled  bool
	CookieAccept   *[]WildCard
	CookieDrop     *[]WildCard
	NamedCleaner   *[]string
	HashFunc       func() hash.Hash
}

// NewGenerator *RequestIDGenerator を取得する
func (config *RequestIDGeneratorConfig) NewGenerator() *RequestIDGenerator {
	gen := &RequestIDGenerator{
		HashFunc: sha256.New,
	}
	if config.MethodRestrict != nil {
		gen.MethodRestrict(*config.MethodRestrict...)
	}
	if config.PathRestrict != nil {
		gen.PathRestrict(*config.PathRestrict...)
	}
	if config.NewGenerator().PathRestrict() != nil {
		gen.PathExcept(*config.PathExcept...)
	}
	if config.HeaderEnabled {
		if config.HeaderAccept != nil {
			gen.HeaderAccept(*config.HeaderAccept...)
		}
		if config.HeaderDrop != nil {
			gen.HeaderDrop(*config.HeaderDrop...)
		}
	}
	if config.QueryEnabled {
		if config.QueryAccept != nil {
			gen.QueryAccept(*config.QueryAccept...)
		}
		if config.QueryDrop != nil {
			gen.QueryDrop(*config.QueryDrop...)
		}
	}
	if config.CookieEnabled {
		if config.CookieAccept != nil {
			gen.CookieAccept(*config.CookieAccept...)
		}
		if config.QueryDrop != nil {
			gen.CookieDrop(*config.CookieDrop...)
		}
	}
	if config.NamedCleaner != nil {
		for _, v := range *config.NamedCleaner {
			f := GetNamedCleaner(v)
			if f != nil {
				gen.AddRequestCleaner(f)
			}
		}
	}
	return gen
}

// NewDefaultRequestIDGeneratorConfig よく使う RequestIDGeneratorConfig を作成する
func NewDefaultRequestIDGeneratorConfig() *RequestIDGeneratorConfig {
	return &RequestIDGeneratorConfig{
		MethodRestrict: &[]string{"GET", "HEAD", "OPTION"},
		HeaderEnabled:  true,
		HeaderAccept:   &[]WildCard{"Host", "Origin", "Authorization", "Accept-Encoding"},
		QueryEnabled:   true,
		QueryDrop:      &[]WildCard{"utm_*", "gclid", "fbclid"},
		CookieEnabled:  false,
		NamedCleaner:   &[]string{"NormarizeAcceptEncodingGzip"},
	}
}

// NormarizeAcceptEncodingGzip Accept-Encoding を正規化する
func NormarizeAcceptEncodingGzip() func(*Request) *Request {
	return func(r *Request) *Request {
		values := r.Header.Values("Accept-Encoding")
		for _, v := range values {
			if strings.Contains(v, "gzip") {
				r.Header.Set("Accept-Encoding", "gzip")
				return r
			}
		}
		r.Header.Del("Accept-Encoding")
		return r
	}
}

// GetNamedCleaner 名前付きのリクエストクリーナーを取得する
func GetNamedCleaner(name string) func(*Request) *Request {
	switch name {
	case "NormarizeAcceptEncodingGzip":
		return NormarizeAcceptEncodingGzip()
	}
	return nil
}
