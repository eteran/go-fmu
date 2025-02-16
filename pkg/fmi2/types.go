package fmi2

/*
#include "headers/fmi2Functions.h"
*/
import "C"

import "unsafe"

// Pointer to internal FMU state
type FmuState struct {
	state unsafe.Pointer
}

type (
	ComponentEnvironment unsafe.Pointer // Pointer to FMU environment
	ValueReference       uint           // handle to the value of a variable

)

type Status int

const (
	OK      Status = C.fmi2OK
	Warning        = C.fmi2Warning
	Discard        = C.fmi2Discard
	Error          = C.fmi2Error
	Fatal          = C.fmi2Fatal
	Pending        = C.fmi2Pending
)

type Type int

const (
	ModelExchangeType Type = C.fmi2ModelExchange
	CoSimulationType       = C.fmi2CoSimulation
)

type StatusKind int

const (
	DoStepStatus       StatusKind = C.fmi2DoStepStatus
	PendingStatus                 = C.fmi2PendingStatus
	LastSuccessfulTime            = C.fmi2LastSuccessfulTime
	Terminated                    = C.fmi2Terminated
)

type EventInfo struct {
	NewDiscreteStatesNeeded           bool
	TerminateSimulation               bool
	NominalsOfContinuousStatesChanged bool
	ValuesOfContinuousStatesChanged   bool
	NextEventTimeDefined              bool
	NextEventTime                     float64
}
