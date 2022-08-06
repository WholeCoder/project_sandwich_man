package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
fmt.Println("args ==",os.Args)
    if len(os.Args) < 3 {
        fmt.Println("Must specify file to be compressed as first command line parameter.")
        fmt.Println("Must specify new file to be compressed into as second command line parameter.")
        fmt.Println("***********************************")
        fmt.Println("*            Usage                *")
        fmt.Println("* compress.exe infile outfile.cmp *")
        fmt.Println("***********************************")
        return
    } else {
        fmt.Println("Deompressing ->",os.Args[1]," ->", os.Args[2])
    }

	hash := initFrequencyHash("words.txt")
	//hashForEncoding := map[string]*Node{}

	//printHash(hash)
	encodingHash := map[string]string{} //buildEncodingHash(&hashForEncoding)

	initBinaryTree(&hash, &encodingHash)

	//printHashNodePointer(hashForEncoding)

	//fmt.Println(binaryTree)

	// build hash used to encode letters to binary sequences
	//printEncodingHash(encodingHash)

    originalTextBytes, err := ReadInBytesFromFile(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    originalText := string(originalTextBytes)

	compressedText := compressText(&encodingHash, originalText)

	fmt.Println("originalText:  ", originalText)
	fmt.Println("compressed  :  ", compressedText)

	lengthOfCompressedText := len(compressedText)
	byteLengthOfCompressedText := uint64(math.Ceil(float64(lengthOfCompressedText)/8.0) + 8.0) // add 8.0 bytes for this size (byteLengthOfCompressedText

	fmt.Println("lengthOfCompressedText: ", lengthOfCompressedText)
	fmt.Println("byteLengthOfCompressedText: ", byteLengthOfCompressedText)

	bray := make([]byte, byteLengthOfCompressedText)

	fmt.Println(byteLengthOfCompressedText)

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(lengthOfCompressedText))

	count := 0
	for count < 8 {
		bray[count] = b[count]
		count++
	}

	compressedTextAsByteRay := InitNewByteset(bray)
	fmt.Println("bray = ", bray)

	for index, number := range compressedText {
		compressedTextAsByteRay.SetBit(index+64, string(number) == "1")
	}

	fmt.Println("compressedTextAsByteRay = ", compressedTextAsByteRay)
	//    node := hashForEncoding["z"]
	i := uint64(binary.BigEndian.Uint64(compressedTextAsByteRay[0:8]))
	fmt.Println("length of data in bray = ", i)

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
	bytesWritten, err := file.Write(compressedTextAsByteRay)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)


}

