package main

import (
	"fmt"
	"os"
)

func ExtractStartOfFrame(file *os.File) {
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

	fmt.Printf("Precision %d\n", int(segmentDataBuffer[0]))

	Height = int(BigEUint16(segmentDataBuffer[1], segmentDataBuffer[2]))
	fmt.Printf("Line No %d\n", Height)

	Width = int(BigEUint16(segmentDataBuffer[3], segmentDataBuffer[4]))
	fmt.Printf("Samples Per Line %d\n", Width)

	numberOfComponents := int(segmentDataBuffer[5])
	// fmt.Printf("Components %d\n", numberOfComponents)

	for i := 0; i < int(numberOfComponents); i++ {

		index := 6 + i*3

		switch i + 1 {
		case 1:
			C1QT = int(segmentDataBuffer[index+2])
			C1XS = int((segmentDataBuffer[index+1] & 0xF0) >> 4)
			C1YS = int(segmentDataBuffer[index+1] & 0x0F)
		case 2:
			C2QT = int(segmentDataBuffer[index+2])
			C2XS = int((segmentDataBuffer[index+1] & 0xF0) >> 4)
			C2YS = int(segmentDataBuffer[index+1] & 0x0F)
		case 3:
			C3QT = int(segmentDataBuffer[index+2])
			C3XS = int((segmentDataBuffer[index+1] & 0xF0) >> 4)
			C3YS = int(segmentDataBuffer[index+1] & 0x0F)
		}

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
