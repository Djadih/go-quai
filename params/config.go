// Cojyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
	"math/big"

	"github.com/dominant-strategies/go-quai/common"
)

// Genesis hashes to enforce below configs on.
var (
	// Progpow GenesisHashes
	ProgpowColosseumGenesisHash  = common.HexToHash("0x72f90b6f6313e6614dc38a616d3560bd6b7609ae3d14fa9080a680b42df6c4ea")
	ProgpowGardenGenesisHash     = common.HexToHash("0x8e7ffaeb1893e289ec046bd93a5b5de8b9a50a5fa8e2b30775555c5a6ec18a5d")
	ProgpowOrchardGenesisHash    = common.HexToHash("0x8e7ffaeb1893e289ec046bd93a5b5de8b9a50a5fa8e2b30775555c5a6ec18a5d")
	ProgpowLocalGenesisHash      = common.HexToHash("0x817d1fa3848adf89c28ebf70ae52f21e5c1c6a61057660582765277ca30d2243")
	ProgpowLighthouseGenesisHash = common.HexToHash("0xc48a4fda62693510a5b6520a95b8fd8ec82cec550dce1bed8d46b17e1a75546a")

	// Blake3GenesisHashes
	Blake3PowColosseumGenesisHash  = common.HexToHash("0x301114d60dd28adcc51872021d16591e3ed300cbcb0a72f363aeedeef09c2cbb")
	Blake3PowGardenGenesisHash     = common.HexToHash("0xcc765ce0d79736950aeded81e32bbd55a03c4dcec236c1944da32609a237fadc")
	Blake3PowOrchardGenesisHash    = common.HexToHash("0x3b9f95bf2e01aa5a9c09b9b2c102e262a77eebf38bd1c36fdab39f97285233be")
	Blake3PowLocalGenesisHash      = common.HexToHash("0xecbfee9f6bce9583520d16795bc9850783a8a86563c0ce9a5082d86f497b728f")
	Blake3PowLighthouseGenesisHash = common.HexToHash("0x02874a9e708fd75a6ed73bf1d6dbc275a87b8914c3f57da0bbdc959c381d0f67")
)

// Different Network names
const (
	ColosseumName  = "colosseum"
	GardenName     = "garden"
	OrchardName    = "orchard"
	LighthouseName = "lighthouse"
	LocalName      = "local"
	DevName        = "dev"
)

var (
	// ColosseumChainConfig is the chain parameters to run a node on the Colosseum network.
	ProgpowColosseumChainConfig = &ChainConfig{
		ChainID: big.NewInt(9000),
		Progpow: new(ProgpowConfig),
	}

	Blake3PowColosseumChainConfig = &ChainConfig{
		ChainID:   big.NewInt(9000),
		Blake3Pow: new(Blake3powConfig),
	}

	// GardenChainConfig contains the chain parameters to run a node on the Garden test network.
	ProgpowGardenChainConfig = &ChainConfig{
		ChainID: big.NewInt(12000),
		Progpow: new(ProgpowConfig),
	}

	Blake3PowGardenChainConfig = &ChainConfig{
		ChainID:   big.NewInt(12000),
		Blake3Pow: new(Blake3powConfig),
	}

	// OrchardChainConfig contains the chain parameters to run a node on the Orchard test network.
	ProgpowOrchardChainConfig = &ChainConfig{
		ChainID: big.NewInt(15000),
		Progpow: new(ProgpowConfig),
	}

	Blake3PowOrchardChainConfig = &ChainConfig{
		ChainID:   big.NewInt(15000),
		Blake3Pow: new(Blake3powConfig),
	}

	// LighthouseChainConfig contains the chain parameters to run a node on the Lighthouse test network.
	ProgpowLighthouseChainConfig = &ChainConfig{
		ChainID:   big.NewInt(17000),
		Blake3Pow: new(Blake3powConfig),
		Progpow:   new(ProgpowConfig),
	}

	Blake3PowLighthouseChainConfig = &ChainConfig{
		ChainID:   big.NewInt(17000),
		Blake3Pow: new(Blake3powConfig),
	}

	// LocalChainConfig contains the chain parameters to run a node on the Local test network.
	ProgpowLocalChainConfig = &ChainConfig{
		ChainID: big.NewInt(1337),
		Progpow: new(ProgpowConfig),
	}

	Blake3PowLocalChainConfig = &ChainConfig{
		ChainID:   big.NewInt(1337),
		Blake3Pow: new(Blake3powConfig),
	}

	// AllProgpowProtocolChanges contains every protocol change introduced
	// and accepted by the Quai core developers into the Progpow consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllProgpowProtocolChanges = &ChainConfig{big.NewInt(1337), "progpow", new(Blake3powConfig), new(ProgpowConfig), common.Location{}, common.Hash{}, false}

	TestChainConfig = &ChainConfig{big.NewInt(1), "progpow", new(Blake3powConfig), new(ProgpowConfig), common.Location{}, common.Hash{}, false}
	TestRules       = TestChainConfig.Rules(new(big.Int))
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection
	// Various consensus engines
	ConsensusEngine    string
	Blake3Pow          *Blake3powConfig `json:"blake3pow,omitempty"`
	Progpow            *ProgpowConfig   `json:"progpow,omitempty"`
	Location           common.Location
	DefaultGenesisHash common.Hash
	IndexAddressUtxos  bool
}

// SetLocation sets the location on the chain config
func (cfg *ChainConfig) SetLocation(location common.Location) {
	cfg.Location = location
}

// Blake3powConfig is the consensus engine configs for proof-of-work based sealing.
type Blake3powConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *Blake3powConfig) String() string {
	return "blake3pow"
}

// ProgpowConfig is the consensus engine configs for proof-of-work based sealing.
type ProgpowConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *ProgpowConfig) String() string {
	return "progpow"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Blake3Pow != nil:
		engine = c.Blake3Pow
	case c.Progpow != nil:
		engine = c.Progpow
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v, Engine: %v, Location: %v}",
		c.ChainID,
		engine,
		c.Location,
	)
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID *big.Int
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID: new(big.Int).Set(chainID),
	}
}
