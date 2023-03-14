package blind_watermark

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
)

const (
	d = 32
)

func AddWatermark(img image.Image, wm []byte) (image.Image, error) {
	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	pix := make([][]float64, RGBAImage.Bounds().Max.X)
	for i := 0; i < RGBAImage.Bounds().Max.X; i++ {
		row := make([]float64, img.Bounds().Max.Y)
		for j := 0; j < RGBAImage.Bounds().Max.Y; j++ {
			r, _, _, _ := RGBAImage.At(i, j).RGBA()
			row[j] = float64(uint8(r))
		}
		pix[i] = row
	}

	//ll, lh, hl, hh := Dwt2(pix)

	//blocks := SwitchToBlocks(ll)
	//
	//blocks = EmbedWm(blocks, wm)
	//
	////fmt.Println(bytes.NewBuffer(append(bytes.NewBufferString(wm).Bytes(), []byte{0,0}...)).Bytes())
	////fmt.Println(blocks)
	//ll = RestoreSourceData(blocks)
	//
	//pix = Idwt2(ll, lh, hl, hh)
	blocks := SwitchToBlocks(pix)
	blocks = EmbedWm(blocks, wm)
	pix = RestoreSourceData(blocks)

	for i, _ := range pix {
		for j, w := range pix[i] {
			if w > 255 {
				//fmt.Println(i, j, w)
				pix[i][j] = 255
			}
			if w < 0 {
				pix[i][j] = 0
			}
		}
	}

	for i := 0; i < RGBAImage.Bounds().Max.X; i++ {
		for j := 0; j < RGBAImage.Bounds().Max.Y; j++ {
			if i>= len(pix)||j>= len(pix[i]){
				continue
			}
			_, g, b, a := RGBAImage.At(i, j).RGBA()
			RGBAImage.Set(i, j, color.RGBA{
				R: uint8(pix[i][j] + 0.1),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
			//fmt.Println(RGBAImage.At(i,j).RGBA())
		}
	}

	return RGBAImage, nil
}

func ExtractWaterMask(img image.Image) ([]byte, error) {
	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	pix := make([][]float64, RGBAImage.Bounds().Max.X)
	for i := 0; i < RGBAImage.Bounds().Max.X; i++ {
		row := make([]float64, img.Bounds().Max.Y)
		for j := 0; j < RGBAImage.Bounds().Max.Y; j++ {
			r, _, _, _ := RGBAImage.At(i, j).RGBA()
			row[j] = float64(uint8(r))
		}
		pix[i] = row
	}

	//ll, _, _, _ := Dwt2(pix)

	//blocks := SwitchToBlocks(ll)

	wm,err:= ExtractWm(pix)

	return wm, err
}

func EmbedWm(src [][][][]float64, wm []uint8) [][][][]float64 {
	//startSymbol := []byte{1, 1, 1, 1}
	//wm = append(append(startSymbol, wm...), []byte{0, 0, 0, 0}...)
	for i := 0; i < len(src); i++ {
		for j := 0; j < len(src[i]); j++ {
			src[i][j] = embedToBlock(src[i][j],wm)
			//pos := ((i*len(src[i]) + j) / 8) % len(wm)
			//phase := (i*len(src[i]) + j) % 8
			////fmt.Println(i, j, pos, phase)
			////if pos >= len(wm) {
			////	//ExtractWm(src)
			////	return src
			////}
			//
			//bit := (wm[pos] >> phase) & 1
			////fmt.Println(bit)
			////fmt.Println((float64(uint64(src[i][j][0][0]/36)) + (1.0/4) + (1.0/2)*float64(b[i*len(src)+j])), (float64(uint64(src[i][j][0][0]/36)) + (1/4) + (1/2)*float64(b[i*len(src)+j])) * 36)
			////src[i][j][0][0] = (float64(uint64(src[i][j][0][0]/36)) + 1.0/4 + 1.0/2*bit) * 36
			////fmt.Println(src[i][j][0][0])
			//src[i][j] = embedOneBitInBlock(src[i][j], bit)
			//fmt.Println(src[i][j][0][0], bit, extractBitFromBlock(src[i][j]))
			//extractBitFromBlock(src[i][j])
		}
	}
	//fmt.Println(src)
	return src
}

func embedToBlock(block [][]float64, wm []uint8) [][]float64 {
	pos :=0
	phase:=0
	for i:=0;i< len(block);i++{
		for j:=0;j< len(block[i]);j++{
			if i==0|| j==0 || i==len(block)-1 || j==len(block[i])-1{
				block[i][j] = (float64(uint64(block[i][j]/d)) + 1.0/4 + 1.0/2*float64(1)) * d
			}else {
				bit:= uint8(0)
				if pos< len(wm){
					bit = (wm[pos] >> phase) & 1
				}
				phase++
				if phase==8{
					pos++
					phase=0
				}
				block[i][j] = (float64(uint64(block[i][j]/d)) + 1.0/4 + 1.0/2*float64(bit)) * d
			}
		}
	}
	return block
}

func embedOneBitInBlock(block [][]float64, bit uint8) [][]float64 {
	//block = Dct2(block)

	block[0][0] = (float64(uint64(block[0][0]/d)) + 1.0/4 + 1.0/2*float64(bit)) * d
	//fmt.Println(block[0][0], uint64(block[0][0])%36)
	//block = Idct2(block)

	return block
}

func extractBitFromBlock(block [][]float64) uint8 {
	//block = Dct2(block)
	//fmt.Println(block[0][0], uint64(block[0][0])%36)
	//for i:=0;i<stride;i++{
	//	for j:=0;j<stride;j++{
	//		if FloatToUint8(block[i][j])%d == d/4|| FloatToUint8(block[i][j])%d==d*3/4{
	//			return 1
	//		}
	//	}
	//}
	if uint64(block[0][0])%d > d/2 {
		return 1
	}
	return 0
}

func extractFromBlock(block [][]float64) (bool,[]uint8) {
	for i:=0;i<stride;i++{
		if uint64(block[i][0])%d <= d/2 || uint64(block[0][i])%d <= d/2 || uint64(block[stride-1][i])%d <= d/2 ||uint64(block[i][stride-1])%d <= d/2{
			return false, nil
		}
 	}

 	wm := make([]uint8,stride*stride)
 	pos := 0
 	phase:=0
 	for i:=1;i<stride-1;i++{
 		for j:=1;j<stride-1;j++{
 			bit := uint8(0)
 			if uint64(block[i][j])%d > d/2{
 				bit = 1
			}
			wm[pos] |= bit << phase

			phase++
			if phase==8{
				pos++
				phase=0
			}
		}
	}

	for pos>=0{
		if wm[pos]==0{
			pos--
		}else {
			break
		}
	}

	return true,wm[:pos+1]
}

func ExtractWm(src [][]float64) ([]uint8,error) {
	for i := 0; i < len(src); i++ {
		startCnt := 0
		for j := 0; j < len(src[i]); j++ {
			if uint64(src[i][j])%d > d/2{
				startCnt++
			}else {
				startCnt=0
			}
			if startCnt>=stride{
				block := make([][]float64,stride)
				for m:=0;m<stride;m++{
					block[m] = src[i+m][j-stride+1:j+1]
				}
				isWm,wm := extractFromBlock(block)
				if isWm{
					return wm,nil
				}
			}
		}
	}

	return nil,errors.New("could not find watermark")
}
