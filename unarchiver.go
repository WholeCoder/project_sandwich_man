package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	archive_name := os.Args[1]
	readInBytesForArchive, err := ioutil.ReadFile(archive_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("Length of bytes from archive: ", readInBytesForArchive)

	sizeOfSliceOfStringFilenames := uint64(binary.BigEndian.Uint64(readInBytesForArchive[:8]))
	fmt.Println("Length of files in archive: ", sizeOfSliceOfStringFilenames)

	var tempStringSliceForFilenames string = string(readInBytesForArchive[8 : sizeOfSliceOfStringFilenames+8])
	var sliceOfFilenames = []string{}
	err = json.Unmarshal([]byte(tempStringSliceForFilenames), &sliceOfFilenames)
	if err != nil {
		panic(err)
	}

	var sliceOfFileLengths = []uint64{}

	count := 8 + int(sizeOfSliceOfStringFilenames)
	fmt.Println("Files In Archive (", len(sliceOfFilenames), " ): ")
	for ; count < 8+int(sizeOfSliceOfStringFilenames); count += 8 {
		sliceOfFileLengths = append(sliceOfFileLengths, uint64(binary.BigEndian.Uint64(readInBytesForArchive[count:count+8])))
		fmt.Println("\t", sliceOfFilenames[len(sliceOfFileLengths)-1], "   length =", sliceOfFileLengths[len(sliceOfFileLengths)-1])
	}

	sliceOfFileNamesForDecompression := []string{}

	for i := 0; i < len(sliceOfFilenames); i++ {
		// Get the next file slice.
		size_of_compressed_file := uint64(binary.BigEndian.Uint64(readInBytesForArchive[count : count+8]))
		count += 8 // Move past the size of the compressed file

		// Open and write the next file slice.
		file, err := os.OpenFile(
			sliceOfFilenames[i]+".comp",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			log.Fatal(err)
		}

		sliceOfFileNamesForDecompression = append(sliceOfFileNamesForDecompression, sliceOfFilenames[i]+".comp")
		defer file.Close()

		bytesWritten, err := file.Write([]byte(readInBytesForArchive[count : count+int(size_of_compressed_file)]))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d bytes.\n", bytesWritten)

		// increment the count to position it on the next slice to be written
		count += int(size_of_compressed_file)
	}

	fmt.Println("\n*** Decompressing ***")
	for i := 0; i < len(sliceOfFileNamesForDecompression); i++ {
		decompress_main(sliceOfFileNamesForDecompression[i], sliceOfFilenames[i])
		os.Remove(sliceOfFileNamesForDecompression[i])
	}
	fmt.Println("\n   Finished Decompression ***")
	fmt.Println("***********************")
}
