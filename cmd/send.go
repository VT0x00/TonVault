package cmd

import "github.com/spf13/cobra"

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send tokens",
	Long:  `Send TON coins or Jetton tokens to any address.`,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
