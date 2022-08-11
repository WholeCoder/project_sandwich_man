package main

import (
	"encoding/binary"
	"encoding/json"
    "io/ioutil"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("args ==", os.Args)
	if len(os.Args) < 3 {
		fmt.Println("Must specify file to be decompressed as first command line parameter.")
		fmt.Println("Must specify new file to be decompressed into as second command line parameter.")
		fmt.Println("************************************")
		fmt.Println("*            Usage                 *")
		fmt.Println("* decompress.exe infile.cmp outfile*")
		fmt.Println("************************************")
		return
	} else {
		fmt.Println("Compressing ->", os.Args[1], " ->", os.Args[2])
	}

    readInBytesForHashUnmarshalling, err := ioutil.ReadFile(os.Args[1])

	sizeOfHashReadFromDiskInBytes := uint64(binary.BigEndian.Uint64(readInBytesForHashUnmarshalling[:8]))

	var s2 string = string(readInBytesForHashUnmarshalling[8:sizeOfHashReadFromDiskInBytes])
	var encodingHash = map[string]Node{}
	err = json.Unmarshal([]byte(s2), encodingHash)
	if err != nil {
		panic(err)
	}

    // *** fix Parent pointers on node tree pionted to by encodingHash ***


	fmt.Println("read in bytes for hash unmarshalling: ", readInBytesForHashUnmarshalling)

	sizeOfCompressedTextReadFromDiskInBytes := uint64(binary.BigEndian.Uint64(readInBytesForHashUnmarshalling[sizeOfHashReadFromDiskInBytes+8:sizeOfHashReadFromDiskInBytes+8+8]))
	fmt.Println("sizeOfCompressedTextReadFromDiskInBytes2 =", sizeOfCompressedTextReadFromDiskInBytes)

	bitsetReadIn := InitNewByteset(readInBytesForHashUnmarshalling[8+sizeOfHashReadFromDiskInBytes:8+sizeOfHashReadFromDiskInBytes+int(sizeOfCompressedTextReadFromDiskInBytes)])

	decoding := ""
	var idx int = 0
	for idx < int(sizeReadFromDiskInBytes2) {
		br := root
		for len(br.Letter_s) > 1 {
			currentBit := bitsetReadIn.GetBit(idx)
			if currentBit {
				br = br.Right
			} else {
				br = br.Left
			}
			idx++
		}
		decoding = decoding + br.Letter_s
	}

	// Open a new file for writing only
	file, err := os.OpenFile(
		os.Args[2],
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	bytesWritten, err := file.Write([]byte(decoding))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)

}
