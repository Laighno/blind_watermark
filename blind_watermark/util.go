package blind_watermark

const (
	stride = 4
)

//切分block,二维转四维
func SwitchToBlocks(src [][]float64) (blocks [][][][]float64) {

	rLen := len(src)
	cLen := len(src[0])
	for i := 0; i*stride < rLen; i++ {
		blockRow := make([][][]float64, 0)
		for j := 0; j*stride < cLen; j++ {
			block := [4][]float64{}

			for m := 0; m < stride; m++ {
				rowInBlock := [4]float64{}
				for n := 0; n < stride; n++ {
					rowInBlock[n] = src[i*stride+m][j*stride+n]
				}
				block[m] = rowInBlock[:]
			}
			blockRow = append(blockRow, block[:])
		}
		blocks = append(blocks, blockRow)
	}
	return blocks
}

//还原数据
func RestoreSourceData(blocks [][][][]float64) (src [][]float64) {

	stride := 4
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


func switchRowAndColumns(data [][]float64)  [][]float64{
	ret := make([][]float64, len(data))
	for j:=0;j<len(data);j++{
		ret[j] = make([]float64, len(data[j]))
	}


	for i:=0;i<len(data);i++{
		for j:=0;j<len(data[i]);j++{
			ret[j][i] = data[i][j]
		}
	}
	return ret
}