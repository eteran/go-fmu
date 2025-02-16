package fmi2

/*
#include "headers/fmi2Functions.h"
*/
import "C"

type Component struct {
	fmu       *Fmu2
	component C.fmi2Component
}

/* SetDebugLogging controls the debug logging that is output via the logger callback function by the FMU.
 * If loggingOn = true, debug logging is enabled for the log categories specified in
 * categories, otherwise it is disabled.
 * If len(categories) == 0, loggingOn applies to all log categories.
 * The allowed values of categories are defined in the modelDescription.xml file via element "<LogCategories>".
 */
func (c *Component) SetDebugLogging(loggingOn bool, categories []string) error {
	return c.fmu.SetDebugLogging(c, loggingOn, categories)
}

/* SetupExperiment informs the FMU to setup the experiment.
 * This function must be called after Instantiate and before EnterInitializationMode is called.
 * Arguments toleranceDefined and tolerance depend on the FMU type.
 */
func (c *Component) SetupExperiment(tStart float64, opts ...SetupExperimentOption) error {
	return c.fmu.SetupExperiment(c, tStart, opts...)
}

/* EnterInitializationMode informs the FMU to enter Initialization Mode.
 * Before calling this function, all variables with attribute <ScalarVariable initial = "exact" or "approx">
 * can be set with the “SetXXX” functions. Setting other variables is not allowed.
 * Furthermore, SetupExperiment must be called at least once before calling EnterInitializationMode,
 * in order that startTime is defined. */
func (c *Component) EnterInitializationMode() error {
	return c.fmu.EnterInitializationMode(c)
}

/* Informs the FMU to exit Initialization Mode. For fmuType = ModelExchange,
 * this function switches off all initialization equations, and the FMU enters Event Mode implicitly;
 * that is, all continuous-time and active discrete-time equations are available.
 */
func (c *Component) ExitInitializationMode() error {
	return c.fmu.ExitInitializationMode(c)
}

/* Terminate informs the FMU that the simulation run is terminated.
 * After calling this function, the final values of all variables can be inquired with the 2GetXXX(..) functions.
 * It is not allowed to call this function after one of the functions returned with a status flag of Error or Fatal.
 */
func (c *Component) Terminate() error {
	return c.fmu.Terminate(c)
}

/* Reset is called by the environment to reset the FMU after a simulation run.
 * The FMU goes into the same state as if Instantiate would have been called.
 * All variables have their default values. Before starting a new run, SetupExperiment and
 * EnterInitializationMode have to be called. */
func (c *Component) Reset() error {
	return c.fmu.Terminate(c)
}

/* GetReal gets actual values of variables by providing their variable references. */
func (c *Component) GetReal(vr []ValueReference) ([]float64, error) {
	return c.fmu.GetReal(c, vr)
}

/* GetInteger gets actual values of variables by providing their variable references. */
func (c *Component) GetInteger(vr []ValueReference) ([]int, error) {
	return c.fmu.GetInteger(c, vr)
}

/* GetBoolean gets actual values of variables by providing their variable references. */
func (c *Component) GetBoolean(vr []ValueReference) ([]bool, error) {
	return c.fmu.GetBoolean(c, vr)
}

/* GetString gets actual values of variables by providing their variable references. */
func (c *Component) GetString(vr []ValueReference) ([]string, error) {
	return c.fmu.GetString(c, vr)
}

/* SetReal sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (c *Component) SetReal(vr []ValueReference, value []float64) error {
	return c.fmu.SetReal(c, vr, value)
}

/* SetInteger sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (c *Component) SetInteger(vr []ValueReference, value []int) error {
	return c.fmu.SetInteger(c, vr, value)
}

/* SetBoolean sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (c *Component) SetBoolean(vr []ValueReference, value []bool) error {
	return c.fmu.SetBoolean(c, vr, value)
}

/* SetString sets parameters, inputs, and start values, and re-initializes caching of variables that depend on these variables */
func (c *Component) SetString(vr []ValueReference, value []string) error {
	return c.fmu.SetString(c, vr, value)
}

/* GetFMUstate makes a copy of the internal FMU state and returns a pointer to this copy (FMUstate).
 * If on entry *FMUstate == nil, a new allocation is required.
 * If *FMUstate != nil, then *FMUstate points to a previously returned FMUstate that has not been modified since.
 * In particular, FreeFMUstate had not been called with this FMUstate as an argument.
 */
func (c *Component) GetFMUstate() (FmuState, error) {
	return c.fmu.GetFMUstate(c)
}

/* SetFMUstate copies the content of the previously copied FMUstate back and uses it as actual new FMU state.
 * The FMUstate copy still exists.
 */
func (c *Component) SetFMUstate(state FmuState) error {
	return c.fmu.SetFMUstate(c, state)
}

