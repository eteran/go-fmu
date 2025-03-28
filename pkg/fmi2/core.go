package fmi2

/*
#cgo LDFLAGS: -ldl
#include "core.h"
*/
import "C"

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"unsafe"
)

//export goLogger
func goLogger(env ComponentEnvironment, instanceName C.fmi2String, status C.fmi2Status, category C.fmi2String, message *C.cchar_t) {

	_ = env
	fmt.Printf("[Name: %s, Status: %d, Category: %s] %s\n", C.GoString(instanceName), int(status), C.GoString(category), C.GoString(message))
}

//export goStepFinished
func goStepFinished(env ComponentEnvironment, status Status) {

	_ = env
	_ = status
	fmt.Println("Step Finished!")
}

type Fmu2 struct {
	Directory    string
	moduleHandle unsafe.Pointer

	// Common Functions
	getVersionPtr               unsafe.Pointer
	getTypesPlatformPtr         unsafe.Pointer
	setDebugLoggingPtr          unsafe.Pointer
	instantiatePtr              unsafe.Pointer
	freeInstancePtr             unsafe.Pointer
	setupExperimentPtr          unsafe.Pointer
	enterInitializationModePtr  unsafe.Pointer
	exitInitializationModePtr   unsafe.Pointer
	terminatePtr                unsafe.Pointer
	resetPtr                    unsafe.Pointer
	getRealPtr                  unsafe.Pointer
	getIntegerPtr               unsafe.Pointer
	getBooleanPtr               unsafe.Pointer
	getStringPtr                unsafe.Pointer
	setRealPtr                  unsafe.Pointer
	setIntegerPtr               unsafe.Pointer
	setBooleanPtr               unsafe.Pointer
	setStringPtr                unsafe.Pointer
	getFMUstatePtr              unsafe.Pointer
	setFMUstatePtr              unsafe.Pointer
	freeFMUstatePtr             unsafe.Pointer
	serializedFMUstateSizePtr   unsafe.Pointer
	serializeFMUstatePtr        unsafe.Pointer
	deSerializeFMUstatePtr      unsafe.Pointer
	getDirectionalDerivativePtr unsafe.Pointer

	// Functions for FMI2 for Model Exchange
	enterEventModePtr                unsafe.Pointer
	newDiscreteStatesPtr             unsafe.Pointer
	enterContinuousTimeModePtr       unsafe.Pointer
	completedIntegratorStepPtr       unsafe.Pointer
	setTimePtr                       unsafe.Pointer
	setContinuousStatesPtr           unsafe.Pointer
	getDerivativesPtr                unsafe.Pointer
	getEventIndicatorsPtr            unsafe.Pointer
	getContinuousStatesPtr           unsafe.Pointer
	getNominalsOfContinuousStatesPtr unsafe.Pointer

	// Functions for FMI2 for Co-Simulation
	setRealInputDerivativesPtr  unsafe.Pointer
	getRealOutputDerivativesPtr unsafe.Pointer
	doStepPtr                   unsafe.Pointer
	cancelStepPtr               unsafe.Pointer
	getStatusPtr                unsafe.Pointer
	getRealStatusPtr            unsafe.Pointer
	getIntegerStatusPtr         unsafe.Pointer
	getBooleanStatusPtr         unsafe.Pointer
	getStringStatusPtr          unsafe.Pointer
}

func (f *Fmu2) Close() error {
	if f.moduleHandle != nil {
		C.dlclose(f.moduleHandle)
		f.moduleHandle = nil
	}
	return nil
}

/* GetVersion returns the version of the "fmi2Functions.h" header file which was used to compile the functions of the FMU.
 * The function returns "fmiVersion" which is defined in this header file. The standard header file as documented in this
 * specification has version "2.0" (so this function usually returns "2.0").
 */
func (f *Fmu2) GetVersion() string {
	return C.GoString(C.GetVersion(f.getVersionPtr))
}

/* GetTypesPlatform returns the string to uniquely identify the "fmi2TypesPlatform.h" header file used for compilation of the functions of the FMU.
 * The function returns a pointer to a static string specified by "fmi2TypesPlatform" defined in this header file. The standard header file,
 * as documented in this specification, has fmi2TypesPlatform set to "default" (so this function usually returns "default").
 */
func (f *Fmu2) GetTypesPlatform() string {
	return C.GoString(C.GetVersion(f.getTypesPlatformPtr))
}

/* SetDebugLogging controls the debug logging that is output via the logger callback function by the FMU.
 * If loggingOn = true, debug logging is enabled for the log categories specified in
 * categories, otherwise it is disabled.
 * If len(categories) == 0, loggingOn applies to all log categories.
 * The allowed values of categories are defined in the modelDescription.xml file via element "<LogCategories>".
 */
