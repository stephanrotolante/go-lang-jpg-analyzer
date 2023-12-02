package main

import (
	"fmt"
	"os"
)

func QuantizationTable(file *os.File, qMap map[int][][]byte) {
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

	for i := 0; i < 8; i++ {
		array := make([]byte, 8)
		for j := 0; j < 8; j++ {
			array[j] = tableData[i*8+j]
		}
		qMap[int(tableType)] = append(qMap[int(tableType)], array)
	}

	fmt.Println("")
}
