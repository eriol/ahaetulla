package cli

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"noa.mornie.org/eriol/ahaetulla/ble"
)

var scanTime float32

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan to find BLE devices.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		devices := ble.ScanUntilTimeout(scanTime)
		for _, device := range devices {
			fmt.Printf("%s (RSSI=%d): %s\n", device.Address.String(), device.RSSI, device.LocalName())
		}
	},
}

func init() {
	scanCmd.Flags().Float32Var(&scanTime, "scan-time", 5., "scan duration")
	rootCmd.AddCommand(scanCmd)
	log.SetLevel(log.ErrorLevel)
}
