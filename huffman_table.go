package main

import (
	"fmt"
	"os"
)

func HuffmanTable(file *os.File) {
	segmentLengthBuffer := make([]byte, 2)

	_, err = ReadFunc(file, segmentLengthBuffer)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	segmentLength := BigEUint16(segmentLengthBuffer[0], segmentLengthBuffer[1])

	fmt.Printf("Length %d\n", segmentLength)

	segmentDataBuffer := make([]byte, 1)

	n, err = ReadFunc(file, segmentDataBuffer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	class := int((segmentDataBuffer[0] & 0xF0) >> 4)

	if class == 0 {
		fmt.Printf("Class %s\n", "DC")
	} else {
		fmt.Printf("Class %s\n", "AC")
	}

	fmt.Printf("Destination %d\n", int(segmentDataBuffer[0]&0x0F))

	frquencyCountBuffer := make([]byte, 16)

	n, err = ReadFunc(file, frquencyCountBuffer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	symbolBuffer := make([]byte, segmentLength-2-1-16)

	n, err = ReadFunc(file, symbolBuffer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("%d\n", len(frquencyCountBuffer))
	fmt.Printf("%d\n", len(symbolBuffer))

	offset := 0
	for i := 0; i < 16; i++ {
		freq := int(frquencyCountBuffer[i])

		// fmt.Printf("%d %d\n", i+1, freq)
		fmt.Printf("%02d (%02d): ", i+1, freq)

		for j := 0; j < int(freq); j++ {
			fmt.Printf("%02x ", symbolBuffer[offset+j])
		}

		offset += int(freq)
		fmt.Println("")

	}

	fmt.Println("")
}