/* FreeFMUstate frees all memory and other resources allocated with the GetFMUstate call for this FMUstate.
 * The input argument to this function is the FMUstate to be freed.
 * If a nil pointer is provided, the call is ignored.
 * The function returns a null pointer in argument FMUstate. */
func (c *Component) FreeFMUstate(state FmuState) error {
	return c.fmu.FreeFMUstate(c, state)
}

/* SerializeFMUstate serializes the data which is referenced by pointer FMUstate and copies this data in to the returned byte slice. */
func (c *Component) SerializeFMUstate(state *FmuState) ([]byte, error) {
	return c.fmu.SerializeFMUstate(c, state)
}

/* DeSerializeFMUstate deserializes the byte slice, constructs a copy of the FMU state and returns FMUstate, the pointer to this copy. [The simulation is restarted at this state, when calling fmi2SetFMUState with FMUstate.] */
func (c *Component) DeSerializeFMUstate(state *FmuState, serializedState []byte) error {
	return c.fmu.DeSerializeFMUstate(c, state, serializedState)
}

/* GetEventIndicators computes event indicators at the current time instant and for the current states.
 * A state event is triggered when the domain of an event indicator changes from zj > 0 to zj ≤ 0 or vice versa.
 * The FMU must guarantee that at an event restart zj ≠ 0, for example, by shifting zj with a small value.
 * Furthermore, zj should be scaled in the FMU with its nominal value (so all elements of the returned slice
 * "eventIndicators" should be in the order of "one").
 * The event indicators are returned as a slice with "ni" elements.
 */
func (c *Component) GetEventIndicators(ni int) ([]float64, error) {
	return c.fmu.GetEventIndicators(c, ni)
}

/* GetDerivatives computes the state derivatives at the current time instant and for the current states.
 * The derivatives are returned as a slice with "nx" elements.
 */
func (c *Component) GetDerivatives(nx int) ([]float64, error) {
	return c.fmu.GetDerivatives(c, nx)
}

/* SetTime sets a new time instant and re-initializes caching of variables that depend on time,
 * provided the newly provided time value is different to the previously set time value
 * (variables that depend solely on constants or parameters need not to be newly computed in the sequel,
 * but the previously computed values can be reused).
 */
func (c *Component) SetTime(time float64) error {
	return c.fmu.SetTime(c, time)
}

/* EnterContinuousTimeMode causes the model to enter Continuous-Time Mode and all discrete-time equations
 * to become inactive and all relations to become "frozen". This function has to be called when changing from
 * Event Mode (after the global event iteration in Event Mode over all involved FMUs and other models has converged)
 * into Continuous-Time Mode.
 */
func (c *Component) EnterContinuousTimeMode() error {
	return c.fmu.EnterContinuousTimeMode(c)
}

/* DoStep causes the computation of a time step to be started. */
func (c *Component) DoStep(currentCommunicationPoint float64, communicationStepSize float64, noSetFMUStatePriorToCurrentPoint bool) error {
	return c.fmu.DoStep(c, currentCommunicationPoint, communicationStepSize, noSetFMUStatePriorToCurrentPoint)
}

/* CancelStep can be called if DoStep returned Pending in order to stop the current asynchronous execution.
 * The master calls this function if, for example, the co-simulation run is stopped by the user or one of the slaves.
 * Afterwards only calls to Reset, FreeInstance, or SetFMUstate are valid to exit the step Canceled state.
 */
func (c *Component) CancelStep() error {
	return c.fmu.CancelStep(c)
}

/* GetStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (c *Component) GetStatus(s StatusKind) (Status, error) {
	return c.fmu.GetStatus(c, s)
}

/* GetRealStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (c *Component) GetRealStatus(s StatusKind) (float64, error) {
	return c.fmu.GetRealStatus(c, s)
}

/* GetIntegerStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (c *Component) GetIntegerStatus(s StatusKind) (int, error) {
	return c.fmu.GetIntegerStatus(c, s)
}

/* GetBooleanStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (c *Component) GetBooleanStatus(s StatusKind) (bool, error) {
	return c.fmu.GetBooleanStatus(c, s)
}

/* GetStringStatus informs the master about the actual status of the simulation run.
 * Which status information is to be returned is specified by the argument StatusKind.
 * It depends on the capabilities of the slave which status information can be given by the slave.
 * If a status is required which cannot be retrieved by the slave it returns Discard.
 */
func (c *Component) GetStringStatus(s StatusKind) (string, error) {
	return c.fmu.GetStringStatus(c, s)
}

/* EnterEventMode causes the model to enter Event Mode from the Continuous-Time Mode and
 * discrete-time equations may become active (and relations are not "frozen").
 */
func (c *Component) EnterEventMode() error {
	return c.fmu.EnterEventMode(c)
}