func (f *Fmu2) SetDebugLogging(c *Component, loggingOn bool, categories []string) error {

	cats := Transform(categories, func(i int, v string) C.fmi2String { return C.CString(v) })

	defer func() {
		for _, v := range cats {
			C.free(unsafe.Pointer(v))
		}
	}()

	if status := Status(C.SetDebugLogging(f.setDebugLoggingPtr, c.component, toBool(loggingOn), C.size_t(len(categories)), &cats[0])); status != OK {
		return fmt.Errorf("error setting debug logging: %v", status)
	}

	return nil
}

/* FreeInstance disposes the given instance, unloads the loaded model, and frees all the allocated memory and other resources that have been allocated by the functions of the FMU interface.
 * If the passed component is nil, the function call is ignored (does not have an effect).
 */
func (f *Fmu2) FreeInstance(c *Component) {
	C.FreeInstance(f.freeInstancePtr, c.component)
}

/* Instantiate returns a new instance of an FMU. If nil is returned, then instantiation failed.
 * In that case, "functions.logger" is called with detailed information about the reason.
 * An FMU can be instantiated many times (provided capability flag canBeInstantiatedOnlyOncePerProcess = false).
 */
func (f *Fmu2) Instantiate(instanceName string, fmuType Type, fmuGuid string, resourceLocation string, visible bool, loggingOn bool) *Component {

	cInstanceName := C.CString(instanceName)
	defer C.free(unsafe.Pointer(cInstanceName))
	cType := C.fmi2Type(fmuType)
	cGuid := C.CString(fmuGuid)
	defer C.free(unsafe.Pointer(cGuid))
	cResourceLocation := C.CString(resourceLocation)
	defer C.free(unsafe.Pointer(cResourceLocation))
	cVisible := toBool(visible)
	cLoggingOn := toBool(loggingOn)

	// TODO(eteran): set the logger callback properly
	// TODO(eteran): set the step finished callback properly

	comp := C.Instantiate(f.instantiatePtr, cInstanceName, cType, cGuid, cResourceLocation, nil, cVisible, cLoggingOn)
	if comp == nil {
		return nil
	}

	return &Component{fmu: f, component: comp}
}

type SetupExperimentOption func(*SetupExperimentOptions)

type SetupExperimentOptions struct {
	relativeToleranceDefined bool
	relativeTolerance        float64
	tStopDefined             bool
	tStop                    float64
}

func WithRelativeTolerance(tolerance float64) SetupExperimentOption {
	return func(o *SetupExperimentOptions) {
		o.relativeToleranceDefined = true
		o.relativeTolerance = tolerance
	}
}

func WithStopTime(tStop float64) SetupExperimentOption {
	return func(o *SetupExperimentOptions) {
		o.tStopDefined = true
		o.tStop = tStop
	}
}

/* SetupExperiment informs the FMU to setup the experiment.
 * This function must be called after Instantiate and before EnterInitializationMode is called.
 * Arguments toleranceDefined and tolerance depend on the FMU type.
 */
func (f *Fmu2) SetupExperiment(c *Component, tStart float64, opts ...SetupExperimentOption) error {

	options := &SetupExperimentOptions{
		relativeToleranceDefined: false,
		relativeTolerance:        0.0,
		tStopDefined:             false,
		tStop:                    0.0,
	}

	for _, opt := range opts {
		opt(options)
	}

	if status := Status(C.SetupExperiment(f.setupExperimentPtr, c.component, toBool(options.relativeToleranceDefined), C.fmi2Real(options.relativeTolerance), C.fmi2Real(tStart), toBool(options.tStopDefined), C.fmi2Real(options.tStop))); status != OK {
		return fmt.Errorf("error setting up experiment: %v", status)
	}

	return nil
}

/* EnterInitializationMode informs the FMU to enter Initialization Mode.
 * Before calling this function, all variables with attribute <ScalarVariable initial = "exact" or "approx">
 * can be set with the “SetXXX” functions. Setting other variables is not allowed.
 * Furthermore, SetupExperiment must be called at least once before calling EnterInitializationMode,
 * in order that startTime is defined. */
func (f *Fmu2) EnterInitializationMode(c *Component) error {
	if status := Status(C.EnterInitializationMode(f.setupExperimentPtr, c.component)); status != OK {
		return fmt.Errorf("error entering initialization mode: %v", status)
	}

	return nil
}

/* Informs the FMU to exit Initialization Mode. For fmuType = ModelExchange,
 * this function switches off all initialization equations, and the FMU enters Event Mode implicitly;
 * that is, all continuous-time and active discrete-time equations are available.
 */
func (f *Fmu2) ExitInitializationMode(c *Component) error {
	if status := Status(C.ExitInitializationMode(f.setupExperimentPtr, c.component)); status != OK {
		return fmt.Errorf("error exiting initialization mode: %v", status)
	}

	return nil
}

