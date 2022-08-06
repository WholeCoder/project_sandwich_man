package main

import "fmt"

func main() {
    n := &Node{test:&Node{test:&Node{test: nil, value: "c"}, value: "b"}, value: "a"}

    for n != nil {
        fmt.Println(n.value)
        n = n.test
    }
}

type Node struct {
    test *Node
    value string
}
