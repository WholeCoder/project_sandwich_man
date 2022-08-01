package main

type Node struct {
    Left  *Node
    Data  float64 // This is actually the probability of finding this letter.
    Letter_s string // This will contain 1 letter if it is a leaf or more if internal.
    Right *Node

    Parent *Node

    ChildNodeRorL string // Contains 1 or 0 depending if it is a right or left child of parent node.  0 mean 0 and 1 means 1
    AlreadyUsedToBuildBinaryTree bool
}

type BinarySearchTree struct {
    Root *Node
}
