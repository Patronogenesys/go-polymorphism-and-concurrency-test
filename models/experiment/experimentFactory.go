package experiment

import (
	"modellingSystems/devicefailureExperiment/config"
	"modellingSystems/devicefailureExperiment/models/clock"
	"modellingSystems/devicefailureExperiment/models/device"
	"modellingSystems/devicefailureExperiment/models/experiment/strategies"
)

// TODO: add a device factory to the experiment init []devices
type DefaultFactory struct {
	// Constants (Inputs)
	ReplacementStrategy strategies.Strategy
	ResultChan          chan Result
}

func (f DefaultFactory) New() Interface {
	factory := device.DefaultArrayFactory{}
	// cerate DeviceNumberDefault devices
	devices := factory.NewArray(config.DeviceNumberDefault)
	return &Experiment{
		DeviceArrayFactory:  factory,
		Devices:             devices,
		Clock:               clock.DefaultExperimentClockFactory{}.New(),
		AwaitPrice:          config.AwaitPriceDefault,
		ReplaceDuration:     config.ReplaceDurationDefault,
		ReplacementStrategy: f.ReplacementStrategy,
		ResultChan:          f.ResultChan,
	}
}

//type ParameterisedFactory struct {
//	// Constants (Inputs)
//	DeviceFactory       models.Factory[device.Device]
//	NumberOfDevices     int
//	ExperimentDuration  int
//	ReplacePrice        int
//	AwaitPrice          int
//	ReplaceDuration     int
//	ReplacementStrategy strategies.Strategy
//
//	// Output
//	ResultChan chan Result
//}
//
//func (f ParameterisedFactory) New() Experiment {
//	factory := f.DeviceFactory
//	// cerate NumberOfDevices devices
//	devices := make([]device.Device, f.NumberOfDevices)
//	for i := range devices {
//		devices[i] = factory.New()
//	}
//
//	return Experiment{
//		DeviceArrayFactory:  factory,
//		Devices:             devices,
//		Clock:               clock.ExperimentClockFactory{Duration: f.ExperimentDuration}.New(),
//		AwaitPrice:          f.AwaitPrice,
//		ReplaceDuration:     f.ReplaceDuration,
//		ReplacementStrategy: f.ReplacementStrategy,
//		ResultChan:          f.ResultChan,
//	}
//}
