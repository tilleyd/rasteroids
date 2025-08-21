package math

import (
	m "math"
	"math/rand"
)

func Cosf(v float32) float32 {
	return float32(m.Cos(float64(v)))
}

func Sinf(v float32) float32 {
	return float32(m.Sin(float64(v)))
}

func RandRangef(lower float32, upper float32) float32 {
	return lower + rand.Float32()*(upper-lower)
}

func RandRangei(lower int, upper int) int {
	return lower + rand.Intn(upper-lower)
}

func RandBool() bool {
	return rand.Int()&0b1 == 1
}
