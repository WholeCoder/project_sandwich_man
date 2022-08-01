package main

import (
    "strings"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	hash := initFrequencyHash("words.txt")
    hashForEncoding := map[string]*Node{}

	//printHash(hash)
    binaryTree := initBinaryTree(hash, hashForEncoding)

    //printHashNodePointer(hashForEncoding)

    fmt.Println(binaryTree)

    // build hash used to encode letters to binary sequences
    encodingHash := buildEncodingHash(hashForEncoding)

}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to open file:  %s", err)
	}
}

func buildEncodingHash(hashForEncoding map[string]*Node) map[string]string {
    encodingHash := map[string]string{}

    for key, value := range hashForEncoding {
        encodingHash[key] = buildEncoding(value)
    }

    return encodingHash
}

func buildEncoding(node *Node) string {
    encoding := ""

    for node != nil {
        encoding = node.ChildNodeRorL + encoding
    }

    return encoding
}

func initBinaryTree(hash map[string]Node, hashForEncoding map[string]*Node) BinarySearchTree {

    for len(hash) > 1 {
        // findFreeMinNode will remove the nodes from the hash
        nextNode := findFreeMinNode(&hash)
        secondNode := findFreeMinNode(&hash)

        if len((*nextNode).Letter_s) == 1 {
            hashForEncoding[(*nextNode).Letter_s] = nextNode
        }

        if len((*secondNode).Letter_s) == 1 {
            hashForEncoding[(*secondNode).Letter_s] = secondNode
        }

        newNode := createNewNodeFrom(nextNode, secondNode)
        hash[newNode.Letter_s] = *newNode
    }

    var n Node
    for _, value := range hash {
        n = value
    }

    bSearchTree := BinarySearchTree{Root: &n}
    return bSearchTree
}

func debugCountHowManyLeftNodes(node *Node) {
    fmt.Println("----------")

    for node != nil {
        fmt.Println("\t", node.Letter_s)
        if node.Right != nil && !true{
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
        if value.Data< minValue.Data { // Data is the Probability
            minKey = key
            minValue = value
        }
    }


    hashMinValue := (*hash)[minKey]
    delete(*hash, minKey)

    nodeMinValue := Node{Left: hashMinValue.Left,
                         Data: hashMinValue.Data, // Data is Probability
                         Letter_s: minKey,
                         Right: hashMinValue.Right,

                         Parent: nil,

                         ChildNodeRorL: "" }
    return &nodeMinValue
}

func createNewNodeFrom(node1, node2 *Node) *Node {
    newNode := Node{Left:node1, Data: node1.Data+node2.Data, Letter_s: node1.Letter_s + node2.Letter_s, Right:node2, Parent:nil}

    node1.Parent = &newNode
    node1.ChildNodeRorL = "0"

    node2.Parent = &newNode
    node2.ChildNodeRorL = "1"

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
