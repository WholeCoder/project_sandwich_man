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
	if len(os.Args) < 3 {
		fmt.Println("Must specify file to be compressed as first command line parameter.")
		fmt.Println("Must specify new file to be compressed into as second command line parameter.")
		fmt.Println("************************************")
		fmt.Println("*            Usage                 *")
		fmt.Println("* ",os.Args[0],"infile outfile.cmp *")
		fmt.Println("************************************")
		return
	} else {
		fmt.Println("Compressing ->", os.Args[1], " ->", os.Args[2])
	}

	hash := initFrequencyHash(os.Args[1])

	encodingHash := map[string]string{}

	initBinaryTree(&hash, &encodingHash)

	originalTextBytes, err := ReadInBytesFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	originalText := string(originalTextBytes)

    fmt.Println("**************************************************** original Text ***********************************************")
    fmt.Println("originalText: ", originalText)
    fmt.Println("len( originalText ) =", len(originalText))
    fmt.Println("******************************************************************************************************************")

	compressedText := compressText(&encodingHash, originalText)

	lengthOfCompressedTextInBytes := int(math.Ceil(float64(len(compressedText))/8.0))
    lengthOfCompressedTextInBits := len(compressedText)
    fmt.Println("3.  length of compressed text in bytes =", lengthOfCompressedTextInBytes)
    fmt.Println("4.  length of compressed text in bits =", lengthOfCompressedTextInBits)
    fmt.Println("**************************************************** compressed text as string of 0/1's *************************")
    fmt.Println("5.  compressed text =", compressedText)
    fmt.Println("*****************************************************************************************************************")
	// Marshall - initFrequencyHash returns the hash with Node as value that are nil
	hashForDecompression := initFrequencyHashWithFloat64ForValues(os.Args[1])

	// write this to fileInBytesInMemory
	hashMarshalled, err := json.Marshal(hashForDecompression)
	if err != nil {
		panic(err)
	}

	// write this to fileInBytesInMemory
	marshalledHashDecompressionLength := len(hashMarshalled)

    fmt.Println("1.  length of marshalled float64 hash =", marshalledHashDecompressionLength)
    fmt.Println("2.  hashForDecompressioned marshalled =", string(hashMarshalled))

	// use this for lenght of file in bytes
	byteLengthOfCompressedTextWithAdditional := uint64(uint64(math.Ceil(float64(lengthOfCompressedTextInBits)/8.0)) + 8.0 + 8.0 + uint64(marshalledHashDecompressionLength) + 8.0) // add 8.0 bytes for this size byteLengthOfCompressedText and add 8.0 for length of marshalledHashDecompressionLength (8) plus lenght of hashMarshalled and plus length of compressed text in bits (8)

	// This is actually the contents of the file (write it to the file).
	fileInBytesInMemory := make([]byte, byteLengthOfCompressedTextWithAdditional)

	marshalledHashDecompressionLengthMarshalled := getBytesForInt(marshalledHashDecompressionLength)
	// hashMarshalled
	lengthOfCompressedTextInBytesMarshalled := getBytesForInt(lengthOfCompressedTextInBytes)
	lengthOfCompressedTextInBitsMarshalled := getBytesForInt(lengthOfCompressedTextInBits)
	// compressedTextAsByteRay

	count := 0
	for count < 8 {
		fileInBytesInMemory[count] = marshalledHashDecompressionLengthMarshalled[count]
		count++
	}

	for count < 8+len(hashMarshalled) {
		fileInBytesInMemory[count] = hashMarshalled[count-8]
        count++
	}

	for count < 8+len(hashMarshalled)+8 {
		fileInBytesInMemory[count] = lengthOfCompressedTextInBytesMarshalled[count-8-len(hashMarshalled)]
		count++
	}

	for count < 8+len(hashMarshalled)+8+8 {
		fileInBytesInMemory[count] = lengthOfCompressedTextInBitsMarshalled[count-8-8-len(hashMarshalled)]
		count++
	}

	compressedTextAsByteRay := InitNewByteset(fileInBytesInMemory)

	for index, number := range compressedText {
        if index > lengthOfCompressedTextInBits {
            break
        }
		compressedTextAsByteRay.SetBit(count*8+index, string(number) == "1")
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
	bytesWritten, err := file.Write(compressedTextAsByteRay)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)

}

func getBytesForInt(length int) []byte {

	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, uint64(length))

	return b
}
