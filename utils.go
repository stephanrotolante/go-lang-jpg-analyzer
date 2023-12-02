package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func ReadFunc(file *os.File, buffer []byte) (int, error) {
	n, err := file.Read(buffer)

	if n == 0 {
		fmt.Println("No Bytes read")
		return 0, errors.New("No Bytes Read")
	}

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return n, nil
}

func BigEUint16(arg1, arg2 byte) int {

	return int(binary.BigEndian.Uint16([]byte{
		arg1,
		arg2,
	}))
}
