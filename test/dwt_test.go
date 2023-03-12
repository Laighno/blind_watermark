package test

import (
	"blind_watermark/blind_watermark"
	"testing"
)

func Test_Dwt(t *testing.T) {
	d1 := [][]float64{{1, 2, 4, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}}
	ll, lh, hl, hh := blind_watermark.Dwt2(d1)

	t.Log(ll, lh, hl, hh)

	t.Log(blind_watermark.Idwt2(ll, lh, hl, hh))
}