/* Terminate informs the FMU that the simulation run is terminated.
 * After calling this function, the final values of all variables can be inquired with the 2GetXXX(..) functions.
 * It is not allowed to call this function after one of the functions returned with a status flag of Error or Fatal.
 */
func (f *Fmu2) Terminate(c *Component) error {
	if status := Status(C.Terminate(f.setupExperimentPtr, c.component)); status != OK {
		return fmt.Errorf("error terminating: %v", status)
	}

	return nil
}

/* Reset is called by the environment to reset the FMU after a simulation run.
 * The FMU goes into the same state as if Instantiate would have been called.
 * All variables have their default values. Before starting a new run, SetupExperiment and
 * EnterInitializationMode have to be called. */
func (f *Fmu2) Reset(c *Component) error {
	if status := Status(C.Reset(f.setupExperimentPtr, c.component)); status != OK {
		return fmt.Errorf("error resetting: %v", status)
	}

	return nil
}

/* GetReal gets actual values of variables by providing their variable references. */
func (f *Fmu2) GetReal(c *Component, vr []ValueReference) ([]float64, error) {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	values := make([]C.fmi2Real, len(vr))

	if status := Status(C.GetReal(f.getRealPtr, c.component, &refs[0], C.size_t(len(vr)), &values[0])); status != OK {
		return nil, fmt.Errorf("error getting real values: %v", status)
	}

	result := Transform(values, func(i int, v C.fmi2Real) float64 { return float64(v) })
	return result, nil
}

/* GetInteger gets actual values of variables by providing their variable references. */
func (f *Fmu2) GetInteger(c *Component, vr []ValueReference) ([]int, error) {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	values := make([]C.fmi2Integer, len(vr))

	if status := Status(C.GetInteger(f.getIntegerPtr, c.component, &refs[0], C.size_t(len(vr)), &values[0])); status != OK {
		return nil, fmt.Errorf("error getting integer values: %v", status)
	}

	result := Transform(values, func(i int, v C.fmi2Integer) int { return int(v) })
	return result, nil
}

/* GetBoolean gets actual values of variables by providing their variable references. */
func (f *Fmu2) GetBoolean(c *Component, vr []ValueReference) ([]bool, error) {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	values := make([]C.fmi2Boolean, len(vr))

	if status := Status(C.GetBoolean(f.getBooleanPtr, c.component, &refs[0], C.size_t(len(vr)), &values[0])); status != OK {
		return nil, fmt.Errorf("error getting boolean values: %v", status)
	}

	result := Transform(values, func(i int, v C.fmi2Boolean) bool { return v != 0 })
	return result, nil
}

/* GetString gets actual values of variables by providing their variable references. */
func (f *Fmu2) GetString(c *Component, vr []ValueReference) ([]string, error) {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	values := make([]C.fmi2String, len(vr))

	if status := Status(C.GetString(f.getStringPtr, c.component, &refs[0], C.size_t(len(vr)), &values[0])); status != OK {
		return nil, fmt.Errorf("error getting string values: %v", status)
	}

	result := Transform(values, func(i int, v C.fmi2String) string { return C.GoString(v) })
	return result, nil
}

/* SetReal sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (f *Fmu2) SetReal(c *Component, vr []ValueReference, value []float64) error {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	vals := Transform(value, func(i int, v float64) C.fmi2Real { return C.fmi2Real(v) })

	if status := Status(C.SetReal(f.setRealPtr, c.component, &refs[0], C.size_t(len(vr)), &vals[0])); status != OK {
		return fmt.Errorf("error setting real values: %v", status)
	}

	return nil
}

/* SetInteger sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (f *Fmu2) SetInteger(c *Component, vr []ValueReference, value []int) error {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	vals := Transform(value, func(i int, v int) C.fmi2Integer { return C.fmi2Integer(v) })

	if status := Status(C.SetInteger(f.setIntegerPtr, c.component, &refs[0], C.size_t(len(vr)), &vals[0])); status != OK {
		return fmt.Errorf("error setting integer values: %v", status)
	}

	return nil
}

/* SetBoolean sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (f *Fmu2) SetBoolean(c *Component, vr []ValueReference, value []bool) error {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	vals := Transform(value, func(i int, v bool) C.fmi2Boolean { return toBool(v) })

	if status := Status(C.SetBoolean(f.setBooleanPtr, c.component, &refs[0], C.size_t(len(vr)), &vals[0])); status != OK {
		return fmt.Errorf("error setting boolean values: %v", status)
	}

	return nil
}

/* SetString sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (f *Fmu2) SetString(c *Component, vr []ValueReference, value []string) error {

	refs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	vals := Transform(value, func(i int, v string) C.fmi2String { return C.CString(v) })

	defer func() {
		for _, v := range vals {
			C.free(unsafe.Pointer(v))
		}
	}()

	if status := Status(C.SetString(f.setStringPtr, c.component, &refs[0], C.size_t(len(vr)), &vals[0])); status != OK {
		return fmt.Errorf("error setting string values: %v", status)
	}

	return nil
}

/* GetFMUstate makes a copy of the internal FMU state and returns a pointer to this copy (FMUstate).
 * If on entry *FMUstate == nil, a new allocation is required.
 * If *FMUstate != nil, then *FMUstate points to a previously returned FMUstate that has not been modified since.
 * In particular, FreeFMUstate had not been called with this FMUstate as an argument.
 */
