package fmi2

import "fmt"

func (md *ModelDescription) SimulationType() (Type, error) {
	switch {
	case md.ModelExchange != nil:
		return ModelExchangeType, nil
	case md.CoSimulation != nil:
		return CoSimulationType, nil
	default:
		return 0, fmt.Errorf("model must have one of ModelExchange or CoSimulation defined")
	}
}

func (md *ModelDescription) Validate() error {
	// TODO(eteran): validate against XSD
	return nil
}

func (md *ModelDescription) ValidateVariableNames() error {
	// TODO(eteran): validate variable names
	return nil
}

func (md *ModelDescription) ValidateStructure() error {
	// TODO(eteran): validate model structure
	return nil
}

type ModelDescription struct {
	FmiVersion               string             `xml:"fmiVersion,attr"`
	ModelName                string             `xml:"modelName,attr"`
	Guid                     string             `xml:"guid,attr"`
	Description              string             `xml:"description,attr,omitempty"`
	Author                   string             `xml:"author,attr,omitempty"`
	Version                  string             `xml:"version,attr,omitempty"`
	Copyright                string             `xml:"copyright,attr,omitempty"`
	License                  string             `xml:"license,attr,omitempty"`
	GenerationTool           string             `xml:"generationTool,attr,omitempty"`
	GenerationDateAndTime    string             `xml:"generationDateAndTime,attr,omitempty"`
	VariableNamingConvention string             `xml:"variableNamingConvention,attr,omitempty"`
	NumberOfEventIndicators  uint32             `xml:"numberOfEventIndicators,attr,omitempty"`
	ModelExchange            *ModelExchange     `xml:"ModelExchange"`
	CoSimulation             *CoSimulation      `xml:"CoSimulation"`
	UnitDefinitions          []UnitDefinitions  `xml:"UnitDefinitions"`
	TypeDefinitions          []TypeDefinitions  `xml:"TypeDefinitions"`
	LogCategories            []LogCategories    `xml:"LogCategories"`
	DefaultExperiment        *DefaultExperiment `xml:"DefaultExperiment"`
	VendorAnnotations        []Annotation       `xml:"VendorAnnotations"`
	ModelVariables           *ModelVariables    `xml:"ModelVariables"`
	FmiModelDescription      string             `xml:"fmiModelDescription"`
}

type ModelExchange struct {
	ModelIdentifier                     string       `xml:"modelIdentifier,attr"`
	NeedsExecutionTool                  bool         `xml:"needsExecutionTool,attr,omitempty"`
	CompletedIntegratorStepNotNeeded    bool         `xml:"completedIntegratorStepNotNeeded,attr,omitempty"`
	CanBeInstantiatedOnlyOncePerProcess bool         `xml:"canBeInstantiatedOnlyOncePerProcess,attr,omitempty"`
	CanNotUseMemoryManagementFunctions  bool         `xml:"canNotUseMemoryManagementFunctions,attr,omitempty"`
	CanGetAndSetFMUstate                bool         `xml:"canGetAndSetFMUstate,attr,omitempty"`
	CanSerializeFMUstate                bool         `xml:"canSerializeFMUstate,attr,omitempty"`
	ProvidesDirectionalDerivative       bool         `xml:"providesDirectionalDerivative,attr,omitempty"`
	SourceFiles                         *SourceFiles `xml:"SourceFiles"`
}

type CoSimulation struct {
	ModelIdentifier                        string       `xml:"modelIdentifier,attr"`
	NeedsExecutionTool                     bool         `xml:"needsExecutionTool,attr,omitempty"`
	CanHandleVariableCommunicationStepSize bool         `xml:"canHandleVariableCommunicationStepSize,attr,omitempty"`
	CanInterpolateInputs                   bool         `xml:"canInterpolateInputs,attr,omitempty"`
	MaxOutputDerivativeOrder               uint32       `xml:"maxOutputDerivativeOrder,attr,omitempty"`
	CanRunAsynchronuously                  bool         `xml:"canRunAsynchronuously,attr,omitempty"`
	CanBeInstantiatedOnlyOncePerProcess    bool         `xml:"canBeInstantiatedOnlyOncePerProcess,attr,omitempty"`
	CanNotUseMemoryManagementFunctions     bool         `xml:"canNotUseMemoryManagementFunctions,attr,omitempty"`
	CanGetAndSetFMUstate                   bool         `xml:"canGetAndSetFMUstate,attr,omitempty"`
	CanSerializeFMUstate                   bool         `xml:"canSerializeFMUstate,attr,omitempty"`
	ProvidesDirectionalDerivative          bool         `xml:"providesDirectionalDerivative,attr,omitempty"`
	SourceFiles                            *SourceFiles `xml:"SourceFiles"`
}

type SourceFiles struct {
	File []File `xml:"File"`
}

type File struct {
	Name string `xml:"name,attr"`
}

type UnitDefinitions struct {
	Unit []Unit `xml:"Unit"`
}

type TypeDefinitions struct {
	SimpleType []SimpleType `xml:"SimpleType"`
}

type Category struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr,omitempty"`
}

type LogCategories struct {
	Category []Category `xml:"Category"`
}

type DefaultExperiment struct {
	StartTime *float64 `xml:"startTime,attr,omitempty"`
	StopTime  *float64 `xml:"stopTime,attr,omitempty"`
	Tolerance *float64 `xml:"tolerance,attr,omitempty"`
	StepSize  *float64 `xml:"stepSize,attr,omitempty"`
}

type ModelVariables struct {
	ScalarVariable []ScalarVariable `xml:"ScalarVariable"`
}

