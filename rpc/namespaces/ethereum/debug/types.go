package debug

import (
	"github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type TraceCallConfig struct {
	evmtypes.TraceConfig
	StateOverrides *types.StateOverride
	BlockOverrides *types.BlockOverrides
}