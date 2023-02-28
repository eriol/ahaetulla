package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"noa.mornie.org/eriol/ahaetulla/ble"
)

var address string

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send one or more words over BLE to a specific device.",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		text := strings.Join(args, " ")
		ble.Send(address, text)
	},
}

func init() {
	sendCmd.Flags().StringVar(&address, "device", "", "ble device address")
	_ = sendCmd.MarkFlagRequired("device")
	rootCmd.AddCommand(sendCmd)
}
