package strategies

import (
	"modellingSystems/devicefailureExperiment/models"
	"modellingSystems/devicefailureExperiment/models/device"
)

type Strategy func(devices []device.Device, factory models.ArrayFactory[device.Device]) (wasReplaced []bool)

var (
	// ReplaceOnlyFailed replaces only failed devices
	ReplaceOnlyFailed Strategy = func(devices []device.Device, factory models.ArrayFactory[device.Device]) (wasReplaced []bool) {
		for i := range devices {
			if devices[i].WasFailed {
				devices[i] = factory.NewAt(i)
				wasReplaced = append(wasReplaced, true)
			} else {
				wasReplaced = append(wasReplaced, false)
			}
		}
		return
	}

	// ReplaceFailedAndTheOldest replaces failed devices and the oldest device after the one that failed
	ReplaceFailedAndTheOldest Strategy = func(devices []device.Device, factory models.ArrayFactory[device.Device]) (wasReplaced []bool) {
		// Replace failed devices
		for i := range devices {
			if devices[i].WasFailed {
				devices[i] = factory.NewAt(i)
				wasReplaced = append(wasReplaced, true)
			} else {
				wasReplaced = append(wasReplaced, false)
			}
		}
		// Find the oldest device
		oldestDeviceIndex := 0
		for i, device := range devices {
			if device.WorkingTime > devices[oldestDeviceIndex].WorkingTime {
				oldestDeviceIndex = i
			}
		}
		// Replace the oldest device
		devices[oldestDeviceIndex] = factory.NewAt(oldestDeviceIndex)
		return
	}
)
