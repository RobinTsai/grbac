package utils

import "encoding/json"

func StringifyJson(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		return "[StringifyJson] " + err.Error()
	}
	return string(b)
}

// UniqueStrings 不能保证稳定性
// option 方法用于统一 key
func UniqueStrings(ss []string, allowEmpty bool, option ...func(string) string) []string {
	m := make(map[string]struct{})

	for _, s := range ss {
		for _, o := range option {
			s = o(s)
		}
		m[s] = struct{}{}
	}
	if !allowEmpty {
		delete(m, "")
	}

	newSS := make([]string, 0, len(m))
	for s := range m {
		newSS = append(newSS, s)
	}
	return newSS
}

// Contains ...
func Contains(ss []string, s string) bool {
	for _, s2 := range ss {
		if s == s2 {
			return true
		}
	}
	return false
}

// Remove ...
func Remove(ss []string, s string) []string {
Again:
	for i, s2 := range ss {
		if s == s2 {
			ss = append(ss[0:i], ss[i+1:]...)
			goto Again
		}
	}
	return ss
}
