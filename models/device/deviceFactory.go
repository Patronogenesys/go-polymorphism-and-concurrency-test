package device

import (
	"modellingSystems/devicefailureExperiment/config"
)

//type DefaultFactory struct{}
//
//func (f DefaultFactory) New() Device {
//	return Device{
//		ReplacementCost:      config.ReplacementCostDefault,
//		ProbabilityOfFailure: config.ProbabilityOfFailureDefault,
//	}
//}
//
//type ParameterisedFactory struct {
//	ReplacementCost      int
//	ProbabilityOfFailure func(time int) float64
//}
//
//func (f ParameterisedFactory) New() Device {
//	return Device{
//		ReplacementCost:      f.ReplacementCost,
//		ProbabilityOfFailure: f.ProbabilityOfFailure,
//	}
//}

type DefaultArrayFactory struct{}

// NewArray creates an array of devices with the same probability of failure
// but with different replacement costs repeating with the pattern of the mask
func (f DefaultArrayFactory) NewArray(size int) []Device {
	devices := make([]Device, size)
	for i := range devices {
		devices[i] = f.NewAt(i)
	}
	return devices
}

func (f DefaultArrayFactory) NewAt(index int) Device {
	ReplacementCostMask := config.ReplacementCostMaskDefault
	return Device{
		ReplacementCost:      ReplacementCostMask[index%len(ReplacementCostMask)],
		ProbabilityOfFailure: config.ProbabilityOfFailureDefault,
	}
}
