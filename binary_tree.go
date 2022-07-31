package main

type Node struct {
    Left  *Node
    Data  float64 // This is actually the probability of finding this letter.
    Letter_s string // This will contain 1 letter if it is a leaf or more if internal.
    Right *Node

    Parent *Node
}

type BinarySearchTree struct {
    Root *Node
}
