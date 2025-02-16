package fmi2

import "math"

type Float64IsCloseOption func(*Float64IsCloseOptions)

type Float64IsCloseOptions struct {
	epsilon float64
}

func WithEpsilon(epsilon float64) Float64IsCloseOption {
	return func(o *Float64IsCloseOptions) {
		o.epsilon = epsilon
	}
}

func Float64IsClose(a float64, b float64, opts ...Float64IsCloseOption) bool {

	options := &Float64IsCloseOptions{
		epsilon: 1e-4,
	}

	for _, opt := range opts {
		opt(options)
	}

	if a == b { // Handle the case where a direct comparison works
		return true
	}

	d := math.Abs(a - b)

	// To avoid division by zero
	if b == 0 {
		return d < options.epsilon
	}

	return (d / math.Abs(b)) < options.epsilon
}
