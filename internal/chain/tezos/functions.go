package tezos

import "github.com/aopoltorzhicky/watch_tower/internal/chain"

func toOperationStatus(status string) chain.OperationStatus {
	switch status {
	case "applied":
		return chain.Applied
	case "failed", "backtracked", "skipped":
		return chain.Failed
	default:
		return chain.Pending
	}
}
