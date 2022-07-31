package main

import (
    "encoding/json"
    "strings"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	hash := initFrequencyHash("words.txt")

	printHash(hash)
    binaryTree := initBinaryTree(hash)

    fmt.Println("tree data =", binaryTree.Root.Data)

    bytes, err := json.Marshal(&binaryTree)
    if err != nil {
        fmt.Println("Can't serislize", binaryTree)
        return
    }

    var tr3 BinarySearchTree
    err = json.Unmarshal(bytes, &tr3)
    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Printf("\nbytes = %#v\n", string(bytes))
    fmt.Printf("\nroot = %#v\n", tr3.Root)
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to open file:  %s", err)
	}
}

func initBinaryTree(hash map[string]Node) BinarySearchTree {

    for len(hash) > 1 {
        // findFreeMinNode will remove the nodes from the hash
        nextNode := findFreeMinNode(&hash)
        secondNode := findFreeMinNode(&hash)

        newNode := createNewNodeFrom(nextNode, secondNode)
        hash[newNode.Letter_s] = *newNode
    }

    n := hash["abcdefghijklmnopqrstuvwxyz"]
    bSearchTree := BinarySearchTree{Root: &n}
    return bSearchTree
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

    nodeMinValue := Node{Left: nil,
                         Data: hashMinValue.Data, // Data is Probability
                         Letter_s: minKey,
                         Right: nil,

                         Parent: nil,

                         ChildNodeRorL: "" }
    return &nodeMinValue
}

func createNewNodeFrom(node1, node2 *Node) *Node {
    newNode := Node{Left:node1, Data: node1.Data+node2.Data, Letter_s: node1.Letter_s + node2.Letter_s, Right:node2, Parent:nil}

    node1.Parent = &newNode
    node1.ChildNodeRorL = "L0"

    node2.Parent = &newNode
    node2.ChildNodeRorL = "R1"

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
