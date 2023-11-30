package genutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// AddStub is an experimental function which allows minid to modify its genesis file
func AddStub(
	cdc codec.Codec,
	accAddr sdk.AccAddress,
	appendAcct bool,
	genesisFileURL string,
	moduleName string,
) (error, distributiontypes.GenesisState) {

	appState, appGenesis, err := genutiltypes.GenesisStateFromGenFile(genesisFileURL)
	// don't need 'parent' appGenesis for now
	fmt.Println(appGenesis)
	if err != nil {
		// return error message and the zero valued GenesisState
		return fmt.Errorf("failed to unmarshal genesis state: %w", err), distributiontypes.GenesisState{}
	}

	distGenState := distributiontypes.GetGenesisStateFromAppState(cdc, appState)

	// TODO
	// retrieve the interchainaccounts genesis state and unpack the allowed_messages object--is it a []string or []json.RawMessage?

	return nil, distGenState

}
