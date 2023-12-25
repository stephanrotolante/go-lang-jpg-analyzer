package main

import (
	"fmt"
	"os"
)

func ExtractQuantizationTable(file *os.File, qMap map[int][]byte) {
	segmentLengthBuffer := make([]byte, 2)

	_, err = ReadFunc(file, segmentLengthBuffer)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	segmentLength := BigEUint16(segmentLengthBuffer[0], segmentLengthBuffer[1])

	fmt.Printf("Length %d\n", segmentLength)

	segmentDataBuffer := make([]byte, segmentLength-2)

	n, err = ReadFunc(file, segmentDataBuffer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableType := int(segmentDataBuffer[0])

	fmt.Printf("Luminance/Chrominance %d\n", tableType)

	tableData := segmentDataBuffer[1:]

	array := make([]byte, 64)
	for h := 0; h < 64; h++ {
		array[h] = tableData[h]
	}

	qMap[int(tableType)] = array

	fmt.Printf("Table Data\n")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%02d ", int(array[i*8+j]))

		}
		fmt.Printf("\n")

	}

	fmt.Println("")
}
