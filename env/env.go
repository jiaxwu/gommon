package env

import (
	"errors"
	"os"
	"strconv"

	"golang.org/x/exp/constraints"
)

// 环境变量不存在
var ErrNotPresent = errors.New("env not present")

// 获取环境变量，可以设置解析和校验函数
func GetEnv[T any](name string, parseAndValidate func(s, name string) (T, error)) (T, error) {
	return getEnv(name, parseAndValidate)
}

// 获取环境变量，字符串
func GetEnvString(name string) (string, error) {
	return GetEnv(name, func(s, name string) (string, error) {
		return s, nil
	})
}

// 获取环境变量，数值
func GetEnvNumber[T constraints.Signed | constraints.Unsigned](name string) (T, error) {
	return GetEnv(name, func(s, name string) (T, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return T(i), nil
	})
}

// 获取环境变量，不能不存在或为空
func MustGetEnv[T any](name string, parseAndValidate func(s, name string) (T, error)) T {
	t, err := GetEnv(name, parseAndValidate)
	if err != nil {
		panic(err)
	}
	return t
}

// 获取环境变量，不能不存在或为空，字符串
func MustGetEnvString(name string) string {
	s, err := GetEnvString(name)
	if err != nil {
		panic(err)
	}
	return s
}

// 获取环境变量，不能不存在或为空，数值
func MustGetEnvNumber[T constraints.Signed | constraints.Unsigned](name string) T {
	n, err := GetEnvNumber[T](name)
	if err != nil {
		panic(err)
	}
	return n
}

func getEnv[T any](name string, parseAndValidate func(s, name string) (T, error)) (T, error) {
	s, ok := os.LookupEnv(name)
	if !ok {
		var t T
		return t, ErrNotPresent
	}
	return parseAndValidate(s, name)
}
