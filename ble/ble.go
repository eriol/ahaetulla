package ble

import (
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

// Scan for Bluetooth LE devices for the specified time t.
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

// Return a specific device from its address.
// TODO: right now it's unbounded, ad a max timeout.
func FindDeviceByAddress(adapter *bluetooth.Adapter, address string) bluetooth.ScanResult {

	err := adapter.Enable()
	if err != nil {
		log.Error(err)
	}

	ch := make(chan bluetooth.ScanResult, 1)

	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == address {
			_ = adapter.StopScan()
			ch <- result
		}
	})
	if err != nil {
		log.Error(err)
	}

	return <-ch
}

func Send(address, text string) {
	adapter := bluetooth.DefaultAdapter

	err := adapter.Enable()
	if err != nil {
		log.Error(err)
	}

	r := FindDeviceByAddress(adapter, address)

	var device *bluetooth.Device

	device, err = adapter.Connect(r.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Error(err)
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDNordicUART})
	if err != nil {
		log.Error("Failed to discover the Nordic UART Service: ", err)
	}
	service := services[0]
	chars, err := service.DiscoverCharacteristics(
		[]bluetooth.UUID{bluetooth.CharacteristicUUIDUARTRX, bluetooth.CharacteristicUUIDUARTTX})
	if err != nil {
		log.Error("Failed to discover RX and TX characteristics: ", err)
	}
	rx := chars[0]
	// tx := chars[1]

	// FIXME: this fails sometime.
	_, err = rx.WriteWithoutResponse([]byte(text))
	if err != nil {
		log.Error("Could not send: ", err)
	}

	_ = device.Disconnect()
}
