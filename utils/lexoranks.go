package utils

import (
	"fmt"
	"strings"
)

func mid_char(c1 byte, c2 byte) string {
	return string(rune((int(c1) + int(c2)) / 2))
}

func SortID(behind string, front string) string {
	c1 := behind[len(behind)-1]
	c2 := front[len(front)-1]

	var result string
	if int(c2)-int(c1) <= 1 {
		result = strings.Clone(behind)
		for i := len(behind); i < len(front); i++ {
			result += "a"
		}
		result += mid_char('a', 'z')
	} else {
		result = front[0:len(front)-1] + mid_char(c1, c2)
	}

	fmt.Println(result)
	return result
}
