package commands

import (
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
)

var scanTime float32

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan to find BLE devices.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		scanUntilTimeout(scanTime)
	},
}

func scanUntilTimeout(t float32) {
	adapter := bluetooth.DefaultAdapter
	timer := time.NewTimer(time.Duration(t) * time.Second)
	timerFired := make(chan bool, 1)

	err := adapter.Enable()
	if err != nil {
		log.Error(err)
	}

	go func() {
		<-timer.C
		timerFired <- true
	}()

	go func() {
		for {
			select {
			case _ = <-timerFired:
				adapter.StopScan()
			default:
			}

		}
	}()

	var devices []bluetooth.ScanResult
	addresses := map[string]bool{}
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		deviceAddress := device.Address.String()
		_, already_discovered := addresses[deviceAddress]
		if !already_discovered {
			devices = append(devices, device)
			addresses[deviceAddress] = true
		}
	})
	if err != nil {
		log.Error(err)
	}

	sort.Slice(devices, func(i, j int) bool { return devices[i].RSSI > devices[j].RSSI })

	for _, device := range devices {
		fmt.Printf("%s (RSSI=%d): %s\n", device.Address.String(), device.RSSI, device.LocalName())
	}

}

func init() {
	scanCmd.Flags().Float32Var(&scanTime, "scan-time", 5., "scan duration")
	rootCmd.AddCommand(scanCmd)
	log.SetLevel(log.ErrorLevel)
}
