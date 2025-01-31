package genallocs

import (
	"math/big"
	"strings"
	"testing"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/params"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	genAllocsStr = `
	[
		{
			"Unlock Schedule": 0,
			"Address": "0x0000000000000000000000000000000000000001",
			"Total": 1029384756000000000000000000,
			"Vested": 1029384756000000000000000000
		},
		{
			"Unlock Schedule": 1,
			"Address": "0x0000000000000000000000000000000000000001",
			"Total": 500000000000000000000000,
			"Vested": 500000000000000000000000
		},
		{
			"Unlock Schedule": 2,
			"Address": "0x0000000000000000000000000000000000000002",
			"Total": 7000000000000000000000000,
			"Vested": 7000000000000000000000000
		},
		{
			"Unlock Schedule": 3,
			"Address": "0x0000000000000000000000000000000000000003",
			"Total": 1000000000000000000000000,
			"Vested": 700000000000000000000000
		}
	]`
)

var (
	expectedAllocs = [4]*GenesisAccount{
		{
			UnlockSchedule:  0,
			Address:         common.HexToAddress("0x0000000000000000000000000000000000000001", common.Location{0, 0}),
			TotalBalance:    new(big.Int).Mul(big.NewInt(1029384756), common.Big10e18),
			VestedBalance:   new(big.Int).Mul(big.NewInt(1029384756), common.Big10e18),
			BalanceSchedule: &orderedmap.OrderedMap[uint64, *big.Int]{},
		},

		{
			UnlockSchedule:  1,
			Address:         common.HexToAddress("0x0000000000000000000000000000000000000001", common.Location{0, 0}),
			TotalBalance:    new(big.Int).Mul(big.NewInt(500000), common.Big10e18),
			VestedBalance:   new(big.Int).Mul(big.NewInt(500000), common.Big10e18),
			BalanceSchedule: &orderedmap.OrderedMap[uint64, *big.Int]{},
		},

		{
			UnlockSchedule:  2,
			Address:         common.HexToAddress("0x0000000000000000000000000000000000000002", common.Location{0, 0}),
			TotalBalance:    new(big.Int).Mul(big.NewInt(7000000), common.Big10e18),
			VestedBalance:   new(big.Int).Mul(big.NewInt(7000000), common.Big10e18),
			BalanceSchedule: &orderedmap.OrderedMap[uint64, *big.Int]{},
		},

		{
			UnlockSchedule:  3,
			Address:         common.HexToAddress("0x0000000000000000000000000000000000000003", common.Location{0, 0}),
			TotalBalance:    new(big.Int).Mul(big.NewInt(1000000), common.Big10e18),
			VestedBalance:   new(big.Int).Mul(big.NewInt(700000), common.Big10e18),
			BalanceSchedule: &orderedmap.OrderedMap[uint64, *big.Int]{},
		},
	}
)

func TestReadingGenallocs(t *testing.T) {

	allocs, err := decodeGenesisAllocs(strings.NewReader(genAllocsStr))
	require.NoError(t, err, "Unable to parse genesis file")

	for num, actualAlloc := range allocs {
		expectedAlloc := expectedAllocs[num]

		require.Equal(t, expectedAlloc.UnlockSchedule, actualAlloc.UnlockSchedule)
		require.Equal(t, expectedAlloc.Address, actualAlloc.Address)
		require.Equal(t, expectedAlloc.TotalBalance, actualAlloc.TotalBalance)
		require.Equal(t, expectedAlloc.VestedBalance, actualAlloc.VestedBalance)
	}
}

type expectedAllocValues struct {
	account       *GenesisAccount
	cliffAmount   *big.Int
	cliffBlock    uint64
	numUnlocks    uint64
	quaiPerUnlock *big.Int
}

func calcExpectedValues(account *GenesisAccount) expectedAllocValues {
	unlockSchedule := unlockSchedules[account.UnlockSchedule]

	cliffPercentage := new(big.Int).SetUint64(unlockSchedule.lumpSumPercentage)
	cliffAmount := cliffPercentage.Mul(cliffPercentage, account.TotalBalance)
	cliffAmount.Div(cliffAmount, common.Big100)

	if cliffAmount.Cmp(account.VestedBalance) > 0 {
		cliffAmount = account.VestedBalance
	}

	cliffIndex := unlockSchedule.lumpSumMonth * params.BlocksPerMonth

	var numUnlocks uint64 = 1 // Cliff index
	var quaiPerUnlock *big.Int
	if unlockSchedule.unlockDuration > 0 {
		// Number of unlocks is calculated by dividing the total allocated over duration then
		// calculating how many unlocks are required to reach (vested balance-cliff).

		unlockableBalance := new(big.Int).Sub(account.TotalBalance, cliffAmount)
		quaiPerUnlock = new(big.Int).Div(unlockableBalance, new(big.Int).SetUint64(unlockSchedule.unlockDuration))

		remainingBalance := new(big.Int).Sub(account.VestedBalance, cliffAmount)
		numUnlocks += new(big.Int).Div(remainingBalance, quaiPerUnlock).Uint64()
	}

	return expectedAllocValues{
		account,
		cliffAmount,
		cliffIndex,
		numUnlocks,
		quaiPerUnlock,
	}
}

func TestCalculatingGenallocs(t *testing.T) {
	allocs, err := decodeGenesisAllocs(strings.NewReader(genAllocsStr))
	require.NoError(t, err, "Unable to parse genesis file")

	for allocNum, actualAlloc := range allocs {
		expectedAlloc := expectedAllocs[allocNum]
		expectedAllocValues := calcExpectedValues(expectedAlloc)
		actualAlloc.calculateLockedBalances()

		// Retrieve UnlockSchedule information.
		unlockSchedule := unlockSchedules[expectedAlloc.UnlockSchedule]

		// Verify the correct number of unlocks are allocated.
		require.Equal(t, expectedAllocValues.numUnlocks, uint64(actualAlloc.BalanceSchedule.Len()))

		totalAdded := new(big.Int)
		// for blockNum, actualUnlock := range actualAlloc.BalanceSchedule.Oldest() {
		for actualUnlock := actualAlloc.BalanceSchedule.Oldest(); actualUnlock != nil; actualUnlock = actualUnlock.Next() {
			blockNum := actualUnlock.Key
			// There should never be unlocks after month 0, before month 12.
			if unlockSchedule.lumpSumMonth != 0 {
				// If lumpSum is not TGE, there should not be an unlock at TGE.
				require.NotEqual(t, 0, blockNum)
			}

			// Verify cliff amount and timing is as expected.
			if blockNum == expectedAllocValues.cliffBlock {
				require.Zero(t, expectedAllocValues.cliffAmount.Cmp(actualUnlock.Value))
			} else if blockNum != actualAlloc.BalanceSchedule.Newest().Key {
				// Verify that all unlocks except for the last one (due to rounding) are correct.
				require.Zero(t, expectedAllocValues.quaiPerUnlock.Cmp(actualUnlock.Value))
			}

			totalAdded.Add(totalAdded, actualUnlock.Value)
		}

		// Verify the total added is equal to the total allocated.
		require.Zero(t, expectedAllocValues.account.VestedBalance.Cmp(totalAdded))
		require.Zero(t, actualAlloc.VestedBalance.Cmp(totalAdded))
	}
}
