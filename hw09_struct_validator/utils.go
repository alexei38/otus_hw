package hw09structvalidator

import (
	"reflect"
	"strings"
)

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func readTag(f reflect.StructField) (map[string]string, bool) {
	val, ok := f.Tag.Lookup("validate")
	if !ok {
		return nil, false
	}
	opts := strings.Split(val, "|")
	if len(opts) == 0 {
		return nil, false
	}
	validators := make(map[string]string)
	for _, tv := range opts {
		vs := strings.SplitN(tv, ":", 2)
		if len(vs) == 2 {
			validators[vs[0]] = vs[1]
		} else {
			// у nested нет значения
			validators[vs[0]] = ""
		}
	}
	return validators, true
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
