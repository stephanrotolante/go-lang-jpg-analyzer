package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const ZOOM = 1

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

	file, err := os.OpenFile("color.out", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for y := 0; y < int(math.Floor(float64(Height/8))); y++ {
		for x := 0; x < int(math.Floor(float64(Width/8))); x++ {

			var finalOutput [3][64]int

			for component := 1; component <= NumberOfComponents; component++ {

				var bits = byte(0x00)

				var coeffList [64]int

				dcTableIndex, err := GetDC(component)

				if err != nil {
					panic(err)
				}

				for bitLength := 1; bitLength < 17; bitLength++ {
					bits = (bits << 1) | GetBit()

					for _, MAPPED_SYM_CODE := range HuffmanTables[DC_TABLE][dcTableIndex][bitLength] {

						if MAPPED_SYM_CODE[HUFF_CODE] == bits {
							bits = 0x00
							bitLength += 17
							var coeffLength = int(MAPPED_SYM_CODE[HUFF_SYM] & 0x0F)

							var coeffByte = byte(0x00)
							for coeffCount := 0; coeffCount < coeffLength; coeffCount++ {
								coeffByte = (coeffByte << 1) | GetBit()
							}

							var coeff = int(coeffByte)

							if coeffLength != 0 {
								coeff = AddDCC(component, ExtractCoefficient(coeff, coeffLength))
							}

							coeffList[ZigZag[0]] = coeff * int(Q_MAP[GetQuantTable(component)][0])

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

						bits = (bits << 1) | GetBit()

						for _, MAPPED_SYM_CODE := range HuffmanTables[AC_TABLE][acTableIndex][bitLength] {
							if MAPPED_SYM_CODE[HUFF_CODE] == bits {
								bitLength += 17
								bits = 0x00

								var coeffLength = int(MAPPED_SYM_CODE[HUFF_SYM] & 0x0F)
								var numberOfZeros = int((MAPPED_SYM_CODE[HUFF_SYM] >> 4))

								if MAPPED_SYM_CODE[HUFF_SYM] == 0xF0 {
									coeffIndex += 16
								}

								if MAPPED_SYM_CODE[HUFF_SYM] == 0x00 {
									// Whole MCU is zero
									fmt.Printf("All Zeros %d %02x\n\tAC TABLE %d\n\n", bitLength-17, MAPPED_SYM_CODE[HUFF_SYM], acTableIndex)
									coeffIndex = 64
									break

								}

								if coeffIndex+numberOfZeros > 63 {
									// fmt.Printf("Breakout too high\n")
									break
								}

								if coeffLength == 0 {
									fmt.Printf("Breakout no length\n")
									break
								}

								if coeffLength > 10 {
									// fmt.Printf("Breakout no length\n")
									break
								}

								coeffIndex += numberOfZeros

								var coeffByte = byte(0x00)
								for coeffCount := 0; coeffCount < coeffLength; coeffCount++ {
									coeffByte = (coeffByte << 1) | GetBit()
								}

								if int(coeffByte) == -1 {
									break
								}

								var coeff = ExtractCoefficient(int(coeffByte), coeffLength)

								coeffList[ZigZag[coeffIndex]] = coeff * int(Q_MAP[GetQuantTable(component)][ZigZag[coeffIndex]])
								coeffIndex += 1
								break
							}
						}

					}

				}

				var t = CreateIDCTTable()
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						var sum = 0.00

						for n := 0; n < 8; n++ {
							for m := 0; m < 8; m++ {
								sum += float64(coeffList[n*8+m]) * t[n*8+y] * t[m*8+x]
								// InverseDCT(8, n) * InverseDCT(8, m) * FuncOfA(n) * FuncOfA(m)
							}
						}

						finalOutput[component-1][y*8+x] = int(sum / 4)
					}
				}

			}

			for yy := 0; yy < 8; yy++ {
				for xx := 0; xx < 8; xx++ {
					r, g, b := ColorConvert(finalOutput[0][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					x1 := (x*8 + xx) * ZOOM
					y1 := (y*8 + yy) * ZOOM
					x2 := (x1 + 1) * ZOOM
					y2 := (y1 + 1) * ZOOM

					_, err = writer.WriteString(fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d\n", x1, y1, x2, y2, r, g, b))
					if err != nil {
						log.Fatal(err)
					}
				}

			}

		}
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Left over %d\n", len(ImageData))

}
