package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const ZOOM = .25

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

	_, err = writer.WriteString(fmt.Sprintf("%d:%d\n", int(float64(Height)*ZOOM), int(float64(Width)*ZOOM)))

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
					var r, g, b, x1, x2, y1, y2 int

					r, g, b = ColorConvert(finalOutput[0][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					x1 = int(float64(x*8+xx) * ZOOM)
					y1 = int(float64(y*8+yy) * ZOOM)
					x2 = int(float64(x*8+(xx+1)) * ZOOM)
					y2 = int(float64(y*8+(yy+1)) * ZOOM)

					_, err = writer.WriteString(fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d\n", x1, y1, x2, y2, r, g, b))
					if err != nil {
						log.Fatal(err)
					}

					if C1XS == 1 && C1YS == 1 {
						continue
					}

					r, g, b = ColorConvert(finalOutput[3][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					x1 = int(float64((x+1)*8+xx) * ZOOM)
					y1 = int(float64(y*8+yy) * ZOOM)
					x2 = int(float64((x+1)*8+(xx+1)) * ZOOM)
					y2 = int(float64(y*8+(yy+1)) * ZOOM)

					_, err = writer.WriteString(fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d\n", x1, y1, x2, y2, r, g, b))
					if err != nil {
						log.Fatal(err)
					}

					r, g, b = ColorConvert(finalOutput[4][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					x1 = int(float64(x*8+xx) * ZOOM)
					y1 = int(float64((y+1)*8+yy) * ZOOM)
					x2 = int(float64(x*8+(xx+1)) * ZOOM)
					y2 = int(float64((y+1)*8+(yy+1)) * ZOOM)

					_, err = writer.WriteString(fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d\n", x1, y1, x2, y2, r, g, b))
					if err != nil {
						log.Fatal(err)
					}

					r, g, b = ColorConvert(finalOutput[5][8*yy+xx], finalOutput[1][8*yy+xx], finalOutput[2][8*yy+xx])

					x1 = int(float64((x+1)*8+xx) * ZOOM)
					y1 = int(float64((y+1)*8+yy) * ZOOM)
					x2 = int(float64((x+1)*8+(xx+1)) * ZOOM)
					y2 = int(float64((y+1)*8+(yy+1)) * ZOOM)

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
