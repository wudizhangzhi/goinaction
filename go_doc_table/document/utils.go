package document

import (
	"strings"
)

// 处理选好的内容
func parseParagraphVal(p string) string {
	start := strings.IndexAny(p, sympolCheckedOneline)
	val := ""
	if start != -1 {
		for i, c := range p {
			if i > start {
				if strings.ContainsAny(sympolUncheckOneline, string(c)) ||
					strings.ContainsAny(sympolCheckedOneline, string(c)) {
					break
				}
				val += string(c)
			}
		}
		return val
	}
	return p
}

func contains(s []string, t string) bool {
	for _, i := range s {
		if i == t {
			return true
		}
	}
	return false
}
