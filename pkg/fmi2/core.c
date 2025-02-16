#include "headers/fmi2Functions.h"
#include <dlfcn.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>

typedef const char cchar_t;

extern void goLogger(void *user_data, fmi2String instanceName, fmi2Status status, fmi2String category, cchar_t *message);
extern void goStepFinished(void *user_data, fmi2Status status);

void Logger(fmi2ComponentEnvironment componentEnvironment, fmi2String instanceName, fmi2Status status, fmi2String category, fmi2String message, ...) {

	char buffer[1024];
	va_list ap;
	va_start(ap, message);
	vsnprintf(buffer, sizeof(buffer), message, ap);
	va_end(ap);

	goLogger(componentEnvironment, instanceName, status, category, buffer);
}

void StepFinished(fmi2ComponentEnvironment componentEnvironment, fmi2Status status) {
	goStepFinished(componentEnvironment, status);
}

const char *GetTypesPlatform(void *f) {
	return ((fmi2GetTypesPlatformTYPE *)f)();
}

const char *GetVersion(void *f) {
	return ((fmi2GetVersionTYPE *)f)();
}

fmi2Status SetDebugLogging(void *f, fmi2Component component, fmi2Boolean loggingOn, size_t nCategories, const fmi2String categories[]) {
	return ((fmi2SetDebugLoggingTYPE *)f)(component, loggingOn, nCategories, categories);
}

void FreeInstance(void *f, fmi2Component component) {
	((fmi2FreeInstanceTYPE *)f)(component);
}

fmi2Component Instantiate(void *f, fmi2String instanceName, fmi2Type fmuType, fmi2String fmuGUID, fmi2String fmuResourceLocation, const fmi2CallbackFunctions *functions, fmi2Boolean visible, fmi2Boolean loggingOn) {

	fmi2CallbackFunctions callbacks;
	callbacks.logger               = Logger;
	callbacks.allocateMemory       = calloc;
	callbacks.freeMemory           = free;
	callbacks.stepFinished         = NULL;
	callbacks.componentEnvironment = NULL;

	return ((fmi2InstantiateTYPE *)f)(instanceName, fmuType, fmuGUID, fmuResourceLocation, &callbacks, visible, loggingOn);
}

fmi2Status SetupExperiment(void *f, fmi2Component c, fmi2Boolean relativeToleranceDefined, fmi2Real relativeTolerance, fmi2Real tStart, fmi2Boolean tStopDefined, fmi2Real tStop) {
	return ((fmi2SetupExperimentTYPE *)f)(c, relativeToleranceDefined, relativeTolerance, tStart, tStopDefined, tStop);
}

fmi2Status EnterInitializationMode(void *f, fmi2Component c) {
	return ((fmi2EnterInitializationModeTYPE *)f)(c);
}

fmi2Status ExitInitializationMode(void *f, fmi2Component c) {
	return ((fmi2ExitInitializationModeTYPE *)f)(c);
}

fmi2Status Terminate(void *f, fmi2Component c) {
	return ((fmi2TerminateTYPE *)f)(c);
}

fmi2Status Reset(void *f, fmi2Component c) {
	return ((fmi2ResetTYPE *)f)(c);
}

fmi2Status GetReal(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Real value[]) {
	return ((fmi2GetRealTYPE *)f)(c, vr, nvr, value);
}

fmi2Status GetInteger(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Integer value[]) {
	return ((fmi2GetIntegerTYPE *)f)(c, vr, nvr, value);
}

fmi2Status GetBoolean(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Boolean value[]) {
	return ((fmi2GetBooleanTYPE *)f)(c, vr, nvr, value);
}

fmi2Status GetString(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2String value[]) {
	return ((fmi2GetStringTYPE *)f)(c, vr, nvr, value);
}

fmi2Status GetRealOutputDerivatives(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer order[], fmi2Real value[]) {
	return ((fmi2GetRealOutputDerivativesTYPE *)f)(c, vr, nvr, order, value);
}

fmi2Status SetReal(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Real value[]) {
	return ((fmi2SetRealTYPE *)f)(c, vr, nvr, value);
}

fmi2Status SetInteger(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer value[]) {
	return ((fmi2SetIntegerTYPE *)f)(c, vr, nvr, value);
}

fmi2Status SetBoolean(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Boolean value[]) {
	return ((fmi2SetBooleanTYPE *)f)(c, vr, nvr, value);
}

fmi2Status SetString(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2String value[]) {
	return ((fmi2SetStringTYPE *)f)(c, vr, nvr, value);
}

fmi2Status SetRealInputDerivatives(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer order[], const fmi2Real value[]) {
	return ((fmi2SetRealInputDerivativesTYPE *)f)(c, vr, nvr, order, value);
}

fmi2Status GetFMUstate(void *f, fmi2Component c, void *FMUstate) {
	return ((fmi2GetFMUstateTYPE *)f)(c, FMUstate);
}

