package test

import (
	"blind_watermark/blind_watermark"
	"testing"
)

func Test_Dct(t *testing.T) {
	d1 := [][]float64{{1, 2, 4, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}}
	d2 := blind_watermark.Dct2(d1)

	t.Log(d2)

	t.Log(blind_watermark.Idct2(d2))
}
