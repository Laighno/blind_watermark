package blind_watermark

import "fmt"

const (
	stride = 4
)

// 切分block,二维转四维
func SwitchToBlocks(src [][]float64) (blocks [][][][]float64) {

	rLen := len(src)
	cLen := len(src[0])
	fmt.Println(rLen, cLen)
	for i := 0; i*stride+stride < rLen; i++ {
		blockRow := make([][][]float64, 0)
		for j := 0; j*stride+stride < cLen; j++ {
			block := [stride][]float64{}

			for m := 0; m < stride; m++ {
				rowInBlock := [stride]float64{}
				for n := 0; n < stride; n++ {
					//fmt.Println(i, j, m, n, i*stride+m, j*stride+n)
					//if j*stride+n>cLen{
					//
					//}
					rowInBlock[n] = src[i*stride+m][j*stride+n]
				}
				block[m] = rowInBlock[:]
			}
			blockRow = append(blockRow, block[:])
		}
		//fmt.Println(i, blockRow)
		blocks = append(blocks, blockRow)
	}
	return blocks
}

// 还原数据
func RestoreSourceData(blocks [][][][]float64) (src [][]float64) {
	rBlockLen := len(blocks)
	cBlockLen := len(blocks[0])

	src = make([][]float64, rBlockLen*stride)
	for i := 0; i < rBlockLen; i++ {
		for m := 0; m < stride; m++ {
			row := make([]float64, cBlockLen*stride)
			for j := 0; j < cBlockLen; j++ {
				for n := 0; n < stride; n++ {
					row[j*stride+n] = blocks[i][j][m][n]
				}
			}
			src[i*stride+m] = row
		}
	}
	return src
}

func switchRowAndColumns(data [][]float64) [][]float64 {
	ret := make([][]float64, len(data[0]))
	for j := 0; j < len(data[0]); j++ {
		ret[j] = make([]float64, len(data))
	}

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			ret[j][i] = data[i][j]
		}
	}
	return ret
}

func FloatToUint8(f float64) uint8 {
	return uint8(f + 0.1)
}
