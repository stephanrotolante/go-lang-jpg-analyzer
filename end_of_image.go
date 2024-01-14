package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

var t = CreateIDCTTable()

func ExtractEndOfImage(file *os.File) {

	// Really stupid edgecase
	for len(RawImageData) > 0 {

		if RawImageData[0] == 0xFF {
			if RawImageData[1] == 0x00 {
				for i := 7; i >= 0; i-- {
					bit := (RawImageData[0] >> i) & 0x01
					ImageData = append(ImageData, bit)
				}
				RawImageData = RawImageData[2:]
			}

		} else {
			for i := 7; i >= 0; i-- {
				bit := (RawImageData[0] >> i) & 0x01
				ImageData = append(ImageData, bit)
			}
			RawImageData = RawImageData[1:]
		}

	}

	file, err := os.OpenFile(OutFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	var PaddingSize = Width % 4
	var Size = 14 + 12 + Height*Width*3 + PaddingSize*Height

	writer.WriteByte(byte('B'))
	writer.WriteByte(byte('M'))

	var b1, b2, b3, b4 byte

	b1, b2, b3, b4 = IntToLittleEdian(Size)
	writer.WriteByte(b1)
	writer.WriteByte(b2)
	writer.WriteByte(b3)
	writer.WriteByte(b4)

	b1, b2, b3, b4 = IntToLittleEdian(0)
	writer.WriteByte(b1)
	writer.WriteByte(b2)
	writer.WriteByte(b3)
	writer.WriteByte(b4)

	b1, b2, b3, b4 = IntToLittleEdian(0x1A)
	writer.WriteByte(b1)
	writer.WriteByte(b2)
	writer.WriteByte(b3)
	writer.WriteByte(b4)

	b1, b2, b3, b4 = IntToLittleEdian(12)
	writer.WriteByte(b1)
	writer.WriteByte(b2)
	writer.WriteByte(b3)
	writer.WriteByte(b4)

	b1, b2 = ShortToLittleEdian(Width)
	writer.WriteByte(b1)
	writer.WriteByte(b2)

	b1, b2 = ShortToLittleEdian(Height)
	writer.WriteByte(b1)
	writer.WriteByte(b2)

	b1, b2 = ShortToLittleEdian(1)
	writer.WriteByte(b1)
	writer.WriteByte(b2)

	b1, b2 = ShortToLittleEdian(24)
	writer.WriteByte(b1)
	writer.WriteByte(b2)

	var Pixels = make([][][]int, Height+15)

	for j := 0; j < len(Pixels); j++ {
		Pixels[j] = make([][]int, Width+15)
	}

	for y := 0; y < int(math.Floor(float64((Height+7)/8))); y += C1YS {
		for x := 0; x < int(math.Floor(float64((Width+7)/8))); x += C1XS {

			var finalOutput [6][64]int

			for component := 1; component <= NumberOfComponents; component++ {

				var bits = uint16(0x0000)

				dcTableIndex, err := GetDC(component)

				if err != nil {
					panic(err)
				}

				YS, err := GetYS(component)

				if err != nil {
					panic(err)
				}

				XS, err := GetXS(component)

				if err != nil {
					panic(err)
				}

				for ys := 0; ys < YS; ys++ {
					for xs := 0; xs < XS; xs++ {

						for bitLength := 1; bitLength < 17; bitLength++ {
							bits = (bits << 1) + uint16(GetBit())

							for _, MAPPED_SYM_CODE := range HuffmanTables[DC_TABLE][dcTableIndex][bitLength] {

								if MAPPED_SYM_CODE[HUFF_CODE] == bits {
									bits = 0x0000
									bitLength += 17
									var coeffLength = int(MAPPED_SYM_CODE[HUFF_SYM] & 0x0F)

									var coeffByte = uint16(0x0000)
									for coeffCount := 0; coeffCount < coeffLength; coeffCount++ {
										coeffByte = (coeffByte << 1) + uint16(GetBit())
									}

									var coeff = int(coeffByte)

									if coeffLength != 0 {
										coeff = AddDCC(component, ExtractCoefficient(coeff, coeffLength))
									} else {
										coeff = AddDCC(component, 0)
									}

									finalOutput[GetComponent(component-1, ys, xs)][0] = coeff * int(Q_MAP[GetQuantTable(component)][0])

									break
								}
							}
						}

						for coeffIndex := 1; coeffIndex < 64; coeffIndex += 0 {

							acTableIndex, err := GetAC(component)

							if err != nil {
								panic(err)
							}

							for bitLength := 1; bitLength < 17; bitLength++ {

								bits = (bits << 1) + uint16(GetBit())

								for _, MAPPED_SYM_CODE := range HuffmanTables[AC_TABLE][acTableIndex][bitLength] {
									if MAPPED_SYM_CODE[HUFF_CODE] == bits {
										bitLength += 17
										bits = 0x0000

										var coeffLength = int(MAPPED_SYM_CODE[HUFF_SYM] & 0x0F)
										var numberOfZeros = int((MAPPED_SYM_CODE[HUFF_SYM] >> 4))

										if MAPPED_SYM_CODE[HUFF_SYM] == 0x00 {
											// Whole MCU is zero
											coeffIndex = 64
											break

										}

										if coeffIndex+numberOfZeros > 63 {
											coeffIndex = 64
											break
										}

										coeffIndex += numberOfZeros

										var coeffByte = uint16(0x0000)
										for coeffCount := 0; coeffCount < coeffLength; coeffCount++ {
											coeffByte = (coeffByte << 1) + uint16(GetBit())
										}

										var coeff = 0
										if coeffLength != 0 {
											coeff = ExtractCoefficient(int(coeffByte), coeffLength)
										}

										finalOutput[GetComponent(component-1, ys, xs)][ZigZag[coeffIndex]] = coeff * int(Q_MAP[GetQuantTable(component)][coeffIndex])
										coeffIndex += 1
										break
									}
								}

							}
						}
					}
				}

			}

			for c := 1; c <= 6; c++ {
				var coeffList [64]int
				var componentIndex = c - 1
				for i := 0; i < 64; i++ {
					coeffList[i] = finalOutput[componentIndex][i]
				}
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						var sum = 0.00

						for n := 0; n < 8; n++ {
							for m := 0; m < 8; m++ {
								sum += float64(coeffList[n*8+m]) * t[m*8+x] * t[n*8+y]
							}
						}

						finalOutput[componentIndex][y*8+x] = int(sum / 4)
					}
				}
			}
			for yy := 0; yy < 8; yy++ {
				for xx := 0; xx < 8; xx++ {
					var r, g, b, yOffset, xOffset int

					r, g, b = ColorConvert(finalOutput[0][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					yOffset = (y * 8) + yy
					xOffset = (x*8 + xx)
					Pixels[yOffset][xOffset] = []int{
						r, g, b,
					}

					if C1XS == 1 && C1YS == 1 {
						continue
					}

					r, g, b = ColorConvert(finalOutput[3][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					yOffset = (y * 8) + yy
					xOffset = ((x+1)*8 + xx)
					Pixels[yOffset][xOffset] = []int{
						r, g, b,
					}

					r, g, b = ColorConvert(finalOutput[4][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					yOffset = ((y + 1) * 8) + yy
					xOffset = (x*8 + xx)
					Pixels[yOffset][xOffset] = []int{
						r, g, b,
					}

					r, g, b = ColorConvert(finalOutput[5][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					yOffset = ((y + 1) * 8) + yy
					xOffset = ((x+1)*8 + xx)
					Pixels[yOffset][xOffset] = []int{
						r, g, b,
					}

				}

			}

		}
	}

	for h := Height - 1; h >= 0; h-- {

		for w := 0; w < Width; w++ {

			if len(Pixels[h][w]) == 3 {
				writer.WriteByte(byte(Pixels[h][w][2]))
				writer.WriteByte(byte(Pixels[h][w][1]))
				writer.WriteByte(byte(Pixels[h][w][0]))
			} else {
				writer.WriteByte(0x00)
				writer.WriteByte(0x00)
				writer.WriteByte(0x00)
			}
		}
		for i := 0; i < PaddingSize; i++ {
			writer.WriteByte(0x00)
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Left over %d\n", len(ImageData))

}
