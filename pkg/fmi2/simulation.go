package fmi2

import (
	"errors"
	"math"
	"time"
)

// Find a nice interval that divides t into 500 - 1000 steps
func AutoInterval(t float64) float64 {

	h := math.Pow(10, (math.Round(math.Log10(t)) - 3))

	n_samples := t / h

	switch {
	case n_samples >= 2500:
		h *= 5
	case n_samples >= 2000:
		h *= 4
	case n_samples >= 1000:
		h *= 2
	case n_samples <= 200:
		h /= 5
	case n_samples <= 250:
		h /= 4
	case n_samples <= 500:
		h /= 2
	}

	return h
}

type SimulationOptions struct {
	Validate                bool              // validate the FMU and start values
	StartTime               *float64          // simulation start time (nil: use default experiment or 0 if not defined)
	StopTime                *float64          // simulation stop time (nil: use default experiment or start_time + 1 if not defined)
	Solver                  string            // solver to use for model exchange ('Euler' or 'CVode')
	StepSize                *float64          // step size for the 'Euler' solver
	RelativeTolerance       *float64          // relative tolerance for the 'CVode' solver and FMI 2.0 co-simulation FMUs
	OutputInterval          *float64          // interval for sampling the output
	RecordEvents            bool              // record outputs at events (model exchange only)
	FmiType                 string            // FMI type for the simulation ("": determine from FMU)
	StartValues             map[string]any    // mapping of variable name -> value pairs
	ApplyDefaultStartValues bool              // apply the start values from the model description (deprecated)
	Timeout                 *float64          // timeout for the simulation
	DebugLogging            bool              // enable the FMU's debug logging
	SetInputDerivatives     bool              // set the input derivatives (FMI 2.0 Co-Simulation only)
	Visible                 bool              // interactive mode (True) or batch mode (False)
	ModelDescription        *ModelDescription // the previously loaded model description (experimental)
	RemotePlatform          string            // platform to use for remoting server ('auto': determine automatically if current platform is not supported, "": no remoting; experimental)
	EarlyReturnAllowed      bool              // allow early return in FMI 3.0 Co-Simulation
	UseEventMode            bool              // use event mode in FMI 3.0 Co-Simulation if the FMU supports it
	Initialize              bool              // initialize the FMU
	Terminate               bool              // terminate the FMU
	SetStopTime             bool              // communicate the stop time to the FMU instance
	FmuInstance             *Component        // the previously instantiated FMU (experimental)

	// TODO(eteran):
	/*
		input                  a structured numpy array that contains the input (see :class:`Input`)
		output                 list of variables to record (None: record outputs)
		fmi_call_logger        callback function to log FMI calls
		logger                 callback function passed to the FMU (experimental)
		step_finished          callback to interact with the simulation (experimental)
		fmu_state              the FMU state or serialized FMU state to initialize the FMU
	*/
}

func SimulateCS(model_description *ModelDescription, fmu *Component, startTime *float64, stopTime *float64, relativeTolerance *float64, start_values map[string]any, apply_default_start_values bool /*input_signals int, output int,*/, outputInterval *float64, timeout *float64 /*step_finished int,*/, setInputDerivatives bool, use_event_mode bool, early_return_allowed bool, validate bool, initialize bool, terminate bool, set_stop_time bool) error {

	if setInputDerivatives && !model_description.CoSimulation.CanInterpolateInputs {
		return errors.New("parameter set_input_derivatives is True but the FMU cannot interpolate inputs")
	}

	if outputInterval == nil {
		interval := AutoInterval(*stopTime - *startTime)
		outputInterval = &interval
	}

	simStart := time.Now()

	canHandleVariableStepSize := model_description.CoSimulation.CanHandleVariableCommunicationStepSize

	//input = Input(fmu=fmu, modelDescription=model_description, signals=input_signals, set_input_derivatives=set_input_derivatives)

	currentTime := startTime

	if initialize {

		options := []SetupExperimentOption{}
		if stopTime != nil {
			options = append(options, WithStopTime(*stopTime))
		}
		if relativeTolerance != nil {
			options = append(options, WithRelativeTolerance(*relativeTolerance))
		}

		fmu.SetupExperiment(*startTime, options...)
		//start_values = apply_start_values(fmu, model_description, start_values, settable=settable_in_instantiated)
		fmu.EnterInitializationMode()
		//start_values = apply_start_values(fmu, model_description, start_values, settable=settable_in_initialization_mode)
		//input.apply(current_time)
		fmu.ExitInitializationMode()

	}

	stepCount := 0.0

	for {
		if timeout != nil && time.Since(simStart).Seconds() > *timeout {
			break
		}

		if *currentTime >= *stopTime {
			break
		}

		nextRegularPoint := (*startTime) + (stepCount+1.0)*(*outputInterval)

		nextCommunicationPoint := nextRegularPoint

		/*
			nextInputEventTime = input.nextEvent(current_time)
			if canHandleVariableStepSize && nextCommunicationPoint > nextInputEventTime && !FloatIsClose(nextCommunicationPoint, nextInputEventTime) {
				nextCommunicationPoint = nextInputEventTime
			}
		*/

		if nextCommunicationPoint > *stopTime && !Float64IsClose(nextCommunicationPoint, *stopTime) {
			if canHandleVariableStepSize {
				nextCommunicationPoint = *stopTime
			} else {
				break
			}
		}

		//inputEvent = FloatIsClose(nextCommunicationPoint, nextInputEventTime)

		stepSize := nextCommunicationPoint - *currentTime

		// input.apply(time, continuous=True, discrete=use_event_mode, after_event=use_event_mode)

		if err := fmu.DoStep(*currentTime, stepSize, false); err != nil {

			if true { //exception.status == fmi2Discard:

				terminateSimulation, err := fmu.GetBooleanStatus(Terminated)
				if err != nil {
					return err
				}

				if terminateSimulation {
					cTime, err := fmu.GetRealStatus(LastSuccessfulTime)
					if err != nil {
						return err
					}

					*currentTime = cTime
					//recorder.sample(time, force=True)
					break
				}
			} else {
				return err
			}
		}

		*currentTime = nextCommunicationPoint

		if Float64IsClose(*currentTime, nextRegularPoint) {
			stepCount += 1.0
		}

		/*
			if step_finished != nil && ! step_finished(time, recorder) {
				break
			}
		*/

	}

	if terminate {
		fmu.Terminate()
	}

	return nil
}

