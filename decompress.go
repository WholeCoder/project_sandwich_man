package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
fmt.Println("args ==",os.Args)
    if len(os.Args) < 3 {
        fmt.Println("Must specify file to be decompressed as first command line parameter.")
        fmt.Println("Must specify new file to be decompressed into as second command line parameter.")
        fmt.Println("************************************")
        fmt.Println("*            Usage                 *")
        fmt.Println("* decompress.exe infile.cmp outfile*")
        fmt.Println("************************************")
        return
    } else {
        fmt.Println("Compressing ->",os.Args[1]," ->", os.Args[2])
    }

	hash := initFrequencyHash("words.txt")
	//hashForEncoding := map[string]*Node{}

	//printHash(hash)
	encodingHash := map[string]string{} //buildEncodingHash(&hashForEncoding)

	binaryTree := initBinaryTree(&hash, &encodingHash)

	//printHashNodePointer(hashForEncoding)

	//fmt.Println(binaryTree)

	// build hash used to encode letters to binary sequences
	//printEncodingHash(encodingHash)

    readInBytes, err := RetrieveROM(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("read in bytes: ", readInBytes)

    sizeReadFromDiskInBits := uint64(binary.BigEndian.Uint64(readInBytes[0:8]))
	fmt.Println("sizeReadFromDiskInBits = ", sizeReadFromDiskInBits)

    bitsetReadIn := InitNewByteset(readInBytes[8:])

    decoding := ""
    var idx int = 0
    for idx < int(sizeReadFromDiskInBits) {
        br := binaryTree.Root
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

func compressText(encodingHash *map[string]string, originalText string) string {
	compressed := ""
	for _, letter := range originalText {
		compressed = compressed + (*encodingHash)[strings.ToLower(string(letter))]
	}

	return compressed
}

func printOutPathOfNodeToRoot(node *Node) {

	fmt.Println("\n^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ START")
	count := 1
	for node != nil {
		fmt.Println(count, ": ", node.Letter_s)
		node = node.Parent
		count++
	}
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^END\n")
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to open file:  %s", err)
	}
}

func printEncodingHash(encodingHash map[string]string) {
	fmt.Println("------------------")

	for key, value := range encodingHash {
		fmt.Println("encodingHash[ ", key, " ] = ", value)
	}

	fmt.Println("------------------")
}

func buildEncodingHash(hashForEncoding *map[string]*Node) *map[string]string {
	encodingHash := map[string]string{}

	for key, value := range *hashForEncoding {
		if len(key) == 1 {
			encodingHash[key] = buildEncoding(value)
		}
	}

	return &encodingHash
}

func buildEncoding(node *Node) string {
	encoding := ""
	n := node

	count := 1
	for n != nil {
		encoding = n.ChildNodeRorL + encoding
		fmt.Println(n.Letter_s+"\tencoding:  ", encoding, "\tcount =", count)
		n = n.Parent
		count++
	}
	return encoding
}

func initBinaryTree(hash *map[string]Node, encodingHash *map[string]string) BinarySearchTree {

	for len(*hash) > 1 {
		// findFreeMinNode will remove the nodes from the hash
		nextNode := findFreeMinNode(hash)
		nextNode.ChildNodeRorL = "0"

		secondNode := findFreeMinNode(hash)
		secondNode.ChildNodeRorL = "1"

		//        if len((*nextNode).Letter_s) == 1 {
		//            (*hashForEncoding)[(*nextNode).Letter_s] = nextNode
		//        }

		//        if len((*secondNode).Letter_s) == 1 {
		//            (*hashForEncoding)[(*secondNode).Letter_s] = secondNode
		//        }

		newNode := createNewNodeFrom(nextNode, secondNode)
		fmt.Println("Letter_s:  ", newNode.Letter_s, "\t Child is:  ", nextNode.Left, "\t", secondNode.Left)
		(*hash)[newNode.Letter_s] = *newNode
	}

	var n Node
	for _, value := range *hash {
		fmt.Println("------------------------------Should only be printed once")
		n = value
	}

	fixBinaryTree(&n) // sorry folks!!
	fixEncodingHash(&n, encodingHash)

	bSearchTree := BinarySearchTree{Root: &n}
	return bSearchTree
}

func fixEncodingHash(node *Node, encodingHash *map[string]string) {
	if node == nil {
		return
	}

	if len(node.Letter_s) == 1 {
		(*encodingHash)[node.Letter_s] = buildEncoding(node)
	}
	fixEncodingHash(node.Left, encodingHash)
	fixEncodingHash(node.Right, encodingHash)
}

func fixBinaryTree(n *Node) {
	if n.Left != nil {
		n.Left.Parent = n
		fixBinaryTree(n.Left)
	}
	if n.Right != nil {
		n.Right.Parent = n
		fixBinaryTree(n.Right)
	}
}

func printNodeDetails(n *Node) {
	if n == nil {
		fmt.Println("-----CHAIN STOPPED----")
		return
	}
	fmt.Println("Letter_s: ", n.Letter_s, " ChildNodeRorL: ", n.ChildNodeRorL)
	printNodeDetails(n.Parent)
	fmt.Println("\n\n")
}

func findAndReturnANode(n *Node, nodeName string) *Node {
	if n == nil {
		return nil
	}
	if n.Letter_s == nodeName {
		return n
	}
	n1 := findAndReturnANode(n.Left, nodeName)
	if n1 != nil {
		return n1
	}
	n2 := findAndReturnANode(n.Right, nodeName)
	if n2 != nil {
		return n2
	}
	return nil
}

func printOutWholeTreeInOrder(n *Node, encodingHash *map[string]string) {
	if n == nil {
		return
	}
	printOutWholeTreeInOrder(n.Left, encodingHash)
	fmt.Println()
	//fmt.Println(n.Letter_s)

	if len(n.Letter_s) == 1 {
		(*encodingHash)[n.Letter_s] = buildEncoding(n)

		//printOutPathOfNodeToRoot(n)
	}

	if n.Left != nil {
		//fmt.Println("\t",n.Left.Letter_s)
	}

	if n.Right != nil {
		//fmt.Println("\t",n.Right.Letter_s)
	}
	if n.Parent != nil {
		//fmt.Println("Parent: ", n.Parent.Left, n.Parent.Data, n.Parent.Letter_s, n.Parent.Right, n.Parent, n.ChildNodeRorL, n.AlreadyUsedToBuildBinaryTree)
	}

	printOutWholeTreeInOrder(n.Right, encodingHash)
}

func debugCountHowManyLeftNodes(node *Node) {
	fmt.Println("----------")

	for node != nil {
		fmt.Println("\t", node.Letter_s)
		if node.Right != nil && !true {
			node = node.Right
		} else {
			node = node.Left
		}
	}

	fmt.Println("----------")
}

// Find and remove node from hash.  Return the node
func findFreeMinNode(hash *map[string]Node) *Node {
	var minKey string
	var minValue Node

	for key, value := range *hash {
		minKey = key
		minValue = value
		break
	}

	for key, value := range *hash {
		if value.Data < minValue.Data { // Data is the Probability
			minKey = key
			minValue = value
		}
	}

	hashMinValue := (*hash)[minKey]
	nodeMinValue := Node{Left: hashMinValue.Left,
		Data:     hashMinValue.Data, // Data is Probability
		Letter_s: minKey,
		Right:    hashMinValue.Right,

		Parent: nil,

		ChildNodeRorL: hashMinValue.ChildNodeRorL}
	delete(*hash, minKey)

	return &nodeMinValue
}

func createNewNodeFrom(node1, node2 *Node) *Node {
	newNode := Node{Left: node1, Data: node1.Data + node2.Data, Letter_s: node1.Letter_s + node2.Letter_s, Right: node2, Parent: nil}

	node1.Parent = &newNode
	node1.ChildNodeRorL = "0"

	node2.Parent = &newNode
	node2.ChildNodeRorL = "1"
	if len(node1.ChildNodeRorL) > 1 {
		fmt.Println(">>>>>>>>>>>>\t", node1.ChildNodeRorL)
	}
	if len(node2.ChildNodeRorL) > 1 {
		fmt.Println(">>>>>>>>>>>>\t", node2.ChildNodeRorL)
	}
	return &newNode
}

func initFrequencyHash(fileName string) map[string]Node {
	readFile, err := os.Open(fileName)
	check(err)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lines []string

	for fileScanner.Scan() {
		lines = append(lines, strings.ToLower(fileScanner.Text()))
	}

	hash := map[string]int{}
	for _, eachline := range lines {
		letters := strings.Split(eachline, "")
		for _, letter := range letters {
			if strings.ContainsAny(letter, "abcdefghijklmnopqrstuvwxyz") {
				hash[letter]++
			}
		}
	}

	totalLetters := 0
	for _, value := range hash {
		totalLetters += value
	}

	freqNodemap := map[string]Node{}

	for key, value := range hash {
		freqNodemap[key] = Node{Data: float64(value) / float64(totalLetters), AlreadyUsedToBuildBinaryTree: false}
	}

	return freqNodemap
}

func printHash(hash map[string]Node) {
	keys := "abcdefghijklmnopqrstuvwxyz"

	for _, akey := range []rune(keys) {
		fmt.Println("hash[ ", string(akey), " ] = ", hash[string(akey)])
	}
}
func printHashNodePointer(hash map[string]*Node) {
	keys := "abcdefghijklmnopqrstuvwxyz"

	for _, akey := range []rune(keys) {
		fmt.Println("hash[ ", string(akey), " ] = ", *hash[string(akey)])
	}
}
