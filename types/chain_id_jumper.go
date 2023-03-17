package types

import (
	"fmt"
)

const (
	CHAINID_NUMBER = 666
	CHAINID_EPOCH  = 1
)

// Modify ChainID to be from actual cosmos init chain to ethereum compatible
func ChainIDJumper(chainID string) string {
	chainID = fmt.Sprintf("%s_%d-%d", chainID, CHAINID_NUMBER, CHAINID_EPOCH)
	return chainID
}
