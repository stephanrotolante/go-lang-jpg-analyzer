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

var C1XS, C2XS, C3XS int

var C1YS, C2YS, C3YS int

var NumberOfComponents int

var Height, Width int

var RawImageData = []byte{}
var ImageData = []byte{}

var C1DCC, C2DCC, C3DCC int

var Q_MAP = make(map[int][]byte)

var HuffmanTables = make(map[int][]HuffmanTable)

var OutFile string

func main() {

	var filePath string
	flag.StringVar(&filePath, "file", "", "path to jpg file")
	flag.StringVar(&OutFile, "output", "", "output file name")
	flag.Parse()

	if fileType := path.Ext(filePath); fileType != ".jpeg" && fileType != ".jpg" {
		panic(errors.New(fmt.Sprintf("file is not correct type %s", fileType)))
	}

	if OutFile == "" {
		panic(errors.New("No output file specified"))
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

				/*
				* End of file is reached. We ran out of bytes to read now we are handling the execption.
				* If 0xFFD9 is not found then we ran into an issue
				 */

				// Final 3 Bytes
				var f3B = RawImageData[len(RawImageData)-3:]

				// Odd Number
				if 0xFFD9 == BigEUint16(f3B[0], f3B[1]) {
					RawImageData = RawImageData[:len(RawImageData)-3]
					currentBytes[0] = 0xFF
					currentBytes[1] = 0xD9
				}
				// Even Number
				if 0xFFD9 == BigEUint16(f3B[1], f3B[2]) {
					RawImageData = RawImageData[:len(RawImageData)-2]
					currentBytes[0] = 0xFF
					currentBytes[1] = 0xD9
				}
			} else {

				panic(err)
			}
		}

		switch BigEUint16(currentBytes[0], currentBytes[1]) {
		// Application Header Segnment
		case 0xFFE0, 0xFFE1, 0xFFE2, 0xFFE3, 0xFFE4, 0xFFE5, 0xFFE6, 0xFFE7, 0xFFE8, 0xFFE9, 0xFFEA, 0xFFEB, 0xFFEC, 0xFFED, 0xFFEE, 0xFFEF:
			fmt.Printf("APP%d\n", int(currentBytes[1]&0x01))
			ExtractApplicationHeader(file)

		// Quantization Table Segnment
		case 0xFFDB:
			fmt.Println("QT")
			ExtractQuantizationTable(file, Q_MAP)

		// Start of Frame
		case 0xFFC0:
			fmt.Println("SOF")
			ExtractStartOfFrame(file)

		// Huffman Table
		case 0xFFC4:
			fmt.Println("HUF")
			ExtractReadHuffmanTable(file, HuffmanTables)

		// Start of Scan
		case 0xFFDA:
			fmt.Println("SOC")
			ExtractStartOfScan(file)

		// End of Image
		case 0xFFD9:
			fmt.Println("EOI")

			ExtractEndOfImage(file)

			os.Exit(0)

		default:

			RawImageData = append(RawImageData, currentBytes...)

		}

	}

}
