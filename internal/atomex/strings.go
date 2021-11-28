package atomex

import "fmt"

func toQuotedBytes(str string) []byte {
	return []byte(fmt.Sprintf("%q", str))
}
