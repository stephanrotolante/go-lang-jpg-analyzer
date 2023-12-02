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

func main() {

	qMap := make(map[int][][]byte)

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
			fmt.Println(err)
			panic(err)
		}

		switch BigEUint16(currentBytes[0], currentBytes[1]) {
		// 0xFFE0 Application Header Segnment
		case 65504:
			fmt.Println("APP0")
			ApplicationHeader(file)

		// 0xFFDB Quantization Table Segnment
		case 65499:
			fmt.Println("QT")
			QuantizationTable(file, qMap)

		// 0xFFC0 Start of Frame
		case 65472:
			fmt.Println("SOF")
			StartOfFrame(file)

		// 0xFFC4 Huffman Table
		case 65476:
			fmt.Println("HUF")
			HuffmanTable(file)

		// 0xFFDA Start of Scan
		case 65498:
			fmt.Println("SOC")
			StartOfScan(file)

			return
		// 0xFFDA End of Image
		case 65497:
			fmt.Println("EOI")
			EndOfImage(file)
		default:
			// fmt.Println("Exit")

		}

	}

}
