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
	ProgpowColosseumGenesisHash  = common.HexToHash("0x079e89a92bbd044ee6d981d0644c9da1dc960e72d90df212c77e8b9cd12f9ecd")
	ProgpowGardenGenesisHash     = common.HexToHash("0xa4f23733a64776e5d6cd4a37e6ffdfd3ac6b078ec28d03fd8120d48a008b7b4b")
	ProgpowOrchardGenesisHash    = common.HexToHash("0x7eb19c56d86851b39ca4eeb93f5f0b40dee8735f041380c95c88da2ce2bc8899")
	ProgpowLighthouseGenesisHash = common.HexToHash("0x7eb19c56d86851b39ca4eeb93f5f0b40dee8735f041380c95c88da2ce2bc8899")
	ProgpowLocalGenesisHash      = common.HexToHash("0xa6f04684d8de37101324a9f42002424a05e85c610c1eaae5eb7abe4d2a047c48")

	// Blake3GenesisHashes
	Blake3PowColosseumGenesisHash  = common.HexToHash("0xee89c5f775980556c2bf8dc988e477f16005da57f84f2a1000dbe4693ad5f53d")
	Blake3PowGardenGenesisHash     = common.HexToHash("0x46dfa6724020658b3e03cfdf894bd8e990f46d1e9b1b66b49598788e43eaf50f")
	Blake3PowOrchardGenesisHash    = common.HexToHash("0xf60de17f1ae6cbae820d14599eb95581f5c18799f84904520c264be9cfff64c4")
	Blake3PowLighthouseGenesisHash = common.HexToHash("0xf60de17f1ae6cbae820d14599eb95581f5c18799f84904520c264be9cfff64c4")
	Blake3PowLocalGenesisHash      = common.HexToHash("0x36872ac4bc745d9bb86ace9859fd4d12f97b92c6ae32535eb122274ccf9df29c")
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
