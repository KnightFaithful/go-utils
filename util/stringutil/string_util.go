package stringutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Object2String(object interface{}) string {
	returnByte, _ := json.Marshal(object)
	return string(returnByte)
}

func Println(list ...string) {
	for _, str := range list {
		fmt.Println(str)
	}
}

func JoinIgnoreEmpty(in []string, sep string) string {
	inIgnoreEmpty := make([]string, 0)
	for _, s := range in {
		if s == "" {
			continue
		}
		inIgnoreEmpty = append(inIgnoreEmpty, s)
	}
	return strings.Join(inIgnoreEmpty, sep)
}

func JoinIgnoreEmptyWith(in []string, sep string, surround string) string {
	inIgnoreEmpty := make([]string, 0)
	for _, s := range in {
		if s == "" {
			continue
		}
		inIgnoreEmpty = append(inIgnoreEmpty, surround+s+surround)
	}
	return strings.Join(inIgnoreEmpty, sep)
}