fmi2Status SetFMUstate(void *f, fmi2Component c, void *FMUstate) {
	return ((fmi2SetFMUstateTYPE *)f)(c, FMUstate);
}

fmi2Status FreeFMUstate(void *f, fmi2Component c, void *FMUstate) {
	return ((fmi2FreeFMUstateTYPE *)f)(c, FMUstate);
}

fmi2Status SerializedFMUstateSize(void *f, fmi2Component c, void *FMUstate, size_t *size) {
	return ((fmi2SerializedFMUstateSizeTYPE *)f)(c, FMUstate, size);
}

fmi2Status SerializeFMUstate(void *f, fmi2Component c, void *FMUstate, void *serializedState, size_t size) {
	return ((fmi2SerializeFMUstateTYPE *)f)(c, FMUstate, serializedState, size);
}

fmi2Status DeSerializeFMUstate(void *f, fmi2Component c, const void *serializedState, size_t size, void *FMUstate) {
	return ((fmi2DeSerializeFMUstateTYPE *)f)(c, serializedState, size, FMUstate);
}

fmi2Status DoStep(void *f, fmi2Component c, fmi2Real currentCommunicationPoint, fmi2Real communicationStepSize, fmi2Boolean noSetFMUStatePriorToCurrentPoint) {
	return ((fmi2DoStepTYPE *)f)(c, currentCommunicationPoint, communicationStepSize, noSetFMUStatePriorToCurrentPoint);
}

fmi2Status CancelStep(void *f, fmi2Component c) {
	return ((fmi2CancelStepTYPE *)f)(c);
}

fmi2Status GetStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Status *value) {
	return ((fmi2GetStatusTYPE *)f)(c, s, value);
}

fmi2Status GetRealStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Real *value) {
	return ((fmi2GetRealStatusTYPE *)f)(c, s, value);
}

fmi2Status GetIntegerStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Integer *value) {
	return ((fmi2GetIntegerStatusTYPE *)f)(c, s, value);
}

fmi2Status GetBooleanStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Boolean *value) {
	return ((fmi2GetBooleanStatusTYPE *)f)(c, s, value);
}

fmi2Status GetStringStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2String *value) {
	return ((fmi2GetStringStatusTYPE *)f)(c, s, value);
}

fmi2Status EnterEventMode(void *f, fmi2Component c) {
	return ((fmi2EnterEventModeTYPE *)f)(c);
}

fmi2Status EnterContinuousTimeMode(void *f, fmi2Component c) {
	return ((fmi2EnterContinuousTimeModeTYPE *)f)(c);
}

fmi2Status SetTime(void *f, fmi2Component c, fmi2Real time) {
	return ((fmi2SetTimeTYPE *)f)(c, time);
}

fmi2Status GetDerivatives(void *f, fmi2Component c, fmi2Real derivatives[], size_t nx) {
	return ((fmi2GetDerivativesTYPE *)f)(c, derivatives, nx);
}

fmi2Status GetEventIndicators(void *f, fmi2Component c, fmi2Real eventIndicators[], size_t ni) {
	return ((fmi2GetEventIndicatorsTYPE *)f)(c, eventIndicators, ni);
}

fmi2Status NewDiscreteStates(void *f, fmi2Component c, fmi2EventInfo *eventInfo) {
	return ((fmi2NewDiscreteStatesTYPE *)f)(c, eventInfo);
}

fmi2Status GetContinuousStates(void *f, fmi2Component c, fmi2Real states[], size_t nx) {
	return ((fmi2GetContinuousStatesTYPE *)f)(c, states, nx);
}

fmi2Status SetContinuousStates(void *f, fmi2Component c, const fmi2Real x[], size_t nx) {
	return ((fmi2SetContinuousStatesTYPE *)f)(c, x, nx);
}

fmi2Status GetNominalsOfContinuousStates(void *f, fmi2Component c, fmi2Real x_nominal[], size_t nx) {
	return ((fmi2GetNominalsOfContinuousStatesTYPE *)f)(c, x_nominal, nx);
}

fmi2Status GetDirectionalDerivative(void *f, fmi2Component c, const fmi2ValueReference z_ref[], size_t nz, const fmi2ValueReference v_ref[], size_t nv, const fmi2Real dv[], fmi2Real dz[]) {
	return ((fmi2GetDirectionalDerivativeTYPE *)f)(c, z_ref, nz, v_ref, nv, dv, dz);
}

fmi2Status CompletedIntegratorStep(void *f, fmi2Component c, fmi2Boolean noSetFMUStatePriorToCurrentPoint, fmi2Boolean *enterEventMode, fmi2Boolean *terminateSimulation) {
	return ((fmi2CompletedIntegratorStepTYPE *)f)(c, noSetFMUStatePriorToCurrentPoint, enterEventMode, terminateSimulation);
}