type Unknown struct {
	Index            uint32 `xml:"index,attr"`
	Dependencies     uint32 `xml:"dependencies,attr,omitempty"`
	DependenciesKind string `xml:"dependenciesKind,attr,omitempty"`
}

type InitialUnknowns struct {
	Unknown []Unknown `xml:"Unknown"`
}

type ModelStructure struct {
	Outputs         []VariableDependency `xml:"Outputs"`
	Derivatives     []VariableDependency `xml:"Derivatives"`
	InitialUnknowns []InitialUnknowns    `xml:"InitialUnknowns"`
}

type VariableDependency struct {
	Unknown []Unknown `xml:"Unknown"`
}

type Real struct {
	RealAttributes
	DeclaredType string  `xml:"declaredType,attr,omitempty"`
	Start        float64 `xml:"start,attr,omitempty"`
	Derivative   uint32  `xml:"derivative,attr,omitempty"`
	Reinit       bool    `xml:"reinit,attr,omitempty"`
}

type Integer struct {
	IntegerAttributes
	DeclaredType string `xml:"declaredType,attr,omitempty"`
	Start        int    `xml:"start,attr,omitempty"`
}

type Boolean struct {
	DeclaredType string `xml:"declaredType,attr,omitempty"`
	Start        bool   `xml:"start,attr,omitempty"`
}

type String struct {
	DeclaredType string `xml:"declaredType,attr,omitempty"`
	Start        string `xml:"start,attr,omitempty"`
}

type Enumeration struct {
	DeclaredType string `xml:"declaredType,attr"`
	Quantity     string `xml:"quantity,attr,omitempty"`
	Min          int    `xml:"min,attr,omitempty"`
	Max          int    `xml:"max,attr,omitempty"`
	Start        int    `xml:"start,attr,omitempty"`
}

type RealType struct {
	FmiRealAttributes *RealAttributes
}

type IntegerType struct {
	FmiIntegerAttributes *IntegerAttributes
}

// RealAttributes is Set to true, e.g., for crank angle. If true and variable is a state, relative tolerance should be zero on this variable.
type RealAttributes struct {
	Quantity         string  `xml:"quantity,attr,omitempty"`
	Unit             string  `xml:"unit,attr,omitempty"`
	DisplayUnit      string  `xml:"displayUnit,attr,omitempty"`
	RelativeQuantity bool    `xml:"relativeQuantity,attr,omitempty"`
	Min              float64 `xml:"min,attr,omitempty"`
	Max              float64 `xml:"max,attr,omitempty"`
	Nominal          float64 `xml:"nominal,attr,omitempty"`
	Unbounded        bool    `xml:"unbounded,attr,omitempty"`
}

// IntegerAttributes is max >= min required
type IntegerAttributes struct {
	Quantity string `xml:"quantity,attr,omitempty"`
	Min      int    `xml:"min,attr,omitempty"`
	Max      int    `xml:"max,attr,omitempty"`
}

type ScalarVariable struct {
	Name                               string       `xml:"name,attr"`
	ValueReference                     uint32       `xml:"valueReference,attr"`
	Description                        string       `xml:"description,attr,omitempty"`
	Causality                          string       `xml:"causality,attr,omitempty"`
	Variability                        string       `xml:"variability,attr,omitempty"`
	Initial                            string       `xml:"initial,attr,omitempty"`
	CanHandleMultipleSetPerTimeInstant bool         `xml:"canHandleMultipleSetPerTimeInstant,attr,omitempty"`
	Real                               *Real        `xml:"Real"`
	Integer                            *Integer     `xml:"Integer"`
	Boolean                            *Boolean     `xml:"Boolean"`
	String                             *String      `xml:"String"`
	Enumeration                        *Enumeration `xml:"Enumeration"`
	Annotations                        []Annotation `xml:"Annotations"`
}

type Item struct {
	Name        string `xml:"name,attr"`
	Value       int    `xml:"value,attr"`
	Description string `xml:"description,attr,omitempty"`
}

type EnumerationType struct {
	Quantity string `xml:"quantity,attr,omitempty"`
	Item     *Item  `xml:"Item"`
}

// SimpleType is Type attributes of a scalar variable
type SimpleType struct {
	Name        string           `xml:"name,attr"`
	Description string           `xml:"description,attr,omitempty"`
	Real        *RealType        `xml:"Real"`
	Integer     *IntegerType     `xml:"Integer"`
	Boolean     *Boolean         `xml:"Boolean"`
	String      *String          `xml:"String"`
	Enumeration *EnumerationType `xml:"Enumeration"`
}

type BaseUnit struct {
	Kg     int     `xml:"kg,attr,omitempty"`
	M      int     `xml:"m,attr,omitempty"`
	S      int     `xml:"s,attr,omitempty"`
	A      int     `xml:"A,attr,omitempty"`
	K      int     `xml:"K,attr,omitempty"`
	Mol    int     `xml:"mol,attr,omitempty"`
	Cd     int     `xml:"cd,attr,omitempty"`
	Rad    int     `xml:"rad,attr,omitempty"`
	Factor float64 `xml:"factor,attr,omitempty"`
	Offset float64 `xml:"offset,attr,omitempty"`
}

type DisplayUnit struct {
	Name   string  `xml:"name,attr"`
	Factor float64 `xml:"factor,attr,omitempty"`
	Offset float64 `xml:"offset,attr,omitempty"`
}

// Unit is Unit definition (with respect to SI base units) and default display units
type Unit struct {
	Name        string       `xml:"name,attr"`
	BaseUnit    *BaseUnit    `xml:"BaseUnit"`
	DisplayUnit *DisplayUnit `xml:"DisplayUnit"`
}

type Tool struct {
	Name string `xml:"name,attr"`
}

type Annotation struct {
	Tool *Tool `xml:"Tool"`
}
