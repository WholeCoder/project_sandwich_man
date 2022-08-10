package main

import (
	"encoding/binary"
	"encoding/json"
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

	hash := initFrequencyHash("words.txt")
	//hashForEncoding := map[string]*Node{}

	//printHash(hash)
	encodingHash := map[string]string{} //buildEncodingHash(&hashForEncoding)

	root := initBinaryTree(&hash, &encodingHash)
	if err != nil {
		log.Fatal(err)
	}

	sizeReadFromDiskInBytes := uint64(binary.BigEndian.Uint64(readInBytesForHashUnmarshalling[0:8]))

	var s2 string = string(readInBytesForHashUnmarshalling[8:sizeReadFromDiskInBytes])
	var encodingHash = map[string]Node{}
	err = json.Unmarshal([]byte(s2), encodingHash)
	if err != nil {
		panic(err)
	}

	fmt.Println("read in bytes: ", readInBytes)

	sizeReadFromDiskInBits := uint64(binary.BigEndian.Uint64(readInBytes[0:8]))
	fmt.Println("sizeReadFromDiskInBits = ", sizeReadFromDiskInBits)

	bitsetReadIn := InitNewByteset(readInBytes[8:])

	decoding := ""
	var idx int = 0
	for idx < int(sizeReadFromDiskInBits) {
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