func SimulateFmu(filename string, options SimulationOptions) error {

	/*
		if fmu_instance is None and platform not in platforms and remote_platform is None:
		raise Exception(f"The current platform ({platform}) is not supported by the FMU.")
	*/

	/*
		can_sim, remote_platform = can_simulate(platforms, remote_platform)
		if not can_sim:
		raise Exception(f"The FMU cannot be simulated on the current platform ({platform}).")
	*/

	if options.ModelDescription == nil {
		md, err := ReadModelDescription(filename, &ValidationOptions{Validate: options.Validate})
		if err != nil {
			return err
		}

		options.ModelDescription = md
	}

	if options.FmiType == "" {
		if options.FmuInstance != nil {
			// options.FmiType = options.FmuInstance.FmiType
			// TODO(eteran): support this?
		} else {
			switch {
			case options.ModelDescription.ModelExchange != nil:
				options.FmiType = "ModelExchange"
			case options.ModelDescription.CoSimulation != nil:
				options.FmiType = "CoSimulation"
			}
		}
	}

	if options.FmiType != "ModelExchange" && options.FmiType != "CoSimulation" {
		return errors.New("FmiType must be one of 'ModelExchange' or 'CoSimulation'")
	}

	if !options.Initialize {
		if options.FmiType != "CoSimulation" {
			return errors.New("if initialize is False, the interface type must be 'CoSimulation'")
		}

		// TODO(eteran): support FmuState
		if options.FmuInstance == nil {
			return errors.New("if initialize is False, FmuInstance or FmuState must be provided")
		}
	}

	experiment := options.ModelDescription.DefaultExperiment
	if options.StartTime == nil {
		if experiment != nil && experiment.StartTime != nil {
			options.StartTime = experiment.StartTime
		} else {
			startTime := 0.0
			options.StartTime = &startTime
		}
	}

	if options.StopTime == nil {
		if experiment != nil && experiment.StopTime != nil {
			options.StopTime = experiment.StopTime
		} else {
			stopTime := *options.StartTime + 1.0
			options.StopTime = &stopTime
		}
	}

	if options.RelativeTolerance == nil && experiment != nil {
		options.RelativeTolerance = experiment.Tolerance
	}

	if options.StepSize == nil {
		totalTime := *options.StopTime - *options.StartTime
		stepSize := math.Pow(10.0, math.Round(math.Log10(totalTime))-3)
		options.StepSize = &stepSize
	}

	if options.OutputInterval == nil && options.FmiType == "CoSimulation" {
		/*
			coSimulation := options.ModelDescription.CoSimulation
				if coSimulation != nil && coSimulation.FixedInternalStepSize != nil {
					options.OutputInterval = coSimulation.FixedInternalStepSize
				} else */
		if experiment != nil && experiment.StepSize != nil {
			options.OutputInterval = experiment.StepSize
		}

		if options.OutputInterval != nil {
			for (*options.StopTime-*options.StartTime) / *options.OutputInterval > 1000 {
				*options.OutputInterval *= 2
			}
		}
	}

	fmu, err := New(filename)
	if err != nil {
		return err
	}

	defer fmu.Close()

	switch options.FmiType {
	case "ModelExchange":
		//result = simulateME(model_description, fmu, start_time, stop_time, solver, step_size, relative_tolerance, start_values, apply_default_start_values, input, output, output_interval, record_events, timeout, step_finished, validate, set_stop_time)
	case "CoSimulation":

		md, err := ReadModelDescription(filename, &ValidationOptions{Validate: options.Validate})
		if err != nil {
			return err
		}

		comp := fmu.Instantiate(
			md.ModelName,
			CoSimulationType,
			md.Guid,
			fmu.ResourceLocation(),
			options.Visible,
			options.DebugLogging)

		result := SimulateCS(
			options.ModelDescription,
			comp,
			options.StartTime,
			options.StopTime,
			options.RelativeTolerance,
			options.StartValues,
			options.ApplyDefaultStartValues,
			//input,
			//output,
			options.OutputInterval,
			options.Timeout,
			//step_finished,
			options.SetInputDerivatives,
			options.UseEventMode,
			options.EarlyReturnAllowed,
			options.Validate,
			options.Initialize,
			options.Terminate,
			options.SetStopTime)

		_ = result
	}

	return nil
}
