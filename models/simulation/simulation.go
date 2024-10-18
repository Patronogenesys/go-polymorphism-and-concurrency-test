package simulation

import (
	"fmt"
	"modellingSystems/devicefailureExperiment/config"
	"modellingSystems/devicefailureExperiment/models"
	"modellingSystems/devicefailureExperiment/models/experiment"
	"modellingSystems/devicefailureExperiment/models/experiment/strategies"
	"sync"
)

type Interface interface {
	Run() <-chan Result
}

type Result struct {
	Name                string
	Errors              []error
	ExperimentCount     int
	avgReplacingExpense float64
	avgAwaitExpense     float64
	avgTotalExpense     float64
	avgFailedDevices    float64
}

func (r Result) String() string {
	result := fmt.Sprintf("Name: %v\nExperimentCount: %.2e\n", r.Name, float64(r.ExperimentCount))
	if len(r.Errors) > 0 {
		result += "Errors: ["
	}
	for _, err := range r.Errors {
		result += err.Error() + ",\n"
	}
	if len(r.Errors) > 0 {
		result += "]\n"
	}

	result += fmt.Sprintf("avgReplacingExpense: %v, avgAwaitExpense: %v, avgFailedDevices: %v",
		r.avgReplacingExpense,
		r.avgAwaitExpense,
		r.avgFailedDevices,
	)
	result += "\n"
	return result
}

type Simulation struct {
	Name                 string
	ExperimentNumber     int
	ExperimentFactory    models.Factory[experiment.Interface]
	ExperimentResultChan chan experiment.Result
	ResultChan           chan Result
}

func (s *Simulation) Run() <-chan Result {
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(s.ExperimentNumber)
		for i := 0; i < s.ExperimentNumber; i++ {
			exp := s.ExperimentFactory.New()
			exp.Run(wg)
		}
		go func() {
			wg.Wait()
			close(s.ExperimentResultChan)
		}()

		simulationResult := Result{Errors: make([]error, 0), ExperimentCount: s.ExperimentNumber, Name: s.Name}
		for res := range s.ExperimentResultChan {
			if res.Error != nil {
				simulationResult.Errors = append(simulationResult.Errors, res.Error)
				continue
			}

			simulationResult.avgReplacingExpense += float64(res.ReplacingExpense)
			simulationResult.avgAwaitExpense += float64(res.AwaitExpense)
			simulationResult.avgFailedDevices += float64(res.FailedDevices)
		}
		simulationResult.avgReplacingExpense /= float64(s.ExperimentNumber)
		simulationResult.avgAwaitExpense /= float64(s.ExperimentNumber)
		simulationResult.avgFailedDevices /= float64(s.ExperimentNumber)

		s.ResultChan <- simulationResult
	}()
	return s.ResultChan
}

type DefaultFactory struct {
	Name                string
	ReplacementStrategy strategies.Strategy
	ResultChan          chan Result
}

func (f DefaultFactory) New() Simulation {
	expResultChan := make(chan experiment.Result)
	return Simulation{
		Name:                 f.Name,
		ExperimentNumber:     config.ExperimentNumberDefault,
		ExperimentFactory:    experiment.DefaultFactory{ReplacementStrategy: f.ReplacementStrategy, ResultChan: expResultChan},
		ExperimentResultChan: expResultChan,
		ResultChan:           f.ResultChan,
	}
}
