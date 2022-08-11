package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Used
func compressText(encodingHash *map[string]string, originalText string) string {
	compressed := ""
	for _, letter := range originalText {
		compressed = compressed + (*encodingHash)[string(letter)]
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

// Used
func initBinaryTree(hash *map[string]Node, encodingHash *map[string]string) *Node {

	for len(*hash) > 1 {
		// findFreeMinNode will remove the nodes from the hash
		nextNode := findFreeMinNode(hash)
		nextNode.ChildNodeRorL = "0"

		secondNode := findFreeMinNode(hash)
		secondNode.ChildNodeRorL = "1"

		newNode := createNewNodeFrom(nextNode, secondNode)
		(*hash)[newNode.Letter_s] = *newNode
	}

	var n Node
	for _, value := range *hash { // Runs Once.
		n = value
	}

	fixBinaryTree(&n) // sorry folks!!
	fixEncodingHash(&n, encodingHash)
	return &n
}

// Used
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

// Used
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

// Used
func createNewNodeFrom(node1, node2 *Node) *Node {
	newNode := Node{Left: node1, Data: node1.Data + node2.Data, Letter_s: node1.Letter_s + node2.Letter_s, Right: node2, Parent: nil}
	return &newNode
}

// Used
func initFrequencyHash(fileName string) map[string]Node {

	dat, err := ioutil.ReadFile(fileName)
	check(err)
	asString := string(dat)

	hash := map[string]int{}

	for _, ch := range asString {
		hash[string(ch)] += 1
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

// Used
func initFrequencyHashWithFloat64ForValues(fileName string) map[string]float64 {

	dat, err := ioutil.ReadFile(fileName)
	check(err)
	asString := string(dat)

	hash := map[string]int{}

	for _, ch := range asString {
		hash[string(ch)] += 1
	}

	totalLetters := 0
	for _, value := range hash {
		totalLetters += value
	}

	freqNodemap := map[string]float64{}

	for key, value := range hash {
		freqNodemap[key] = float64(value) / float64(totalLetters)
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
