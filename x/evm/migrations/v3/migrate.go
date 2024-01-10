package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// "github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore sets the default for GrayGlacierBlock and MergeNetsplitBlock in ChainConfig parameter.
func MigrateStore(ctx sdk.Context, paramstore *paramtypes.Subspace) error {
	// doing nothing
	return nil
}
