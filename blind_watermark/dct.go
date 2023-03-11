package blind_watermark

import (
	"math"
)

func dct(x []float64) []float64 {
	ret := make([]float64, len(x))
	for k, _ := range x {
		for n, _ := range x {
			q := math.Sqrt(1 / float64(len(x)))
			if k != 0 {
				q = math.Sqrt(2 / float64(len(x)))
			}
			ret[k] += q * x[n] * math.Cos(math.Pi*(float64(n)+0.5)*float64(k)/float64(len(x)))
		}
	}
	return ret
}

func idct(xk []float64) []float64 {
	ret := make([]float64, len(xk))
	for n, _ := range xk {
		for k, _ := range xk {
			q := math.Sqrt(1 / float64(len(xk)))
			if k != 0 {
				q = math.Sqrt(2 / float64(len(xk)))
			}
			ret[n] += q * xk[k] * math.Cos(math.Pi*(float64(n)+0.5)*float64(k)/float64(len(xk)))
		}
	}
	return ret
}

func Dct2(x [][]float64) [][]float64 {
	for i, v := range x {
		x[i] = dct(v)
	}
	x = switchRowAndColumns(x)
	for i, v := range x {
		x[i] = dct(v)
	}
	x = switchRowAndColumns(x)
	return x
}

func Idct2(x [][]float64) [][]float64 {
	x = switchRowAndColumns(x)
	for i, v := range x {
		x[i] = idct(v)
	}
	x = switchRowAndColumns(x)
	for i, v := range x {
		x[i] = idct(v)
	}
	return x
}
