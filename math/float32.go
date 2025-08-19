package math

import m "math"

func Cosf(v float32) float32 {
	return float32(m.Cos(float64(v)))
}

func Sinf(v float32) float32 {
	return float32(m.Sin(float64(v)))
}
