// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"strconv"
	"strings"

	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	nodeTypes "github.com/celenium-io/astria-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

func (module *Module) parseConstants(appState nodeTypes.AppState, consensus pkgTypes.ConsensusParams, data *parsedData) {
	// block
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlock,
		Name:   "block_max_bytes",
		Value:  strconv.FormatInt(consensus.Block.MaxBytes, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameBlock,
		Name:   "block_max_gas",
		Value:  strconv.FormatInt(consensus.Block.MaxGas, 10),
	})

	// evidence
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_age_num_blocks",
		Value:  strconv.FormatInt(consensus.Evidence.MaxAgeNumBlocks, 10),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_age_duration",
		Value:  consensus.Evidence.MaxAgeDuration.String(),
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameEvidence,
		Name:   "max_bytes",
		Value:  strconv.FormatInt(consensus.Evidence.MaxBytes, 10),
	})

	// validator
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameValidator,
		Name:   "pub_key_types",
		Value:  strings.Join(consensus.Validator.PubKeyTypes, ","),
	})

	// version
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameVersion,
		Name:   "app",
		Value:  strconv.FormatUint(consensus.Version.AppVersion, 10),
	})

	// generic
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "authority_sudo_key",
		Value:  appState.AuthoritySudoKey,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "native_asset_base_denomination",
		Value:  appState.NativeAssetBaseDenomination,
	})
	data.constants = append(data.constants, storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   "ibc_sudo_address",
		Value:  appState.IbcSudoAddress,
	})
}
