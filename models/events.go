package models

type ExperimentTickEventType int

const (
	TimeIncrement ExperimentTickEventType = iota
	Stop
)

type ExperimentTickEvent struct {
	Type ExperimentTickEventType
}
