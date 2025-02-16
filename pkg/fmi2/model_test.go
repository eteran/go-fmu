package fmi2_test

import (
	"go-fmu/pkg/fmi2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadModelRectifier(t *testing.T) {

	const filename = "../../examples/Rectifier.fmu"
	const delta = 1e-4

	md, err := fmi2.ReadModelDescription(filename, nil)

	require.NoError(t, err)
	require.Equal(t, "2.0", md.FmiVersion)
	require.Equal(t, "Rectifier", md.ModelName)
	require.Equal(t, "Model Rectifier", md.Description)
	require.Equal(t, uint32(6), md.NumberOfEventIndicators)
	require.Equal(t, "2017-01-19T18:38:22Z", md.GenerationDateAndTime)
	require.Equal(t, "structured", md.VariableNamingConvention)
	require.Equal(t, "MapleSim (1196527/1196706/1196706)", md.GenerationTool)
	require.Nil(t, md.ModelExchange)
	require.NotNil(t, md.CoSimulation)
	require.NotNil(t, md.DefaultExperiment)
	require.InDelta(t, 0.0, *(md.DefaultExperiment.StartTime), delta)
	require.InDelta(t, 0.1, *(md.DefaultExperiment.StopTime), delta)
	require.Nil(t, md.DefaultExperiment.Tolerance)
	require.InDelta(t, 1e-7, *(md.DefaultExperiment.StepSize), delta)
	require.NotNil(t, md.ModelVariables)
	require.NotNil(t, md.ModelVariables.ScalarVariable)
	require.Len(t, md.ModelVariables.ScalarVariable, 63)
}

func TestReadModelBall(t *testing.T) {

	const filename = "../../examples/Ball.fmu"
	const delta = 1e-4

	md, err := fmi2.ReadModelDescription(filename, nil)

	require.NoError(t, err)
	require.Equal(t, "2.0", md.FmiVersion)
	require.Equal(t, "Ball", md.ModelName)
	require.Equal(t, "", md.Description)
	require.Equal(t, uint32(2), md.NumberOfEventIndicators)
	require.Equal(t, "2022-01-13T16:44:31Z", md.GenerationDateAndTime)
	require.Equal(t, "structured", md.VariableNamingConvention)
	require.Equal(t, "Dymola Version 2022x (64-bit), 2021-12-08", md.GenerationTool)
	require.NotNil(t, md.ModelExchange)
	require.Nil(t, md.CoSimulation)
	require.NotNil(t, md.DefaultExperiment)
	require.InDelta(t, 0.0, *(md.DefaultExperiment.StartTime), delta)
	require.InDelta(t, 1.0, *(md.DefaultExperiment.StopTime), delta)
	require.InDelta(t, 0.0, *(md.DefaultExperiment.Tolerance), delta)
	require.Nil(t, md.DefaultExperiment.StepSize)
	require.NotNil(t, md.ModelVariables)
	require.NotNil(t, md.ModelVariables.ScalarVariable)
	require.Len(t, md.ModelVariables.ScalarVariable, 8)
}
