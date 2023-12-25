package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

type HuffmanTable = map[int][][]byte

const HUFF_CODE = 0
const HUFF_SYM = 1

const DC_TABLE = 0
const AC_TABLE = 1

func ReadFunc(file *os.File, buffer []byte) (int, error) {
	n, err := file.Read(buffer)

	if n == 0 {
		return 0, errors.New("No Bytes Read")
	}

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return n, nil
}

func BigEUint16(arg1, arg2 byte) int {

	return int(binary.BigEndian.Uint16([]byte{
		arg1,
		arg2,
	}))
}

func Clamp(value int) int {
	if value > 255 {
		return 255
	}

	if value < 0 {
		return 0
	}

	return value

}

func ColorConvert(Y, Cb, Cr int) (int, int, int) {

	R := int(float64(Y)+(float64(Cr)*1.402)) + 128
	G := int(float64(Y)-(0.344*float64(Cb))-(0.714*float64(Cr))) + 128
	B := int(float64(Y)+(1.772*float64(Cb))) + 128

	return Clamp(R), Clamp(G), Clamp(B)
}

func GetBit() byte {

	var b = (ImageData[0] & 0x01)

	ImageData = ImageData[1:]

	return b
}

func ExtractCoefficient(currentCoeff, coeffBitLength int) int {

	if currentCoeff < (1 << (coeffBitLength - 1)) {
		return currentCoeff - ((1 << coeffBitLength) - 1)
	}

	return currentCoeff
}

func GetQuantTable(component int) int {

	switch component {
	case 1:
		return C1QT
	case 2:
		return C2QT
	case 3:
		return C3QT
	}

	return -1

}

func GetAC(component int) (int, error) {
	switch component {
	case 1:
		return C1AC, nil
	case 2:
		return C2AC, nil
	case 3:
		return C3AC, nil
	}

	return -1, errors.New("Unable to find AC table index")
}

func GetDC(component int) (int, error) {
	switch component {
	case 1:
		return C1DC, nil
	case 2:
		return C2DC, nil
	case 3:
		return C3DC, nil
	}

	return -1, errors.New("Unable to find DC table index")
}

func AddDCC(component, newValue int) int {
	switch component {
	case 1:
		C1DCC = C1DCC + newValue

		return C1DCC
	case 2:
		C2DCC = C2DCC + newValue
		return C2DCC
	case 3:
		C3DCC = C3DCC + newValue
		return C3DCC
	}

	return 0

}

var ZigZag = []int{
	0, 1, 5, 6, 14, 15, 27, 28,
	2, 4, 7, 13, 16, 26, 29, 42,
	3, 8, 12, 17, 25, 30, 41, 43,
	9, 11, 18, 24, 31, 40, 44, 53,
	10, 19, 23, 32, 39, 45, 52, 54,
	20, 22, 33, 38, 46, 51, 55, 60,
	21, 34, 37, 47, 50, 56, 59, 61,
	35, 36, 48, 49, 57, 58, 62, 63,
}

func CreateIDCTTable() []float64 {
	var table []float64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			var c = (float64(2.0*j+1.0) * float64(i) * math.Pi)
			var n = 1.0
			if j == 0 {
				n = 1.0 / math.Sqrt(2.0)
			}
			table = append(table, n*math.Cos(c/16.0))
		}
	}

	return table
}

func FuncOfA(input int) float64 {
	if input == 0 {
		return 1 / math.Sqrt(8)
	}

	return math.Sqrt(2 / 8)
}

func InverseDCT(arg1, arg2 int) float64 {

	return math.Cos(float64((float64(2.0*float64(arg1)+1.0) * float64(arg2) * math.Pi) / 16))
}
