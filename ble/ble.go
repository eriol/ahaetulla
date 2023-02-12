package ble

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

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

	return devices
}
