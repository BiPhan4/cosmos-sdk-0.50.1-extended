package cli

import (
	"bufio"
	"fmt"

	address "cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
)

const (
	flagmoduleName = "module-name"
)

// AddStubCmd returns add-stub cobra Command.
// This command is just for fleshing out how the genesis file can be edited
// The name 'AddStub' is not very descriptive right now.
func AddStubCmd(defaultNodeHome string, addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-stub [address_or_key_name]",
		Short: "return the distribution genesis state",
		Long:  "NA",

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			var kr keyring.Keyring
			addr, err := addressCodec.StringToBytes(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)

				if keyringBackend != "" && clientCtx.Keyring == nil {
					var err error
					kr, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf, clientCtx.Codec)
					if err != nil {
						return err
					}
				} else {
					kr = clientCtx.Keyring
				}

				k, err := kr.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keyring: %w", err)
				}

				addr, err = k.GetAddress()
				if err != nil {
					return err
				}
			}

			moduleNameStr, _ := cmd.Flags().GetString(flagModuleName)

			return genutil.AddStub(clientCtx.Codec, addr, config.GenesisFile(), moduleNameStr)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flagModuleName, "", "module account name")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd

}
