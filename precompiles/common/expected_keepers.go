package common

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	// "github.com/ethereum/go-ethereum/common"
)

type BankKeeper interface {
	SendCoins(sdk.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error
	GetBalance(sdk.Context, sdk.AccAddress, string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}

type EVMKeeper interface {
	// GetCodeHash(sdk.Context, common.Address) common.Hash
	// GetBaseDenom(ctx sdk.Context) string
}

type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	// GetLegacyBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) feemarkettypes.Params
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}