func (f *Fmu2) GetFMUstate(c *Component) (FmuState, error) {

	var state unsafe.Pointer
	ptr := unsafe.Pointer(&state)

	if status := Status(C.GetFMUstate(f.getFMUstatePtr, c.component, ptr)); status != OK {
		return FmuState{}, fmt.Errorf("error getting FMU state: %v", status)
	}

	return FmuState{state}, nil
}

/* SetFMUstate copies the content of the previously copied FMUstate back and uses it as actual new FMU state.
 * The FMUstate copy still exists.
 */
func (f *Fmu2) SetFMUstate(c *Component, state FmuState) error {

	ptr := state.state

	if status := Status(C.SetFMUstate(f.setFMUstatePtr, c.component, ptr)); status != OK {
		return fmt.Errorf("error setting FMU state: %v", status)
	}

	return nil
}

/* FreeFMUstate frees all memory and other resources allocated with the GetFMUstate call for this FMUstate.
 * The input argument to this function is the FMUstate to be freed.
 * If a nil pointer is provided, the call is ignored.
 * The function returns a null pointer in argument FMUstate. */
func (f *Fmu2) FreeFMUstate(c *Component, state FmuState) error {

	s := state.state
	ptr := unsafe.Pointer(&s)

	if status := Status(C.FreeFMUstate(f.freeFMUstatePtr, c.component, ptr)); status != OK {
		return fmt.Errorf("error freeing FMU state: %v", status)
	}

	return nil
}

/* SerializeFMUstate serializes the data which is referenced by pointer FMUstate and copies this data in to the returned byte slice. */
func (f *Fmu2) SerializeFMUstate(c *Component, state *FmuState) ([]byte, error) {

	var sz C.size_t
	if status := Status(C.SerializedFMUstateSize(f.serializedFMUstateSizePtr, c.component, unsafe.Pointer(state.state), &sz)); status != OK {
		return nil, fmt.Errorf("error getting serialized FMU state size: %v", status)
	}

	array := make([]byte, int(sz))
	ptr := (unsafe.Pointer(&array[0]))

	if status := Status(C.SerializeFMUstate(f.serializeFMUstatePtr, c.component, unsafe.Pointer(state.state), ptr, sz)); status != OK {
		return nil, fmt.Errorf("error serializing FMU state: %v", status)
	}

	return array, nil
}

/* DeSerializeFMUstate deserializes the byte slice, constructs a copy of the FMU state and returns FMUstate, the pointer to this copy. [The simulation is restarted at this state, when calling fmi2SetFMUState with FMUstate.] */
func (f *Fmu2) DeSerializeFMUstate(c *Component, state *FmuState, serializedState []byte) error {

	ptr := (unsafe.Pointer(&serializedState[0]))
	if status := Status(C.DeSerializeFMUstate(f.deSerializeFMUstatePtr, c.component, ptr, C.size_t(len(serializedState)), unsafe.Pointer(state.state))); status != OK {
		return fmt.Errorf("error deserializing FMU state: %v", status)
	}

	return nil
}

/* GetEventIndicators computes event indicators at the current time instant and for the current states.
 * A state event is triggered when the domain of an event indicator changes from zj > 0 to zj ≤ 0 or vice versa.
 * The FMU must guarantee that at an event restart zj ≠ 0, for example, by shifting zj with a small value.
 * Furthermore, zj should be scaled in the FMU with its nominal value (so all elements of the returned slice
 * "eventIndicators" should be in the order of "one").
 * The event indicators are returned as a slice with "ni" elements.
 */
func (f *Fmu2) GetEventIndicators(c *Component, ni int) ([]float64, error) {

	indicators := make([]C.fmi2Real, ni)

	if status := Status(C.GetEventIndicators(f.getEventIndicatorsPtr, c.component, &indicators[0], C.size_t(ni))); status != OK {
		return nil, fmt.Errorf("error getting event indicators: %v", status)
	}

	result := Transform(indicators, func(i int, v C.fmi2Real) float64 { return float64(v) })

	return result, nil
}

