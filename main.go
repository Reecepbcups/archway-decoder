package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/spf13/cobra"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/archway-network/archway/app"
	"github.com/archway-network/archway/app/params"
)

// arch-decode tx decode <tx-b64>

func main() {
	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	cfg.Seal()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		// WithAccountRetriever(authtypes.AccountRetriever{}).
		// WithBroadcastMode(flags.BroadcastBlock).
		// WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: "Archway Network Lightweight Decoder",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	// Keeping it as close to the original as possible
	rootCmd.AddCommand(
		txCommand(),
	)
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetDecodeCommand(),
		GetFileDecodeCommand(),
	)

	return cmd
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Decode struct {
	// ID is the unique ID for the SQL database transaction
	ID int `json:"id"`
	// tx is the base64 amino in the input file, and the Decoded JSON in the output file
	Tx string `json:"tx"`
}

type Decodes []Decode

func decodeTx(clientCtx client.Context, wg *sync.WaitGroup, jobs <-chan Decode, results chan<- Decode) {
	defer wg.Done()
	for value := range jobs {
		txBytes, err := base64.StdEncoding.DecodeString(value.Tx)
		if err != nil {
			panic(err)
		}

		tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
		if err != nil {
			panic(err)
		}

		json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
		if err != nil {
			panic(err)
		}

		results <- Decode{
			ID: value.ID,
			Tx: string(json),
		}
	}
}

// GetDecodeCommand returns the decode command to take serialized bytes and turn
// it into a JSON-encoded transaction.
func GetFileDecodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode-file [file] [output-file-name]",
		Short: "Decode a bunch of amino bytre strings in 1 file. Then export",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			start := time.Now()
			// args := os.Args[1:]
			if len(args) < 2 {
				fmt.Println("Usage: ./arch-decoder tx decode-file input.json output.json")
				return
			}

			dat, err := ioutil.ReadFile(args[0])
			check(err)

			var values Decodes
			err = json.Unmarshal(dat, &values)
			check(err)

			clientCtx := client.GetClientContextFromCmd(cmd)

			jobs := make(chan Decode, len(values))
			results := make(chan Decode, len(values))

			var wg sync.WaitGroup

			cores := runtime.NumCPU()
			wg.Add(cores)
			for i := 0; i < cores; i++ {
				go decodeTx(clientCtx, &wg, jobs, results)
			}

			for _, value := range values {
				jobs <- value
			}
			close(jobs)

			newValues := make([]Decode, 0, len(values))
			for i := 0; i < len(values); i++ {
				newValues = append(newValues, <-results)
			}

			wg.Wait()

			output, err := json.Marshal(newValues)
			check(err)
			err = ioutil.WriteFile(args[1], output, 0644)
			check(err)

			fmt.Println("Decode time taken:", time.Since(start))

			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
