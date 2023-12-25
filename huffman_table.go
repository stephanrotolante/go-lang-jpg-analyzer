package main

import (
	"fmt"
	"os"
)

func ExtractReadHuffmanTable(file *os.File, huffmanTables map[int][]HuffmanTable) {
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

	var tableId = int(segmentDataBuffer[0] & 0x0F)

	fmt.Printf("Table ID %d\n", tableId)

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

	fmt.Printf("Number of Symbols %d\n", len(symbolBuffer))

	offset := 0

	var code = 0x00

	NewHuffmanTable := make(HuffmanTable)

	for i := 0; i < 16; i++ {
		freq := int(frquencyCountBuffer[i])

		fmt.Printf("%02d (%02d): ", i+1, freq)

		for j := 0; j < int(freq); j++ {

			sym := symbolBuffer[offset+j]
			NewHuffmanTable[i+1] = append(NewHuffmanTable[i+1], []byte{byte(code), sym})
			fmt.Printf("%02x ", sym)

			switch i {
			case 0:
				fmt.Printf("(%01b) ", code)
			case 1:
				fmt.Printf("(%02b) ", code)
			case 2:
				fmt.Printf("(%03b) ", code)
			case 3:
				fmt.Printf("(%04b) ", code)
			case 4:
				fmt.Printf("(%05b) ", code)
			case 5:
				fmt.Printf("(%06b) ", code)
			case 6:
				fmt.Printf("(%07b) ", code)
			case 7:
				fmt.Printf("(%08b) ", code)
			case 8:
				fmt.Printf("(%09b) ", code)
			case 9:
				fmt.Printf("(%010b) ", code)
			case 10:
				fmt.Printf("(%011b) ", code)
			case 11:
				fmt.Printf("(%012b) ", code)
			case 12:
				fmt.Printf("(%013b) ", code)
			case 13:
				fmt.Printf("(%014b) ", code)
			case 14:
				fmt.Printf("(%015b) ", code)
			case 15:
				fmt.Printf("(%016b) ", code)

			}
			code = code + 1
		}

		code = code << 1

		offset += int(freq)
		fmt.Println("")

	}

	huffmanTables[class] = append(huffmanTables[class], NewHuffmanTable)

	fmt.Println("")
}
