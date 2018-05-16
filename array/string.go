package array

import (
	"fmt"
	"strings"
)

// StringInArray return true if array contain the string
func StringInArray(in string, arr ...string) bool {
	for i := range arr {
		if arr[i] == in {
			return true
		}
	}

	return false
}

// ArrayToString convert array of int to string with a separator
func ArrayToString(a []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")

}
