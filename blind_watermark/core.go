package blind_watermark

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

func AddWatermark(img image.Image,wm []byte) (image.Image,error){
	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	pix := make([][]float64,RGBAImage.Bounds().Max.X)
	for i:=0;i<RGBAImage.Bounds().Max.X;i++ {
		row := make([]float64,img.Bounds().Max.Y)
		for j:=0;j<RGBAImage.Bounds().Max.Y;j++{
			r,_,_,_ := RGBAImage.At(i,j).RGBA()
			row[j] = float64(uint8(r))
		}
		pix[i] = row
	}

	blocks := SwitchToBlocks(pix)

	blocks = EmbedWm(blocks,wm)

	//fmt.Println(bytes.NewBuffer(append(bytes.NewBufferString(wm).Bytes(), []byte{0,0}...)).Bytes())

	pix = RestoreSourceData(blocks)

	for i:=0;i<RGBAImage.Bounds().Max.X;i++ {
		for j:=0;j<RGBAImage.Bounds().Max.Y;j++{
			_,g,b,a:= RGBAImage.At(i,j).RGBA()
			RGBAImage.Set(i,j,color.RGBA{
				R: uint8(pix[i][j]+0.1),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
			//fmt.Println(RGBAImage.At(i,j).RGBA())
		}
	}

	return RGBAImage,nil
}

func ExtractWaterMask(img image.Image) ([]byte, error) {
	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	pix := make([][]float64,RGBAImage.Bounds().Max.X)
	for i:=0;i<RGBAImage.Bounds().Max.X;i++ {
		row := make([]float64,img.Bounds().Max.Y)
		for j:=0;j<RGBAImage.Bounds().Max.Y;j++{
			r,_,_,_ := RGBAImage.At(i,j).RGBA()
			row[j] = float64(r)
		}
		pix[i] = row
	}

	blocks := SwitchToBlocks(pix)

	wm := ExtractWm(blocks)

	return wm,nil
}

func EmbedWm(src [][][][]float64, wm []uint8) [][][][]float64 {
	wm = append(wm, []byte{0,0,0,0}...)
	for i := 0; i < len(src); i++ {
		for j := 0; j < len(src[i]); j++ {
			pos := (i*len(src[i])+j)/8
			phase := (i*len(src[i])+j)%8
			fmt.Println(i,j,pos,phase)
			if pos >= len(wm) {
				return src
			}

			bit := float64((wm[pos]>>phase)&1)
			//fmt.Println(bit)
			//fmt.Println((float64(uint64(src[i][j][0][0]/36)) + (1.0/4) + (1.0/2)*float64(b[i*len(src)+j])), (float64(uint64(src[i][j][0][0]/36)) + (1/4) + (1/2)*float64(b[i*len(src)+j])) * 36)
			src[i][j][0][0] = (float64(uint64(src[i][j][0][0]/36)) + 1.0/4 + 1.0/2*bit) * 36
		}
	}
	//fmt.Println(src)
	return src
}

func ExtractWm(src [][][][]float64) []uint8 {
	wm := make([]uint8, len(src)*len(src[0]))
	endCnt:=0
	for i := 0; i < len(src); i++ {
		for j := 0; j < len(src[i]); j++ {
			pos := (i*len(src[i])+j)/8
			phase := (i*len(src[i])+j)%8

			if uint64(src[i][j][0][0])%36> 18{
				wm[pos] |= 1<<phase
			}

			if phase==7{
				if wm[pos]==0{
					endCnt++
				}else {
					endCnt=0
				}
				if endCnt>=4{
					return wm[:pos-3]
				}
			}
		}
	}

	return wm
}