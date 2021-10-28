package curlbuilder

import "strings"

func escape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}
