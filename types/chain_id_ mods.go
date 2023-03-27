package types

import (
	"fmt"
)

const (
	CHAINID_NUMBER_MAINNET = 6666
	CHAINID_NUMBER_TESTNET = 6665
	CHAINID_LOCALNET       = 666
	CHAINID_EPOCH  = 1
)

// Modify ChainID to be from actual cosmos init chain to ethereum compatible
// Without any modification, the chainID will be "sixnet" then the chainID will be "sixnet_666-1"
func ChainIDJumper(chainID string) string {
	if chainID == "sixnet" {
		chainID = fmt.Sprintf("%s_%d-%d", chainID, CHAINID_NUMBER_MAINNET, CHAINID_EPOCH)
	}else if chainID == "fivenet" {
		chainID = fmt.Sprintf("%s_%d-%d", chainID, CHAINID_NUMBER_TESTNET, CHAINID_EPOCH)
	}else {
		chainID = fmt.Sprintf("%s_%d-%d", chainID, CHAINID_LOCALNET, CHAINID_EPOCH)
	}
	return chainID
}
