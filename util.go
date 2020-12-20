package requestid

import "net/url"

// CloneValues url.Values を複製する
func CloneValues(v url.Values) url.Values {
	v2 := make(map[string][]string, len(v))
	for k, v := range v {
		v2[k] = v
	}
	return v2
}

// StringSlice []string でContainsするのに使う
type StringSlice []string

// Contains ss に s が含まれるか調べる
func (ss StringSlice) Contains(s string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}
