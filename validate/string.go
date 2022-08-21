package validate

// 字符长度在范围内
func Len(s string, min, max int, errFunc func(min, max int) error) error {
	if len(s) < min || len(s) > max {
		return errFunc(min, max)
	}
	return nil
}

// 字符串转字符集
func StringToCharset(s string) map[rune]struct{} {
	m := map[rune]struct{}{}
	for _, c := range s {
		m[c] = struct{}{}
	}
	return m
}

// 字符串在字符集里面
func InCharset(s string, charset map[rune]struct{}, errFunc func(invalidChar rune) error) error {
	for _, c := range s {
		if _, ok := charset[c]; !ok {
			return errFunc(c)
		}
	}
	return nil
}

// 字符串包含字符集
func IncludeCharset(s string, charset map[rune]struct{}, errFunc func() error) error {
	if !includeCharset(s, charset) {
		return errFunc()
	}
	return nil
}

func includeCharset(s string, charset map[rune]struct{}) bool {
	for _, c := range s {
		if _, ok := charset[c]; ok {
			return true
		}
	}
	return false
}

// 字符串包含所有字符集
func IncludeCharsets(s string, charsets []map[rune]struct{}, errFunc func(notIncludedCharset map[rune]struct{}) error) error {
	for _, charset := range charsets {
		if !includeCharset(s, charset) {
			return errFunc(charset)
		}
	}
	return nil
}

// 字符串包含字符集数量
func IncludeCharsetsCount(s string, min, max int, charsets []map[rune]struct{}, errFunc func(min, max, count int) error) error {
	cnt := 0
	for _, charset := range charsets {
		if includeCharset(s, charset) {
			cnt++
		}
	}
	if cnt < min || cnt > max {
		return errFunc(min, max, cnt)
	}
	return nil
}
