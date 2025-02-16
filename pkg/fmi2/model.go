package fmi2

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"path"
	"slices"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValidationOptions struct {
	Validate               bool
	ValidateVariableNames  bool
	ValidateModelStructure bool
}

/*
Read the model description from an FMU without extracting it

Parameters:

	filename filename of the FMU
	options  an instance of ValidationOptions or nil. If nil, uses a default instance where
	         Validate = true, ValidateVariableNames = false, ValidateModelStructure = false

Returns:

	a ModelDescription object
*/
func ReadModelDescription(filename string, options *ValidationOptions) (*ModelDescription, error) {

	if options == nil {
		options = &ValidationOptions{
			Validate:               true,
			ValidateVariableNames:  false,
			ValidateModelStructure: false,
		}
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}

	f, err := r.Open("modelDescription.xml")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var md ModelDescription
	reader := bufio.NewReader(f)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	if err = decoder.Decode(&md); err != nil {
		return nil, err
	}

	is_fmi1 := md.FmiVersion == "1.0"
	is_fmi2 := md.FmiVersion == "2.0"
	is_fmi3 := strings.HasPrefix(md.FmiVersion, "3.")

	if !is_fmi1 && !is_fmi2 && !is_fmi3 {
		return nil, fmt.Errorf("unsupported FMI version: %s", md.FmiVersion)
	}

	if options.Validate {
		if err := md.Validate(); err != nil {
			return nil, err
		}
	}

	if options.ValidateVariableNames {
		if err := md.ValidateVariableNames(); err != nil {
			return nil, err
		}
	}

	if options.ValidateModelStructure {
		if err := md.ValidateStructure(); err != nil {
			return nil, err
		}
	}

	// TODO(eteran): normalize other FMI versions ?
	// FMI3: GUID -> instantiationToken
	// FMI2: len(ModelStructure/Derivatives/Unknown) -> numberOfContinuousStates

	return &md, nil
}

/*
Determine the supported platforms for the FMU

Parameters:

	filename filename of the FMU

Returns:

	a slice of strings representing the supported platforms, or an empty slice on error
*/
func SupportedPlatforms(filename string) []string {

	platforms := make([]string, 0)

	r, err := zip.OpenReader(filename)
	if err != nil {
		return platforms
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range r.File {

		dir := path.Dir(f.Name)
		base := path.Base(dir)

		switch {
		case strings.HasSuffix(f.Name, ".dylib") && strings.HasPrefix(dir, "binaries/"):
			platforms = append(platforms, base)
		case strings.HasSuffix(f.Name, ".so") && strings.HasPrefix(dir, "binaries/"):
			platforms = append(platforms, base)
		case strings.HasSuffix(f.Name, ".dll") && strings.HasPrefix(dir, "binaries/"):
			platforms = append(platforms, base)
		}
	}

	return platforms
}

/*
Format the info for an FMU

Parameters:

	filename     filename of the FMU
	causalities  the causalities of the variables to include

Returns

	The info as a multi line string
*/
func FmuInfo(filename string, causalities []string) (string, error) {

	md, err := ReadModelDescription(filename, nil)
	if err != nil {
		return "", err
	}

	platforms := SupportedPlatforms(filename)

	fmiTypes := make([]string, 0)
	if md.ModelExchange != nil {
		fmiTypes = append(fmiTypes, "Model Exchange")
	}

	if md.CoSimulation != nil {
		fmiTypes = append(fmiTypes, "Co-Simulation")
	}

	var sb strings.Builder

	sb.WriteString("Model Info\n\n")
	sb.WriteString(fmt.Sprintf("  FMI Version        %s\n", md.FmiVersion))
	sb.WriteString(fmt.Sprintf("  FMI Type           %s\n", strings.Join(fmiTypes, ", ")))
	sb.WriteString(fmt.Sprintf("  Model Name         %s\n", md.ModelName))
	sb.WriteString(fmt.Sprintf("  Description        %s\n", md.Description))
	sb.WriteString(fmt.Sprintf("  Platforms          %s\n", strings.Join(platforms, ", ")))
	sb.WriteString(fmt.Sprintf("  Continuous States  %d\n", -1)) // ?
	sb.WriteString(fmt.Sprintf("  Event Indicators   %d\n", md.NumberOfEventIndicators))
	sb.WriteString(fmt.Sprintf("  Variables          %d\n", len(md.ModelVariables.ScalarVariable)))
	sb.WriteString(fmt.Sprintf("  Generation Tool    %s\n", md.GenerationTool))
	sb.WriteString(fmt.Sprintf("  Generation Date    %s\n", md.GenerationDateAndTime))

	if md.DefaultExperiment != nil {
		sb.WriteString("\nDefault Experiment\n\n")
		if md.DefaultExperiment.StartTime != nil {
			sb.WriteString(fmt.Sprintf("  Start Time    %g\n", *md.DefaultExperiment.StartTime))
		}
		if md.DefaultExperiment.StopTime != nil {
			sb.WriteString(fmt.Sprintf("  Stop Time     %g\n", *md.DefaultExperiment.StopTime))
		}
		if md.DefaultExperiment.Tolerance != nil {
			sb.WriteString(fmt.Sprintf("  Tolerance     %g\n", *md.DefaultExperiment.Tolerance))
		}
		if md.DefaultExperiment.StepSize != nil {
			sb.WriteString(fmt.Sprintf("  Step Size     %g\n", *md.DefaultExperiment.StepSize))
		}
	}

	/*
		inputs := make([]string, 0)
		outputs := make([]string, 0)

		for _, v := range md.ModelVariables.ScalarVariable {
			switch v.Causality {
			case "input":
				inputs = append(inputs, v.Name)
			case "output":
				outputs = append(outputs, v.Name)
			}
		}
	*/

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Variables (%s)\n\n", strings.Join(causalities, ", ")))
	sb.WriteString(fmt.Sprintf("  %-18s %-10s %-12s %-8s %s\n", "Name", "Causality", "Start Value", "Unit", "Description"))
	for _, v := range md.ModelVariables.ScalarVariable {
		if !slices.Contains(causalities, v.Causality) {
			continue
		}

		name := v.Name
		if len(name) > 18 {
			name = "..." + name[18-3:]
		}

		startValue := ""
		units := ""

		switch {
		case v.Real != nil:
			startValue = fmt.Sprint(v.Real.Start)
			units = v.Real.Unit
		case v.Boolean != nil:
			startValue = fmt.Sprint(v.Boolean.Start)
			units = v.Boolean.DeclaredType
		case v.Integer != nil:
			startValue = fmt.Sprint(v.Integer.Start)
			units = v.Integer.DeclaredType
		case v.String != nil:
			startValue = v.String.Start
			units = v.String.DeclaredType
		}

		sb.WriteString(fmt.Sprintf("  %-18s %-10s %-12s %-8s %s\n", name, v.Causality, startValue, units, v.Description))

	}

	return sb.String(), nil

}

/*
Dump the info for an FMU

Parameters:

	filename     filename of the FMU
	causalities  the causalities of the variables to include
*/
func Dump(filename string, causalities []string) error {
	info, err := FmuInfo(filename, causalities)
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}

func FmiInfo(filename string) (string, []string, error) {
	md, err := ReadModelDescription(filename, nil)
	if err != nil {
		return "", nil, err
	}

	fmiTypes := make([]string, 0)
	if md.ModelExchange != nil {
		fmiTypes = append(fmiTypes, "Model Exchange")
	}

	if md.CoSimulation != nil {
		fmiTypes = append(fmiTypes, "Co-Simulation")
	}

	return md.FmiVersion, fmiTypes, nil
}
