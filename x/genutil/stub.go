package genutil

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// AddStub is an experimental function which allows minid to modify its genesis file
func AddStub(
	cdc codec.Codec,
	accAddr sdk.AccAddress,
	appendAcct bool,
	genesisFileURL, amountStr, vestingAmtStr string,
	vestingStart, vestingEnd int64,
	moduleName string,
) error {

	coins, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
	fmt.Println(vestingAmt)
	if err != nil {
		return fmt.Errorf("failed to parse vesting aomunt: %w", err)
	}

	var genAccount authtypes.GenesisAccount

	balances := banktypes.Balance{Address: accAddr.String(), Coins: coins.Sort()}
	baseAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)
	fmt.Println(baseAccount)

	// note: genesis file has to be created from 'init' first
	appState, appGenesis, err := genutiltypes.GenesisStateFromGenFile(genesisFileURL)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}
	// There must be an unmarshalling function for the interchain accounts state somewhere?
	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	if accs.Contains(accAddr) {
		if !appendAcct {
			return fmt.Errorf(" Account %s already exists\nUse 'append' flag to append account at existing address", accAddr)
		}

		genesisB := banktypes.GetGenesisStateFromAppState(cdc, appState)
		for idx, acc := range genesisB.Balances {
			if acc.Address != accAddr.String() {
				continue
			}

			updatedCoins := acc.Coins.Add(coins...)
			bankGenState.Balances[idx] = banktypes.Balance{Address: accAddr.String(), Coins: updatedCoins.Sort()}
			break
		}
	} else {
		// Add the new account to the set of genesis accounts and sanitize the accounts afterwards.
		accs = append(accs, genAccount)
		accs = authtypes.SanitizeGenesisAccounts(accs)

		genAccs, err := authtypes.PackAccounts(accs)
		if err != nil {
			return fmt.Errorf("failed to convert accounts into anys: %w", err)
		}
		authGenState.Accounts = genAccs

		authGenStateBz, err := cdc.MarshalJSON(&authGenState)
		if err != nil {
			return fmt.Errorf("failed to marshal auth genesis state: %w", err)
		}
		appState[authtypes.ModuleName] = authGenStateBz

		bankGenState.Balances = append(bankGenState.Balances, balances)
	}

	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	appGenesis.AppState = appStateJSON
	return ExportGenesisFile(appGenesis, genesisFileURL)

}
