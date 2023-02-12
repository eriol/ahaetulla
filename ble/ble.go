package ble

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

// Scan for BLE devices for the specified time t.
//
// Devices are returned in descendind RSSI order.
func ScanUntilTimeout(t float32) []bluetooth.ScanResult {
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
		close(timerFired)
	}()

	go func() {
		for range timerFired {
		} // Block until we consume timerFired channel.
		_ = adapter.StopScan()
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

	return devices
}
