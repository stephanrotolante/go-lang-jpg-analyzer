package main

import (
	"fmt"
	"os"
)

func ApplicationHeader(file *os.File) {

	segmentLengthBuffer := make([]byte, 2)

	_, err := ReadFunc(file, segmentLengthBuffer)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	segmentLength := BigEUint16(segmentLengthBuffer[0], segmentLengthBuffer[1])

	fmt.Printf("Length %d\n", segmentLength)

	segmentDataBuffer := make([]byte, segmentLength-2)

	_, err = ReadFunc(file, segmentDataBuffer)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Identifier %08x\\%02x\n", segmentDataBuffer[:4], segmentDataBuffer[4:5])

	fmt.Printf("Version %d.%d\n",
		int(segmentDataBuffer[5:6][0]),
		int(segmentDataBuffer[6:7][0]),
	)
	fmt.Printf("Units %d\n", int(segmentDataBuffer[7:8][0]))

	fmt.Printf("X Units %d\n", BigEUint16(segmentDataBuffer[8:10][0], segmentDataBuffer[8:10][1]))

	fmt.Printf("Y Units %d\n", BigEUint16(segmentDataBuffer[10:12][0], segmentDataBuffer[10:12][1]))

	fmt.Printf("X ThumbNail %d\n", int(segmentDataBuffer[12:13][0]))

	fmt.Printf("Y ThumbNail %d\n", int(segmentDataBuffer[13:14][0]))

	fmt.Println("")

}