/* GetDerivatives computes the state derivatives at the current time instant and for the current states.
 * The derivatives are returned as a slice with "nx" elements.
 */
func (f *Fmu2) GetDerivatives(c *Component, nx int) ([]float64, error) {
	derivatives := make([]C.fmi2Real, nx)

	if status := Status(C.GetDerivatives(f.getDerivativesPtr, c.component, &derivatives[0], C.size_t(nx))); status != OK {
		return nil, fmt.Errorf("error getting derivatives: %v", status)
	}

	result := Transform(derivatives, func(i int, v C.fmi2Real) float64 { return float64(v) })

	return result, nil
}

/* SetTime sets a new time instant and re-initializes caching of variables that depend on time,
 * provided the newly provided time value is different to the previously set time value
 * (variables that depend solely on constants or parameters need not to be newly computed in the sequel,
 * but the previously computed values can be reused).
 */
func (f *Fmu2) SetTime(c *Component, time float64) error {
	if status := Status(C.SetTime(f.setTimePtr, c.component, C.fmi2Real(time))); status != OK {
		return fmt.Errorf("error setting time: %v", status)
	}

	return nil
}

/* EnterContinuousTimeMode causes the model to enter Continuous-Time Mode and all discrete-time equations
 * to become inactive and all relations to become "frozen". This function has to be called when changing from
 * Event Mode (after the global event iteration in Event Mode over all involved FMUs and other models has converged)
 * into Continuous-Time Mode.
 */
func (f *Fmu2) EnterContinuousTimeMode(c *Component) error {
	if status := Status(C.EnterContinuousTimeMode(f.enterContinuousTimeModePtr, c.component)); status != OK {
		return fmt.Errorf("error entering continuous time mode: %v", status)
	}

	return nil
}

/* DoStep causes the computation of a time step to be started. */
func (f *Fmu2) DoStep(c *Component, currentCommunicationPoint float64, communicationStepSize float64, noSetFMUStatePriorToCurrentPoint bool) error {

	if f.doStepPtr == nil {
		return fmt.Errorf("doStep function not available")
	}

	if status := Status(C.DoStep(f.doStepPtr, c.component, C.fmi2Real(currentCommunicationPoint), C.fmi2Real(communicationStepSize), toBool(noSetFMUStatePriorToCurrentPoint))); status != OK {
		return fmt.Errorf("error doing step: %v", status)
	}

	return nil
}

/* CancelStep can be called if DoStep returned Pending in order to stop the current asynchronous execution.
 * The master calls this function if, for example, the co-simulation run is stopped by the user or one of the slaves.
 * Afterwards only calls to Reset, FreeInstance, or SetFMUstate are valid to exit the step Canceled state.
 */
func (f *Fmu2) CancelStep(c *Component) error {
	if status := Status(C.CancelStep(f.cancelStepPtr, c.component)); status != OK {
		return fmt.Errorf("error canceling step: %v", status)
	}

	return nil
}

/* GetStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (f *Fmu2) GetStatus(c *Component, s StatusKind) (Status, error) {

	var value C.fmi2Status

	if status := Status(C.CancelStep(f.cancelStepPtr, c.component)); status != OK {
		return Error, fmt.Errorf("error getting status: %v", status)
	}

	return Status(value), nil
}

/* GetRealStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (f *Fmu2) GetRealStatus(c *Component, s StatusKind) (float64, error) {

	var value C.fmi2Real

	if status := Status(C.GetRealStatus(f.getRealStatusPtr, c.component, C.fmi2StatusKind(s), &value)); status != OK {
		return 0, fmt.Errorf("error getting real status: %v", status)
	}

	return float64(value), nil

}

/* GetIntegerStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (f *Fmu2) GetIntegerStatus(c *Component, s StatusKind) (int, error) {
	var value C.fmi2Integer

	if status := Status(C.GetIntegerStatus(f.getIntegerStatusPtr, c.component, C.fmi2StatusKind(s), &value)); status != OK {
		return 0, fmt.Errorf("error getting integer status: %v", status)
	}

	return int(value), nil
}

/* GetBooleanStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (f *Fmu2) GetBooleanStatus(c *Component, s StatusKind) (bool, error) {
	var value C.fmi2Boolean

	if status := Status(C.GetBooleanStatus(f.getBooleanStatusPtr, c.component, C.fmi2StatusKind(s), &value)); status != OK {
		return false, fmt.Errorf("error getting boolean status: %v", status)
	}

	return value != 0, nil
}

/* GetStringStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (f *Fmu2) GetStringStatus(c *Component, s StatusKind) (string, error) {
	var value C.fmi2String

	if status := Status(C.GetStringStatus(f.getStringStatusPtr, c.component, C.fmi2StatusKind(s), &value)); status != OK {
		return "", fmt.Errorf("error getting string status: %v", status)
	}

	str := C.GoString(value)
	C.free(unsafe.Pointer(value))
	return str, nil
}

/* EnterEventMode causes the model to enter Event Mode from the Continuous-Time Mode and
 * discrete-time equations may become active (and relations are not "frozen").
 */
