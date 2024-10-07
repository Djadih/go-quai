package utils

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/common/hexutil"
	"github.com/dominant-strategies/go-quai/core"
	"github.com/dominant-strategies/go-quai/core/rawdb"
	"github.com/dominant-strategies/go-quai/core/state"
	"github.com/dominant-strategies/go-quai/crypto"
	"github.com/dominant-strategies/go-quai/ethdb"
	"github.com/dominant-strategies/go-quai/log"
	"github.com/dominant-strategies/go-quai/params"
)

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database contains incompatible genesis (have %x, new %x)", e.Stored, e.New)
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//	                     genesis == nil       genesis != nil
//	                  +------------------------------------------
//	db has no genesis |  main-net default  |  genesis
//	db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db ethdb.Database, genesis *core.Genesis, nodeLocation common.Location, logger *log.Logger) (*params.ChainConfig, common.Hash, error) {
	return SetupGenesisBlockWithOverride(db, genesis, nodeLocation, 0, logger)
}

func SetupGenesisBlockWithOverride(db ethdb.Database, genesis *core.Genesis, nodeLocation common.Location, startingExpansionNumber uint64, logger *log.Logger) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.AllProgpowProtocolChanges, common.Hash{}, errGenesisNoConfig
	}
	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			logger.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			logger.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db, nodeLocation, startingExpansionNumber)
		if err != nil {
			return genesis.Config, common.Hash{}, err
		}
		return genesis.Config, block.Hash(), nil
	}
	// We have the genesis block in database(perhaps in ancient database)
	// but the corresponding state is missing.
	header := rawdb.ReadHeader(db, 0, stored)
	if _, err := state.New(header.EVMRoot(), header.EtxSetRoot(), header.QuaiStateSize(), state.NewDatabaseWithConfig(db, nil), state.NewDatabaseWithConfig(db, nil), nil, nodeLocation, logger); err != nil {
		if genesis == nil {
			genesis = DefaultGenesisBlock()
		}
		// Ensure the stored genesis matches with the given one.
		hash := genesis.ToBlock(startingExpansionNumber).Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
		block, err := genesis.Commit(db, nodeLocation, startingExpansionNumber)
		if err != nil {
			return genesis.Config, hash, err
		}
		return genesis.Config, block.Hash(), nil
	}
	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := genesis.ToBlock(startingExpansionNumber).Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}
	// Get the existing chain configuration.
	newcfg := genesis.ConfigOrDefault(stored)
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		logger.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.ProgpowColosseumGenesisHash {
		return storedcfg, stored, nil
	}
	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}

	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

// DefaultGenesisBlock returns the Latest default Genesis block.
// Currently it returns the DefaultColosseumGenesisBlock.
func DefaultGenesisBlock() *core.Genesis {
	return DefaultColosseumGenesisBlock("progpow")
}

// DefaultColosseumGenesisBlock returns the Quai Colosseum testnet genesis block.
func DefaultColosseumGenesisBlock(consensusEngine string) *core.Genesis {
	if consensusEngine == "blake3" {
		return &core.Genesis{
			Config:     params.Blake3PowColosseumChainConfig,
			Nonce:      66,
			ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fb"),
			GasLimit:   5000000,
			Difficulty: big.NewInt(2000000),
		}
	}
	return &core.Genesis{
		Config:     params.ProgpowColosseumChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fb"),
		GasLimit:   5000000,
		Difficulty: big.NewInt(1000000000),
	}
}

// DefaultGardenGenesisBlock returns the Garden testnet genesis block.
func DefaultGardenGenesisBlock(consensusEngine string) *core.Genesis {
	if consensusEngine == "blake3" {
		return &core.Genesis{
			Config:     params.Blake3PowGardenChainConfig,
			Nonce:      66,
			ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
			GasLimit:   160000000,
			Difficulty: big.NewInt(500000),
		}
	}
	return &core.Genesis{
		Config:     params.ProgpowGardenChainConfig,
		Nonce:      0,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353539"),
		GasLimit:   5000000,
		Difficulty: big.NewInt(300000000),
	}
}

// DefaultOrchardGenesisBlock returns the Orchard testnet genesis block.
func DefaultOrchardGenesisBlock(consensusEngine string) *core.Genesis {
	if consensusEngine == "blake3" {
		return &core.Genesis{
			Config:     params.Blake3PowOrchardChainConfig,
			Nonce:      66,
			ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fc"),
			GasLimit:   5000000,
			Difficulty: big.NewInt(200000),
		}
	}
	return &core.Genesis{
		Config:     params.ProgpowOrchardChainConfig,
		Nonce:      0,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353536"),
		GasLimit:   5000000,
		Difficulty: big.NewInt(30000000000),
	}
}

// DefaultLighthouseGenesisBlock returns the Lighthouse testnet genesis block.
func DefaultLighthouseGenesisBlock(consensusEngine string) *core.Genesis {
	if consensusEngine == "blake3" {
		return &core.Genesis{
			Config:     params.Blake3PowLighthouseChainConfig,
			Nonce:      66,
			ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fb"),
			GasLimit:   40000000,
			Difficulty: big.NewInt(200000),
		}
	}
	return &core.Genesis{
		Config:     params.ProgpowLighthouseChainConfig,
		Nonce:      0,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353537"),
		GasLimit:   5000000,
		Difficulty: big.NewInt(200000),
	}
}

// DefaultLocalGenesisBlock returns the Local testnet genesis block.
func DefaultLocalGenesisBlock(consensusEngine string) *core.Genesis {
	if consensusEngine == "blake3" {
		return &core.Genesis{
			Config:     params.Blake3PowLocalChainConfig,
			Nonce:      66,
			ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fb"),
			GasLimit:   5000000,
			Difficulty: big.NewInt(500000),
		}
	}
	return &core.Genesis{
		Config:     params.ProgpowLocalChainConfig,
		Nonce:      0,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353535"),
		GasLimit:   5000000,
		Difficulty: big.NewInt(1000),
	}
}

// DeveloperGenesisBlock returns the 'quai --dev' genesis block.
func DeveloperGenesisBlock(period uint64, faucet common.Address) *core.Genesis {
	// Override the default period to the user requested one
	config := *params.AllProgpowProtocolChanges
	// Assemble and return the genesis with the precompiles and faucet pre-funded
	return &core.Genesis{
		Config:     &config,
		ExtraData:  append(append(make([]byte, 32), faucet.Bytes()[:]...), make([]byte, crypto.SignatureLength)...),
		GasLimit:   0x47b760,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Difficulty: big.NewInt(1),
	}
}
