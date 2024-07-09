package keeper

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/evmos/ethermint/x/evm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type bridgeEvent struct {
	Src    common.Address
	Dst    string
	Amount *big.Int
}

const (
	bridgeDiffTreshold = 1
)

// HOOK LOG EVENT FROM EVM
// ! STILL BROKEN AFTER CONSUME GAS AND AMOUNT UNMET
// TODO:: SOLUTION1 SET TRESHOLD OF AMOUNT OF CONVERTION
// TODO:: SOLUTION2 IF USER WANT TO CONVERT ALL TOKEN IF AMOUNT OF USER IN INPUT LESSTHAN IT ACTUAL THEN USE ACTUAL AND CONVERT ALL
func (k *Keeper) CosmosHook(ctx sdk.Context, msg core.Message, logs []*types.Log) error {
	evmParams := k.GetParams(ctx)

	switch destinagtion := msg.To().String(); destinagtion {
	case evmParams.ConverterParams.ConverterContract:
		if !evmParams.ConverterParams.Enable {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Native bridge is on disable state")
		}
		err := k.BridgeEVMxCosmos(ctx, msg, logs)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Native bridge is on disable state")
		}
	default:
		return nil
	}

	return nil
}

func (k *Keeper) BridgeEVMxCosmos(ctx sdk.Context, msg core.Message, logs []*types.Log) error {
	evmParams := k.GetParams(ctx)
	eventName := evmParams.ConverterParams.EventName
	eventAbiString := evmParams.ConverterParams.EventAbi

	// ------------------------------------
	// |                                  |
	// |          CORE CONVERTOR          |
	// |                                  |
	// ------------------------------------

	// eventName = evmParams.ConverterParams.EventName
	event_tuple := evmParams.ConverterParams.EventTuple
	signature := fmt.Sprintf("%v(%v)", eventName, event_tuple)
	topicEventName := hexutil.Encode(k.keccak256([]byte(signature)))
	// unwrap logs by using abi
	// get erc20 abi
	eventAbi, err := abi.JSON(strings.NewReader(eventAbiString))
	if err != nil {
		panic(err)
	}

	_bridgeEvent := bridgeEvent{}

	for _, log := range logs {
		for _, topic := range log.Topics {
			if topic == topicEventName && log.ToEthereum().TxHash.String() != zero_hash {
				err = eventAbi.UnpackIntoInterface(&_bridgeEvent, eventName, log.Data)
				if err != nil {
					panic(err)
				}

				// MINT BURN TOKEN
				// Get the signer address
				signer := sdk.AccAddress(msg.From().Bytes()) // from Eth address to cosmos address

				// check that receiver is cosmos address or ethereum address
				receiver, err := sdk.AccAddressFromBech32(_bridgeEvent.Dst)
				if err != nil {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "receiver address is not cosmos address")
				}
				// check if amount is valid
				intAmount := sdk.NewIntFromBigInt(_bridgeEvent.Amount)
				if intAmount.IsZero() {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount of token is prohibit from module")
				}

				// check if balance and input are valid
				if balance := k.bankKeeper.GetBalance(ctx, signer, "asix"); balance.Amount.LT(intAmount) {
					// if current_balance + 1 >= inputAmount then convert all token of the account

					tresshold_balance := balance.Amount.Add(sdk.NewInt(bridgeDiffTreshold))
					if tresshold_balance.LT(intAmount) {
						return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Amount of token is too high than current balance")
					}
					intAmount = balance.Amount
				}

				supply := k.bankKeeper.GetSupply(ctx, "asix")
				if supply.Amount.LT(intAmount) {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount of token is higher than current total supply")
				}

				//send to module
				convertAmount := sdk.NewCoins(sdk.NewCoin("asix", intAmount))
				if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signer, tokenmngrModuleName, convertAmount); err != nil {
					return sdkerrors.Wrap(types.ErrSendCoinsFromAccountToModule, "Amount of token is too high than current balance due"+err.Error())
				}

				if err := k.bankKeeper.BurnCoins(ctx, tokenmngrModuleName, convertAmount); err != nil {
					return sdkerrors.Wrap(types.ErrBurnCoinsFromModuleAccount, err.Error())
				}

				microSix := sdk.NewCoin("usix", intAmount.QuoRaw(1_000_000_000_000))

				// get the module account balance
				tokenmngrModuleAccount := k.accountKeeper.GetModuleAddress(tokenmngrModuleName)
				moduleBalance := k.bankKeeper.GetBalance(ctx, tokenmngrModuleAccount, "usix")

				// check if module account balance is enough to send
				if moduleBalance.Amount.LT(microSix.Amount) {
					return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "module account balance is not enough to send")
				}

				// send to receiver
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, tokenmngrModuleName, receiver, sdk.NewCoins(microSix),
				); err != nil {
					return sdkerrors.Wrap(types.ErrSendCoinsFromAccountToModule, "unable to send msg.Amounts from module to account despite previously minting msg.Amounts to module account:"+err.Error())
				}
			}
		}
	}
	return nil
}