/*
 * GetRealOutputDerivatives retrieves the n-th derivative of output values.
 * vr is a slice of value references that define the variables whose derivatives shall be retrieved.
 * order contains the order of the respective derivative
 * (1 means the first derivative, 0 is not allowed).
 * returns a slice the actual values of the derivatives.
 * Restrictions on using the function are the same as for the GetReal function.
 */
func (c *Component) GetRealOutputDerivatives(vr []ValueReference, order []int) ([]float64, error) {
	return c.fmu.GetRealOutputDerivatives(c, vr, order)
}

/*
 * SetRealInputDerivatives sets the n-th time derivative of real input variables.
 * vr is a slice of value references that define the variables whose derivatives shall be set.
 * order contains the orders of the respective derivative
 * (1 means the first derivative, 0, is not allowed).
 * value is a slice with the values of the derivatives.
 * Different input variables may have different interpolation order.
 * Restrictions on using the function are the same as for the SetReal function.
 */
func (c *Component) SetRealInputDerivatives(vr []ValueReference, order []int, value []float64) error {
	return c.fmu.SetRealInputDerivatives(c, vr, order, value)
}

/* The FMU is in Event Mode and the super dense time is incremented by this call.
 * If the super dense time before a call to NewDiscreteStates was (tR,tI),
 * then the time instant after the call is (tR,tI+1).
 * If return argument EventInfo->NewDiscreteStatesNeeded = true, the FMU should stay in Event Mode,
 * and the FMU requires to set new inputs to the FMU (SetXXX on inputs) to compute and get the outputs
 * (GetXXX on outputs) and to call fmi2NewDiscreteStates again.
 * Depending on the connection with other FMUs, the environment shall
 * - call Terminate, if TerminateSimulation = true is returned by at least one FMU.
 * - call EnterContinuousTimeMode if all FMUs return NewDiscreteStatesNeeded = false.
 * - stay in Event Mode otherwise.
 */
func (c *Component) NewDiscreteStates() (*EventInfo, error) {
	return c.fmu.NewDiscreteStates(c)
}

/* GetContinuousStates returns the new (continuous) state vector x.
 */
func (c *Component) GetContinuousStates(nx int) ([]float64, error) {
	return c.fmu.GetContinuousStates(c, nx)
}

/* SetContinuousStates sets a new (continuous) state vector and re-initializes caching of
 * variables that depend on the states. (variables that depend solely on constants,
 * parameters, time, and inputs do not need to be newly computed in the sequel,
 * but the previously computed values can be reused).
 * Note that the continuous states might also be changed in Event Mode.
 * Note: Status = Discard is possible.
 */
func (c *Component) SetContinuousStates(x []float64) error {
	return c.fmu.SetContinuousStates(c, x)
}

/* GetNominalsOfContinuousStates returns the nominal values of the continuous states.
 * This function should always be called after calling function NewDiscreteStates if it
 * returns with EventInfo->NominalsOfContinuousStatesChanged = true, since then the nominal
 * values of the continuous states have changed
 * [for example, because the association of the continuous states to variables has changed due to internal dynamic state selection].
 * If the FMU does not have information about the  nominal value of a continuous state i,
 * a nominal value x_nominal[i] = 1.0 should be returned.
 * Note that it is required that x_nominal[i] > 0.0
 */
func (c *Component) GetNominalsOfContinuousStates(nx int) ([]float64, error) {
	return c.fmu.GetNominalsOfContinuousStates(c, nx)
}

/* CompletedIntegratorStep must be called by the environment after every completed step of the integrator
 * provided the capability flag CompletedIntegratorStepNotNeeded = false.
 * Argument noSetFMUStatePriorToCurrentPoint is true if SetFMUState will no longer be called for
 * time instants prior to current time in this simulation run [the FMU can use this flag to flush a result buffer].
 */
func (c *Component) CompletedIntegratorStep(noSetFMUStatePriorToCurrentPoint bool) (bool, bool, error) {
	return c.fmu.CompletedIntegratorStep(c, noSetFMUStatePriorToCurrentPoint)
}

/* GetDirectionalDerivative computes the directional derivatives of an FMU. An FMU has different Modes
 * and in every Mode an FMU might be described by different equations and different unknowns.
 * The precise definitions are given in the mathematical descriptions of Model Exchange and Co-Simulation
 */
func (c *Component) GetDirectionalDerivative(zRef []ValueReference, vRef []ValueReference, dv []float64) ([]float64, error) {
	return c.fmu.GetDirectionalDerivative(c, zRef, vRef, dv)
}

/* FreeInstance disposes the given instance, unloads the loaded model, and frees all the allocated memory
 * and other resources that have been allocated by the functions of the FMU interface.
 */
func (c *Component) FreeInstance() {
	c.fmu.FreeInstance(c)
}
