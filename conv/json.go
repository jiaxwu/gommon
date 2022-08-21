package conv

import (
	"encoding/json"
)

// 必须序列化
func MustMarshal(v any) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes
}

// 必须反序列化
func MustUnmarshal[T any](b []byte) T {
	var t T
	if err := json.Unmarshal(b, &t); err != nil {
		panic(err)
	}
	return t
}

// 必须序列化
func MustMarshalToString(v any) string {
	return string(MustMarshal(v))
}

// 必须反序列化
func MustUnmarshalString[T any](b string) T {
	return MustUnmarshal[T]([]byte(b))
}
