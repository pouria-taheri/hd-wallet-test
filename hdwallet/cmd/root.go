package cmd

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"github.com/spf13/cobra"
)

func enableCli() {
	app.Config.Cmd.Client.Cli.Enabled = true
}

var (
	rootCli *cobra.Command
	genTemplateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate config template file.",
		Long: `generates new empty config template file 
named "template.yaml" for filling important information of wallet to prepare 
file for encryption so application will be able to process the specified wallet.`,
		Run: func(cmd *cobra.Command, args []string) {
			msg := domain.NewMessage()
			msg.SetType(domain.MessageTypeEnumGenerate)
			msg.SetArgs(args)

			enableCli()
			initialize()
			app.Handler.Generate(msg)
		},
	}
	lockCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt current config.",
		Long: `encrypts old unsecure config file. the old config file will be 
overridden by new file.`,
		Run: func(cmd *cobra.Command, args []string) {
			msg := domain.NewMessage()
			msg.SetType(domain.MessageTypeEnumEncrypt)
			msg.SetArgs(args)

			enableCli()
			initialize()
			app.Handler.Encrypt(msg)
		},
	}
	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt the encrypted config using suitable password.",
		Long: `decrypts old secured config file. the old config file will be 
overridden by new file.`,
		Run: func(cmd *cobra.Command, args []string) {
			msg := domain.NewMessage()
			msg.SetType(domain.MessageTypeEnumDecrypt)
			msg.SetArgs(args)

			enableCli()
			initialize()
			app.Handler.Decrypt(msg)
		},
	}
	unlockCmd = &cobra.Command{
		Use:   "unlock",
		Short: "unlock the server remotely using grpc",
		Long: `requests to server for unlocking encrypted file so application 
could resume the process.`,
		Run: func(cmd *cobra.Command, args []string) {
			app.Config.Mode = config.AppModeClient
			msg := domain.NewMessage()
			msg.SetType(domain.MessageTypeEnumUnlock)
			msg.SetArgs(args)

			enableCli()
			initialize()
			app.Handler.Unlock(msg)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCli.Execute()
}

func init() {
	rootCli = &cobra.Command{
		Use:   "",
		Short: "HD Wallet generator & blockchain transaction signer",
		Long: `HD Wallet(Hierarchical Deterministic Wallet) is being used by 
blockchain gateway for managing blockchain accounts and signing transactions.`,
		Run: func(cmd *cobra.Command, args []string) {
			initialize()
			run()
		},
	}
	rootCli.AddCommand(genTemplateCmd)
	rootCli.AddCommand(lockCmd)
	rootCli.AddCommand(decryptCmd)
	rootCli.AddCommand(unlockCmd)
}
