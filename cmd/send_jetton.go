package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendJettonCmd = &cobra.Command{
	Use:   "jetton <jetton_master_address> <to_address>",
	Short: "Send Jetton tokens",
	Long:  `Transfer Jetton tokens to any TON address.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Jetton transfers: coming soon")
		return nil
	},
}

func init() {
	sendCmd.AddCommand(sendJettonCmd)
}
