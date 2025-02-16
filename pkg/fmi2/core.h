
#ifndef CORE_H_
#define CORE_H_

#include "headers/fmi2Functions.h"
#include <dlfcn.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>

typedef const char cchar_t;

extern void Logger(fmi2ComponentEnvironment componentEnvironment, fmi2String instanceName, fmi2Status status, fmi2String category, fmi2String message, ...);
extern void StepFinished(fmi2ComponentEnvironment componentEnvironment, fmi2Status status);
extern const char *GetTypesPlatform(void *f);
extern const char *GetVersion(void *f);
extern fmi2Status SetDebugLogging(void *f, fmi2Component component, fmi2Boolean loggingOn, size_t nCategories, const fmi2String categories[]);
extern void FreeInstance(void *f, fmi2Component component);
extern fmi2Component Instantiate(void *f, fmi2String instanceName, fmi2Type fmuType, fmi2String fmuGUID, fmi2String fmuResourceLocation, const fmi2CallbackFunctions *functions, fmi2Boolean visible, fmi2Boolean loggingOn);
extern fmi2Status SetupExperiment(void *f, fmi2Component c, fmi2Boolean relativeToleranceDefined, fmi2Real relativeTolerance, fmi2Real tStart, fmi2Boolean tStopDefined, fmi2Real tStop);
extern fmi2Status EnterInitializationMode(void *f, fmi2Component c);
extern fmi2Status ExitInitializationMode(void *f, fmi2Component c);
extern fmi2Status Terminate(void *f, fmi2Component c);
extern fmi2Status Reset(void *f, fmi2Component c);
extern fmi2Status GetReal(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Real value[]);
extern fmi2Status GetInteger(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Integer value[]);
extern fmi2Status GetBoolean(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2Boolean value[]);
extern fmi2Status GetString(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, fmi2String value[]);
extern fmi2Status GetRealOutputDerivatives(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer order[], fmi2Real value[]);
extern fmi2Status SetReal(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Real value[]);
extern fmi2Status SetInteger(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer value[]);
extern fmi2Status SetBoolean(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Boolean value[]);
extern fmi2Status SetString(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2String value[]);
extern fmi2Status SetRealInputDerivatives(void *f, fmi2Component c, const fmi2ValueReference vr[], size_t nvr, const fmi2Integer order[], const fmi2Real value[]);
extern fmi2Status GetFMUstate(void *f, fmi2Component c, void *FMUstate);
extern fmi2Status SetFMUstate(void *f, fmi2Component c, void *FMUstate);
extern fmi2Status FreeFMUstate(void *f, fmi2Component c, void *FMUstate);
extern fmi2Status SerializedFMUstateSize(void *f, fmi2Component c, void *FMUstate, size_t *size);
extern fmi2Status SerializeFMUstate(void *f, fmi2Component c, void *FMUstate, void *serializedState, size_t size);
extern fmi2Status DeSerializeFMUstate(void *f, fmi2Component c, const void *serializedState, size_t size, void *FMUstate);
extern fmi2Status DoStep(void *f, fmi2Component c, fmi2Real currentCommunicationPoint, fmi2Real communicationStepSize, fmi2Boolean noSetFMUStatePriorToCurrentPoint);
extern fmi2Status CancelStep(void *f, fmi2Component c);
extern fmi2Status GetStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Status *value);
extern fmi2Status GetRealStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Real *value);
extern fmi2Status GetIntegerStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Integer *value);
extern fmi2Status GetBooleanStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2Boolean *value);
extern fmi2Status GetStringStatus(void *f, fmi2Component c, const fmi2StatusKind s, fmi2String *value);
extern fmi2Status EnterEventMode(void *f, fmi2Component c);
extern fmi2Status EnterContinuousTimeMode(void *f, fmi2Component c);
extern fmi2Status SetTime(void *f, fmi2Component c, fmi2Real time);
extern fmi2Status GetDerivatives(void *f, fmi2Component c, fmi2Real derivatives[], size_t nx);
extern fmi2Status GetEventIndicators(void *f, fmi2Component c, fmi2Real eventIndicators[], size_t ni);
extern fmi2Status NewDiscreteStates(void *f, fmi2Component c, fmi2EventInfo *eventInfo);
extern fmi2Status GetContinuousStates(void *f, fmi2Component c, fmi2Real states[], size_t nx);
extern fmi2Status SetContinuousStates(void *f, fmi2Component c, const fmi2Real x[], size_t nx);
extern fmi2Status GetNominalsOfContinuousStates(void *f, fmi2Component c, fmi2Real x_nominal[], size_t nx);
extern fmi2Status GetDirectionalDerivative(void *f, fmi2Component c, const fmi2ValueReference z_ref[], size_t nz, const fmi2ValueReference v_ref[], size_t nv, const fmi2Real dv[], fmi2Real dz[]);
extern fmi2Status CompletedIntegratorStep(void *f, fmi2Component c, fmi2Boolean noSetFMUStatePriorToCurrentPoint, fmi2Boolean *enterEventMode, fmi2Boolean *terminateSimulation);

#endif
