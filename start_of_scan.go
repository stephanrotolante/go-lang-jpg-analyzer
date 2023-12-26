package main

import (
	"fmt"
	"os"
)

func ExtractStartOfScan(file *os.File) {
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

	NumberOfComponents = int(BigEUint16(0x00, segmentDataBuffer[0]))

	fmt.Printf("Components %d\n", NumberOfComponents)

	for g := 0; g < NumberOfComponents; g++ {
		switch g {
		case 0:
			C1DC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0xF0>>4))
			C1AC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0x0F))

			fmt.Printf("Component %d DC %d AC %d\n", g+1, C1DC, C1AC)
		case 1:
			C2DC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0xF0>>4))
			C2AC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0x0F))
			fmt.Printf("Component %d DC %d AC %d\n", g+1, C2DC, C2AC)
		case 2:
			C3DC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0xF0>>4))
			C3AC = int(BigEUint16(0x00, segmentDataBuffer[2+(g*2)]&0x0F))
			fmt.Printf("Component %d DC %d AC %d\n", g+1, C3DC, C3AC)
		}

	}

	fmt.Println("")
}
