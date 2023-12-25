package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
)

var err error
var n int

var C1DC, C2DC, C3DC int

var C1AC, C2AC, C3AC int

var C1QT, C2QT, C3QT int

var NumberOfComponents int

var Height, Width int

var RawImageData = []byte{}
var ImageData = []byte{}

var C1DCC, C2DCC, C3DCC int

var Q_MAP = make(map[int][]byte)

var HuffmanTables = make(map[int][]HuffmanTable)

func main() {

	var filePath string
	flag.StringVar(&filePath, "file", "", "path to jpg file")
	flag.Parse()

	if fileType := path.Ext(filePath); fileType != ".jpeg" && fileType != ".jpg" {
		panic(errors.New(fmt.Sprintf("file is not correct type %s", fileType)))
	}

	file, err := os.Open(filePath)

	defer file.Close()
	if err != nil {
		fmt.Printf("failed to open file %s\n", filePath)
		panic(err)

	}

	frameStart := make([]byte, 2)

	n, err = ReadFunc(file, frameStart)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if frameStart[0] != 0xff && frameStart[1] != 0xd8 {
		fmt.Println("Corruptted header")
		return
	}

	// Go throught the segments
	for {

		currentBytes := make([]byte, 2)

		_, err = ReadFunc(file, currentBytes)

		if err != nil {

			if len(RawImageData) >= 3 {
				// Odd Number
				var condition1 = RawImageData[len(RawImageData)-3:][0] == 0xFF && RawImageData[len(RawImageData)-3:][1] == 0xD9

				// Even Number
				var condition2 = RawImageData[len(RawImageData)-3:][1] == 0xFF && RawImageData[len(RawImageData)-3:][2] == 0xD9
				if condition1 {
					RawImageData = RawImageData[:len(RawImageData)-3]
					currentBytes[0] = 0xFF
					currentBytes[1] = 0xD9
				}

				if condition2 {
					RawImageData = RawImageData[:len(RawImageData)-2]
					currentBytes[0] = 0xFF
					currentBytes[1] = 0xD9
				}
			} else {

				panic(err)
			}
		}

		switch BigEUint16(currentBytes[0], currentBytes[1]) {
		// 0xFFE0 Application Header Segnment
		case 65504:
			fmt.Println("APP0")
			ExtractApplicationHeader(file)

		// 0xFFDB Quantization Table Segnment
		case 65499:
			fmt.Println("QT")
			ExtractQuantizationTable(file, Q_MAP)

		// 0xFFC0 Start of Frame
		case 65472:
			fmt.Println("SOF")
			ExtractStartOfFrame(file)

		// 0xFFC4 Huffman Table
		case 65476:
			fmt.Println("HUF")
			ExtractReadHuffmanTable(file, HuffmanTables)

		// 0xFFDA Start of Scan
		case 65498:
			fmt.Println("SOC")
			ExtractStartOfScan(file)

		// 0xFFD9 End of Image
		case 65497:
			fmt.Println("EOI")

			ExtractEndOfImage(file)

			os.Exit(0)

		default:

			RawImageData = append(RawImageData, currentBytes...)

		}

	}

}
