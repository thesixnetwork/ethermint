package utils_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/evmos/ethermint/utils"
)

func TestHardFail(t *testing.T) {
	hardFailer := func() {
		panic(utils.DecorateHardFailError(errors.New("some error")))
	}
	panicHandlingFn := func() {
		defer utils.PanicHandler(func(_ any) {})()
		hardFailer()
	}
	require.Panics(t, panicHandlingFn)
}
