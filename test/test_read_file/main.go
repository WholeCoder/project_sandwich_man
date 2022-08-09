package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dat, _ := ioutil.ReadFile(fileName)
	asString := string(dat)

	hash := map[string]int{}

	for _, ch := range asString {
		fmt.Print(string(ch), " ")
		hash[string(ch)] += 1
	}
	fmt.Println()
	for v, k := range hash {

	}
}