func (f *Fmu2) EnterEventMode(c *Component) error {
	if status := Status(C.EnterEventMode(f.enterEventModePtr, c.component)); status != OK {
		return fmt.Errorf("error entering event mode: %v", status)
	}

	return nil
}

func (f *Fmu2) GetRealOutputDerivatives(c *Component, vr []ValueReference, order []int) ([]float64, error) {

	vrs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	orders := Transform(order, func(i int, v int) C.fmi2Integer { return C.fmi2Integer(v) })
	values := make([]C.fmi2Real, len(vr))

	if status := Status(C.GetRealOutputDerivatives(f.getRealOutputDerivativesPtr, c.component, &vrs[0], C.size_t(len(vr)), &orders[0], &values[0])); status != OK {
		return nil, fmt.Errorf("error getting real output derivatives: %v", status)
	}

	result := Transform(values, func(i int, v C.fmi2Real) float64 { return float64(v) })
	return result, nil
}

func (f *Fmu2) SetRealInputDerivatives(c *Component, vr []ValueReference, order []int, value []float64) error {

	vrs := Transform(vr, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	orders := Transform(order, func(i int, v int) C.fmi2Integer { return C.fmi2Integer(v) })
	values := Transform(value, func(i int, v float64) C.fmi2Real { return C.fmi2Real(v) })

	if status := Status(C.SetRealInputDerivatives(f.setRealInputDerivativesPtr, c.component, &vrs[0], C.size_t(len(vr)), &orders[0], &values[0])); status != OK {
		return fmt.Errorf("error setting real input derivatives: %v", status)
	}

	return nil
}

func (f *Fmu2) NewDiscreteStates(c *Component) (*EventInfo, error) {

	var eventInfo C.fmi2EventInfo

	if status := Status(C.NewDiscreteStates(f.newDiscreteStatesPtr, c.component, &eventInfo)); status != OK {
		return nil, fmt.Errorf("error getting new discrete states: %v", status)
	}

	result := &EventInfo{
		NewDiscreteStatesNeeded:           eventInfo.newDiscreteStatesNeeded != 0,
		TerminateSimulation:               eventInfo.terminateSimulation != 0,
		NominalsOfContinuousStatesChanged: eventInfo.nominalsOfContinuousStatesChanged != 0,
		ValuesOfContinuousStatesChanged:   eventInfo.valuesOfContinuousStatesChanged != 0,
		NextEventTimeDefined:              eventInfo.nextEventTimeDefined != 0,
		NextEventTime:                     float64(eventInfo.nextEventTime),
	}

	return result, nil
}

func (f *Fmu2) GetContinuousStates(c *Component, nx int) ([]float64, error) {
	states := make([]C.fmi2Real, nx)

	if status := Status(C.GetContinuousStates(f.getContinuousStatesPtr, c.component, &states[0], C.size_t(nx))); status != OK {
		return nil, fmt.Errorf("error getting continuous states: %v", status)
	}

	result := Transform(states, func(i int, v C.fmi2Real) float64 { return float64(v) })
	return result, nil
}

func (f *Fmu2) SetContinuousStates(c *Component, x []float64) error {

	states := Transform(x, func(i int, v float64) C.fmi2Real { return C.fmi2Real(v) })

	if status := Status(C.SetContinuousStates(f.setContinuousStatesPtr, c.component, &states[0], C.size_t(len(x)))); status != OK {
		return fmt.Errorf("error setting continuous states: %v", status)
	}

	return nil
}

func (f *Fmu2) GetNominalsOfContinuousStates(c *Component, nx int) ([]float64, error) {
	nominals := make([]C.fmi2Real, nx)

	if status := Status(C.GetNominalsOfContinuousStates(f.getNominalsOfContinuousStatesPtr, c.component, &nominals[0], C.size_t(nx))); status != OK {
		return nil, fmt.Errorf("error getting nominals of continuous states: %v", status)
	}

	result := Transform(nominals, func(i int, v C.fmi2Real) float64 { return float64(v) })
	return result, nil
}

func (f *Fmu2) CompletedIntegratorStep(c *Component, noSetFMUStatePriorToCurrentPoint bool) (bool, bool, error) {

	var enterEventMode C.fmi2Boolean
	var terminateSimulation C.fmi2Boolean

	if status := Status(C.CompletedIntegratorStep(f.completedIntegratorStepPtr, c.component, toBool(noSetFMUStatePriorToCurrentPoint), &enterEventMode, &terminateSimulation)); status != OK {
		return false, false, fmt.Errorf("error completing integrator step: %v", status)
	}

	return enterEventMode != 0, terminateSimulation != 0, nil
}

func (f *Fmu2) GetDirectionalDerivative(c *Component, zRef []ValueReference, vRef []ValueReference, dv []float64) ([]float64, error) {

	zRefs := Transform(zRef, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	vRefs := Transform(vRef, func(i int, v ValueReference) C.fmi2ValueReference { return C.fmi2ValueReference(v) })
	dvs := Transform(dv, func(i int, v float64) C.fmi2Real { return C.fmi2Real(v) })
	dz := make([]C.fmi2Real, len(zRef))

	if status := Status(C.GetDirectionalDerivative(f.getDirectionalDerivativePtr, c.component, &zRefs[0], C.size_t(len(zRef)), &vRefs[0], C.size_t(len(vRef)), &dvs[0], &dz[0])); status != OK {
		return nil, fmt.Errorf("error getting directional derivative: %v", status)
	}

	result := Transform(dz, func(i int, v C.fmi2Real) float64 { return float64(v) })
	return result, nil
}

func (f *Fmu2) ResourceLocation() string {
	return filepath.Join(f.Directory, "resources")
}

func resolveFunction(handle unsafe.Pointer, name string) unsafe.Pointer {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	ptr := C.dlsym(handle, str)
	return ptr

}

func New(filename string) (*Fmu2, error) {

	platforms := SupportedPlatforms(filename)
	machine := CurrentMachine()

	if !slices.Contains(platforms, machine.Platform) {
		return nil, fmt.Errorf("unsupported machine: %s", machine)
	}

	md, err := ReadModelDescription(filename, nil)
	if err != nil {
		return nil, err
	}

	library := path.Join(machine.Platform, md.ModelName+"."+machine.LibrarySuffix)

	directory, err := Extract(filename)
	if err != nil {
		return nil, err
	}

	modulePath := path.Join(directory, "binaries", library)
	moduleString := C.CString(modulePath)
	defer C.free(unsafe.Pointer(moduleString))

	handle := C.dlopen(moduleString, C.RTLD_LAZY)

	// Common Functions
	getTypesPlatformPtr := resolveFunction(handle, "fmi2GetTypesPlatform")
	getVersionPtr := resolveFunction(handle, "fmi2GetVersion")
	setDebugLoggingPtr := resolveFunction(handle, "fmi2SetDebugLogging")
	instantiatePtr := resolveFunction(handle, "fmi2Instantiate")
	freeInstancePtr := resolveFunction(handle, "fmi2FreeInstance")
	setupExperimentPtr := resolveFunction(handle, "fmi2SetupExperiment")
	enterInitializationModePtr := resolveFunction(handle, "fmi2EnterInitializationMode")
	exitInitializationModePtr := resolveFunction(handle, "fmi2ExitInitializationMode")
	terminatePtr := resolveFunction(handle, "fmi2Terminate")
	resetPtr := resolveFunction(handle, "fmi2Reset")
	getRealPtr := resolveFunction(handle, "fmi2GetReal")
	getIntegerPtr := resolveFunction(handle, "fmi2GetInteger")
	getBooleanPtr := resolveFunction(handle, "fmi2GetBoolean")
	getStringPtr := resolveFunction(handle, "fmi2GetString")
	setRealPtr := resolveFunction(handle, "fmi2SetReal")
	setIntegerPtr := resolveFunction(handle, "fmi2SetInteger")
	setBooleanPtr := resolveFunction(handle, "fmi2SetBoolean")
	setStringPtr := resolveFunction(handle, "fmi2SetString")
	getFMUstatePtr := resolveFunction(handle, "fmi2GetFMUstate")
	setFMUstatePtr := resolveFunction(handle, "fmi2SetFMUstate")
	freeFMUstatePtr := resolveFunction(handle, "fmi2FreeFMUstate")
	serializedFMUstateSizePtr := resolveFunction(handle, "fmi2SerializedFMUstateSize")
	serializeFMUstatePtr := resolveFunction(handle, "fmi2SerializeFMUstate")
	deSerializeFMUstatePtr := resolveFunction(handle, "fmi2DeSerializeFMUstate")
	getDirectionalDerivativePtr := resolveFunction(handle, "fmi2GetDirectionalDerivative")

	// Functions for FMI2 for Model Exchange
	enterEventModePtr := resolveFunction(handle, "fmi2EnterEventMode")
	newDiscreteStatesPtr := resolveFunction(handle, "fmi2NewDiscreteStates")
	enterContinuousTimeModePtr := resolveFunction(handle, "fmi2EnterContinuousTimeMode")
	completedIntegratorStepPtr := resolveFunction(handle, "fmi2CompletedIntegratorStep")
	setTimePtr := resolveFunction(handle, "fmi2SetTime")
	setContinuousStatesPtr := resolveFunction(handle, "fmi2SetContinuousStates")
	getDerivativesPtr := resolveFunction(handle, "fmi2GetDerivatives")
	getEventIndicatorsPtr := resolveFunction(handle, "fmi2GetEventIndicators")
	getContinuousStatesPtr := resolveFunction(handle, "fmi2GetContinuousStates")
	getNominalsOfContinuousStatesPtr := resolveFunction(handle, "fmi2GetNominalsOfContinuousStates")

	// Functions for FMI2 for Co-Simulation
	setRealInputDerivativesPtr := resolveFunction(handle, "fmi2SetRealInputDerivatives")
	getRealOutputDerivativesPtr := resolveFunction(handle, "fmi2GetRealOutputDerivatives")
	doStepPtr := resolveFunction(handle, "fmi2DoStep")
	cancelStepPtr := resolveFunction(handle, "fmi2CancelStep")
	getStatusPtr := resolveFunction(handle, "fmi2GetStatus")
	getRealStatusPtr := resolveFunction(handle, "fmi2GetRealStatus")
	getIntegerStatusPtr := resolveFunction(handle, "fmi2GetIntegerStatus")
	getBooleanStatusPtr := resolveFunction(handle, "fmi2GetBooleanStatus")
	getStringStatusPtr := resolveFunction(handle, "fmi2GetStringStatus")

	fmu := &Fmu2{
		Directory:    directory,
		moduleHandle: handle,

		getVersionPtr:               getVersionPtr,
		getTypesPlatformPtr:         getTypesPlatformPtr,
		setDebugLoggingPtr:          setDebugLoggingPtr,
		instantiatePtr:              instantiatePtr,
		freeInstancePtr:             freeInstancePtr,
		setupExperimentPtr:          setupExperimentPtr,
		enterInitializationModePtr:  enterInitializationModePtr,
		exitInitializationModePtr:   exitInitializationModePtr,
		terminatePtr:                terminatePtr,
		resetPtr:                    resetPtr,
		getRealPtr:                  getRealPtr,
		getIntegerPtr:               getIntegerPtr,
		getBooleanPtr:               getBooleanPtr,
		getStringPtr:                getStringPtr,
		setRealPtr:                  setRealPtr,
		setIntegerPtr:               setIntegerPtr,
		setBooleanPtr:               setBooleanPtr,
		setStringPtr:                setStringPtr,
		getFMUstatePtr:              getFMUstatePtr,
		setFMUstatePtr:              setFMUstatePtr,
		freeFMUstatePtr:             freeFMUstatePtr,
		serializedFMUstateSizePtr:   serializedFMUstateSizePtr,
		serializeFMUstatePtr:        serializeFMUstatePtr,
		deSerializeFMUstatePtr:      deSerializeFMUstatePtr,
		getDirectionalDerivativePtr: getDirectionalDerivativePtr,

		// Functions for FMI2 for Model Exchange
		enterEventModePtr:                enterEventModePtr,
		newDiscreteStatesPtr:             newDiscreteStatesPtr,
		enterContinuousTimeModePtr:       enterContinuousTimeModePtr,
		completedIntegratorStepPtr:       completedIntegratorStepPtr,
		setTimePtr:                       setTimePtr,
		setContinuousStatesPtr:           setContinuousStatesPtr,
		getDerivativesPtr:                getDerivativesPtr,
		getEventIndicatorsPtr:            getEventIndicatorsPtr,
		getContinuousStatesPtr:           getContinuousStatesPtr,
		getNominalsOfContinuousStatesPtr: getNominalsOfContinuousStatesPtr,

		// Functions for FMI2 for Co-Simulation
		setRealInputDerivativesPtr:  setRealInputDerivativesPtr,
		getRealOutputDerivativesPtr: getRealOutputDerivativesPtr,
		doStepPtr:                   doStepPtr,
		cancelStepPtr:               cancelStepPtr,
		getStatusPtr:                getStatusPtr,
		getRealStatusPtr:            getRealStatusPtr,
		getIntegerStatusPtr:         getIntegerStatusPtr,
		getBooleanStatusPtr:         getBooleanStatusPtr,
		getStringStatusPtr:          getStringStatusPtr,
	}

	runtime.SetFinalizer(fmu, (*Fmu2).Close)

	return fmu, nil
}

func toBool(x bool) C.fmi2Boolean {
	if x {
		return 1
	}
	return 0
}

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func Extract(filename string) (string, error) {
	dir, err := os.MkdirTemp("", "go-fmu-*")
	if err != nil {
		return "", err
	}

	if err := Unzip(filename, dir); err != nil {
		return "", err
	}

	return dir, nil
}
