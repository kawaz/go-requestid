package requestid

import "strings"

// StringMatcher 文字列チェッカ
type StringMatcher func(string) bool

// WildCard 文字列のワイルドカードマッチする
type WildCard string

// Matcher StringMatcher を取得する
func (w WildCard) Matcher() StringMatcher {
	if w == "*" {
		return func(s string) bool {
			return true
		}
	}
	if w[len(w)] == '*' {
		prefix := string(w[:len(w)-1])
		return func(s string) bool {
			return strings.HasPrefix(s, prefix)
		}
	}
	str := string(w)
	return func(s string) bool {
		return s == str
	}
}

// WildCardSlice 複数のワイルドカードを使いやすくする
type WildCardSlice []WildCard

// Matchers ワイルドカードにマッチする StringMatcher を取得する
func (ws WildCardSlice) Matchers() []StringMatcher {
	ms := []StringMatcher{}
	for _, w := range ws {
		ms = append(ms, w.Matcher())
	}
	return ms
}

// AnyMatcher どれかのワイルドカードにマッチする StringMatcher を取得する
func (ws WildCardSlice) AnyMatcher() StringMatcher {
	return func(s string) bool {
		for _, m := range ws.Matchers() {
			if m(s) {
				return true
			}
		}
		return false
	}
}

// AllMatcher 全てのワイルドカードにマッチする StringMatcher を取得する
func (ws WildCardSlice) AllMatcher() StringMatcher {
	return func(s string) bool {
		for _, m := range ws.Matchers() {
			if !m(s) {
				return false
			}
		}
		return true
	}
}

// NotAnyMatcher 全てのワイルドカードにマッチしない StringMatcher を取得する
func (ws WildCardSlice) NotAnyMatcher() StringMatcher {
	return func(s string) bool {
		return ws.AnyMatcher()(s)
	}
}
