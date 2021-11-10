package tezos

import "fmt"

type blockIDHead struct {
	prev string
}

// ID -
func (b blockIDHead) ID() string {
	if b.prev != "" {
		return fmt.Sprintf("head~%s", b.prev)
	}
	return "head"
}
