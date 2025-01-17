package genallocs

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/params"
	"github.com/stretchr/testify/require"
)

const (
	genAllocsStr = `
	[
		{
			"Vest Schedule": 0,
			"Address": "0x0000000000000000000000000000000000000001",
			"Amount": 500000
		},
		{
			"Vest Schedule": 1,
			"Address": "0x0000000000000000000000000000000000000002",
			"Amount": 7000000
		},
		{
			"Vest Schedule": 2,
			"Address": "0x0000000000000000000000000000000000000003",
			"Amount": 1234567
		}
	]`
)

var (
	expectedAllocs = [3]GenesisAccount{
		{
			VestSchedule:    0,
			Address:         common.HexToAddress("0x0000000000000000000000000000000000000001", common.Location{0, 0}),
			TotalBalance:    500000,
			BalanceSchedule: map[uint64]*big.Int{},
		},
		{
			VestSchedule: 1,
			Address:      common.HexToAddress("0x0000000000000000000000000000000000000002", common.Location{0, 0}),
			TotalBalance: 7000000,
			BalanceSchedule: map[uint64]*big.Int{
				0:                                 big.NewInt(7000000 * 25 / 100),
				(12)*params.BlocksPerMonth - 1:    big.NewInt(145833),
				(12+1)*params.BlocksPerMonth - 1:  big.NewInt(145833),
				(12+2)*params.BlocksPerMonth - 1:  big.NewInt(145833),
				(12+3)*params.BlocksPerMonth - 1:  big.NewInt(145833),
				(12+4)*params.BlocksPerMonth - 1:  big.NewInt(145833),
				(12+34)*params.BlocksPerMonth - 1: big.NewInt(145833),
				(12+35)*params.BlocksPerMonth - 1: big.NewInt(145833),
				(12+36)*params.BlocksPerMonth - 1: big.NewInt(145833 + 12), // rounding
			},
		},
		{
			VestSchedule: 2,
			Address:      common.HexToAddress("0x0000000000000000000000000000000000000003", common.Location{0, 0}),
			TotalBalance: 1234567,
			BalanceSchedule: map[uint64]*big.Int{
				0:                                 big.NewInt(0),
				(12)*params.BlocksPerMonth - 1:    big.NewInt(34293),
				(12+1)*params.BlocksPerMonth - 1:  big.NewInt(34293),
				(12+2)*params.BlocksPerMonth - 1:  big.NewInt(34293),
				(12+3)*params.BlocksPerMonth - 1:  big.NewInt(34293),
				(12+4)*params.BlocksPerMonth - 1:  big.NewInt(34293),
				(12+34)*params.BlocksPerMonth - 1: big.NewInt(34293),
				(12+35)*params.BlocksPerMonth - 1: big.NewInt(34293),
				(12+36)*params.BlocksPerMonth - 1: big.NewInt(34293 + 18), // rounding
			},
		},
	}
)

func altCalcBalances(account *GenesisAccount) {
	total := new(big.Int).Mul(big.NewInt(int64(account.TotalBalance)), common.Big10e18)

	vestingSchedule := vestingSchedules[account.VestSchedule]
	tgePercentage := int64(vestingSchedule.tgePercentage * 100)
	tgeAmount := new(big.Int).Div(new(big.Int).Mul(total, big.NewInt(tgePercentage)), big.NewInt(100))
	unlock := new(big.Int).Div(new(big.Int).Sub(total, tgeAmount), new(big.Int).SetUint64(vestingSchedule.vestDuration*12))
	rounded := new(big.Int).Sub(total, new(big.Int).Add(tgeAmount, new(big.Int).Mul(unlock, new(big.Int).SetUint64(vestingSchedule.vestDuration*12))))

	account.BalanceSchedule[0] = tgeAmount
	for unlockMonth := uint64(12); unlockMonth <= vestingSchedule.vestDuration*12; unlockMonth++ {
		account.BalanceSchedule[unlockMonth*params.BlocksPerMonth-1] = unlock
	}
	// account.BalanceSchedule[unlockMonth*params.BlocksPerMonth-1] = account.BalanceSchedule[unlockMonth*params.BlocksPerMonth-1].Add(account.BalanceSchedule[unlockMonth*params.BlocksPerMonth-1], rounded)
	unlockMonth := vestingSchedule.vestDuration*12*params.BlocksPerMonth - 1
	account.BalanceSchedule[unlockMonth].Add(account.BalanceSchedule[unlockMonth], rounded)
}

func TestReadingGenallocs(t *testing.T) {

	allocs, err := decodeGenesisAllocs(strings.NewReader(genAllocsStr))
	require.NoError(t, err, "Unable to parse genesis file")

	for num, actualAlloc := range allocs {
		expectedAlloc := expectedAllocs[num]

		require.Equal(t, expectedAlloc.VestSchedule, actualAlloc.VestSchedule)
		require.Equal(t, expectedAlloc.Address, actualAlloc.Address)
		require.Equal(t, expectedAlloc.TotalBalance, actualAlloc.TotalBalance)
	}
}

func TestCalculatingGenallocs(t *testing.T) {
	allocs, err := decodeGenesisAllocs(strings.NewReader(genAllocsStr))
	require.NoError(t, err, "Unable to parse genesis file")

	for allocNum, actualAlloc := range allocs[:1] {
		actualAlloc.calculateLockedBalances()
		altCalcBalances(&expectedAllocs[allocNum])
		for blockNum, expectedUnlock := range expectedAllocs[allocNum].BalanceSchedule {
			require.Zero(t,
				expectedUnlock.Cmp(actualAlloc.BalanceSchedule[blockNum]),
				fmt.Sprintf("incorrect balance unlock on block %d", blockNum),
			)
		}
	}
}
