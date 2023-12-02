package main

import (
	"fmt"
	"os"
)

func StartOfFrame(file *os.File) {
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

	fmt.Printf("Precision %d\n", BigEUint16(segmentDataBuffer[0], segmentDataBuffer[1]))
	fmt.Printf("Line No %d\n", BigEUint16(segmentDataBuffer[2], segmentDataBuffer[3]))
	fmt.Printf("Samples Per Line %d\n", int(segmentDataBuffer[4]))

	numberOfComponents := int(segmentDataBuffer[5])
	// fmt.Printf("Components %d\n", numberOfComponents)

	for i := 0; i < int(numberOfComponents); i++ {

		index := 6 + i*3

		fmt.Printf("Component %d %dx%d %d\n",
			int(segmentDataBuffer[index]),
			// Mask And Shift
			int((segmentDataBuffer[index+1]&0xF0)>>4),
			int(segmentDataBuffer[index+1]&0x0F),
			int(segmentDataBuffer[index+2]),
		)
	}

	fmt.Println("")
}
