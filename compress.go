package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	fmt.Println("args ==", os.Args)
	if len(os.Args) < 3 {
		fmt.Println("Must specify file to be compressed as first command line parameter.")
		fmt.Println("Must specify new file to be compressed into as second command line parameter.")
		fmt.Println("***********************************")
		fmt.Println("*            Usage                *")
		fmt.Println("* compress.exe infile outfile.cmp *")
		fmt.Println("***********************************")
		return
	} else {
		fmt.Println("Deompressing ->", os.Args[1], " ->", os.Args[2])
	}

	hash := initFrequencyHash(os.Args[1])

	encodingHash := map[string]string{}

	initBinaryTree(&hash, &encodingHash)

	originalTextBytes, err := ReadInBytesFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	originalText := string(originalTextBytes)

	compressedText := compressText(&encodingHash, originalText)

	fmt.Println("originalText:  ", originalText)
	fmt.Println("compressed  :  ", compressedText)

	lengthOfCompressedText := len(compressedText)

	// Marshall
	hashForDecompression := initFrequencyHash(os.Args[1])

	// *** make all parents in tree nil ***
	hashForDecompression = nilOutParentNodes(hashForDecompression)

	// write this to fileInBytesInMemory
	hashMarshalled, err := json.Marshal(hashForDecompression)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashMarshalled)) //{"Name":"Amanda","Age":12}

	// write this to fileInBytesInMemory
	marshalledHashDecompressionLength := len(hashMarshalled)

	// use this for lenght of file in bytes
	byteLengthOfCompressedTextWithAdditional := uint64(math.Ceil(float64(lengthOfCompressedText)/8.0) + 8.0 + marshalledHashDecompressionLength + 8.0) // add 8.0 bytes for this size byteLengthOfCompressedText and add 8.0 for length of marshalledHashDecompressionLength (8) plus lenght of hashMarshalled

	fmt.Println("lengthOfCompressedText: ", lengthOfCompressedText)
	fmt.Println("byteLengthOfCompressedTextWithAdditional: ", byteLengthOfCompressedTextWithAdditional)

	// This is actually the contents of the file (write it to the file).
	fileInBytesInMemory := make([]byte, byteLengthOfCompressedTextWithAdditional)

	fmt.Println(byteLengthOfCompressedTextWithAdditional)

	marshalledHashDecompressionLengthMarshalled := getBytesForInt(mashallHashDecompressionLength)
	// hashMarshalled
	lengthOfCompressedTextMarshalled := getBytesForInt(lengthOfCompressedText)
	// compressedTextAsByteRay

	count := 0
	for count < 8 {
		fileInBytesInMemory[count] = marshalledHashDecompressionLengthMarshalled[count]
		count++
	}

	for count < 8+len(hashMarshalled) {
		fileInBytesInMemory[count] = hashMarshelled[count-8]
	}

	for count < 8+len(hashMarshalled)+8 {
		fileInBytesInMemory[count] = lengthOfCompressedText[count-8-len(hashMarshalled)]
		count++
	}

	compressedTextAsByteRay := InitNewByteset(fileInBytesInMemory)
	fmt.Println("fileInBytesInMemory = ", fileInBytesInMemory)

	for index, number := range compressedText {
		compressedTextAsByteRay.SetBit(index+64+64+len(hashMarshalled)*8, string(number) == "1")
	}

	fmt.Println("compressedTextAsByteRay = ", compressedTextAsByteRay)

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

func getBytesForInt(length int) []byte {

	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, uint64(lengthOfCompressedText))

	return b
}
