package goeland

import "strings"

func SplitAndTrimString(str string) []string {
	list := strings.Split(str, ",")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}
	return list
}
