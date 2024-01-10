package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

type HuffmanTable = map[int][][]uint16

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

func BigEUint16(arg1, arg2 byte) uint16 {

	return binary.BigEndian.Uint16([]byte{
		arg1,
		arg2,
	})
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

// I hate this lol
func GetComponent(c, s1, s2 int) int {

	if c == 0 {
		if s1 == 0 && s2 == 0 {
			return c
		} else if s1 == 0 && s2 == 1 {
			return 3
		} else if s1 == 1 && s2 == 0 {
			return 4
		} else if s1 == 1 && s2 == 1 {
			return 5
		}

		return 3
	}

	return c
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

func GetXS(component int) (int, error) {
	switch component {
	case 1:
		return C1XS, nil
	case 2:
		return C2XS, nil
	case 3:
		return C3XS, nil
	}

	return -1, errors.New("Unable to find AC table index")
}

func GetYS(component int) (int, error) {
	switch component {
	case 1:
		return C1YS, nil
	case 2:
		return C2YS, nil
	case 3:
		return C3YS, nil
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
	0, 1, 8, 16, 9, 2, 3, 10,
	17, 24, 32, 25, 18, 11, 4, 5,
	12, 19, 26, 33, 40, 48, 41, 34,
	27, 20, 13, 6, 7, 14, 21, 28,
	35, 42, 49, 56, 57, 50, 43, 36,
	29, 22, 15, 23, 30, 37, 44, 51,
	58, 59, 52, 45, 38, 31, 39, 46,
	53, 60, 61, 54, 47, 55, 62, 63,
}

func CreateIDCTTable() []float64 {
	var table []float64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {

			table = append(table, FuncOfA(i)*InverseDCT(j, i))
		}
	}

	return table
}

func FuncOfA(input int) float64 {
	if input == 0 {
		return 1.0 / math.Sqrt(2.0)
	}

	return 1.0
}

func InverseDCT(arg1, arg2 int) float64 {

	return math.Cos(float64((float64(2.0*float64(arg1)+1.0) * float64(arg2) * math.Pi) / 16))
}

func IntToLittleEdian(v int) (byte, byte, byte, byte) {

	return byte(v >> 0 & 0xFF), byte(v >> 8 & 0xFF), byte(v >> 16 & 0xFF), byte(v >> 24 & 0xFF)
}

func ShortToLittleEdian(v int) (byte, byte) {

	return byte(v >> 0 & 0xFF), byte(v >> 8 & 0xFF)
}
