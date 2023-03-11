package test

import (
	"blind_watermark/blind_watermark"
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"
)

func Test_AddWatermark(t *testing.T) {
	imageByte, err := os.ReadFile("./2.png")
	if err != nil {
		t.Log(err)
		return
	}
	img, _, err := image.Decode(bytes.NewBuffer(imageByte))
	if err != nil {
		t.Log(err)
		return
	}

	wmImg, _ := blind_watermark.AddWatermark(img, bytes.NewBufferString("不是昨天说的，每个tab分开记录浏览记录，来达到分开去重的效果吗").Bytes())

	t.Log(blind_watermark.ExtractWaterMask(wmImg))

	f, err := os.Create("./result.png")
	png.Encode(f, wmImg)
}

func Test_ExtractWaterMark(t *testing.T) {
	imageByte, err := os.ReadFile("./img.png")
	if err != nil {
		t.Log(err)
		return
	}
	img, _, err := image.Decode(bytes.NewBuffer(imageByte))
	if err != nil {
		t.Log(err)
		return
	}

	wm, err := blind_watermark.ExtractWaterMask(img)

	t.Log(bytes.NewBuffer(wm).String())
}
