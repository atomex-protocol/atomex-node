package exchange

import (
	"fmt"
	"time"
)

// ErrToManyRequests -
type ErrToManyRequests struct {
	RetryAfter time.Time
}

// Error -
func (e ErrToManyRequests) Error() string {
	return fmt.Sprintf("too many requests. retry after: %s", e.RetryAfter.String())
}
