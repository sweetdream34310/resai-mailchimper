package aways

import "strings"

func RemoveEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		str = strings.TrimSpace(str) // remove space
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}
