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
	name_of_new_archive_file := os.Args[1]
	files_to_compress_and_archive := os.Args[2:]

	for _, filename := range files_to_compress_and_archive {
		compress_main(filename, filename+".comp")
	}

	files_to_be_compressed_marshalled, err := json.Marshal(files_to_compress_and_archive)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(
		name_of_new_archive_file,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()

	// Write length of slice of filename strings to the file.
	bytesWritten, err := file.Write(getBytesForInt(len(files_to_be_compressed_marshalled)))
	if err != nil {
		panic(err)
	}

	_, err = file.Write(files_to_be_compressed_marshalled)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)

	files_in_bytes := []byte{}

	for _, filename := range files_to_compress_and_archive {
		readInBytes, err := ioutil.ReadFile(filename + ".comp")
		if err != nil {
			panic(err)
		}
		length_of_compressed_file_with_comrpession_hash := make([]byte, 8)

		binary.BigEndian.PutUint64(length_of_compressed_file_with_comrpession_hash, uint64(len(readInBytes)))

		files_in_bytes = append(files_in_bytes, length_of_compressed_file_with_comrpession_hash...)
		files_in_bytes = append(files_in_bytes, readInBytes...)
	}

	file.Write(files_in_bytes)
	fmt.Println("Compressed Following Files:  ")
	for _, filename := range files_to_compress_and_archive {
		fmt.Println("\t", filename)
	}
	fmt.Println("\nInto Following Archive:  ")
	fmt.Println("\t", name_of_new_archive_file)
}
