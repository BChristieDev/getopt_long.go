/*
	@file      internal/common/common.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

package common

import (
	"strings"
)

func FindIndex[T any](array []T, callbackFn func(T) bool) int {
	for index, element := range array {
		if callbackFn(element) {
			return index
		}
	}

	return -1
}

func IndexOf(str string, searchElement string, fromIndex int) int {
	index := strings.Index(str[fromIndex:], searchElement)

	if index != -1 {
		index += fromIndex
	}

	return index
}

func CharAt(str string, index int) string {
	if index < len(str) {
		return str[index : index+1]
	}

	return ""
}
