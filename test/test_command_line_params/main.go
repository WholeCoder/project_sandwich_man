package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "io"
)

func main() {
    // Defining a boolean flag -l to count iines instead of words​
    lines := flag.Bool("l", false, "Count lines")

    // Parsing the flags provided by the user​
    flag.Parse()

    // Calling the count function to count the number of words (or lines)​
    // received from the Standard Input and printing it out​
    fmt.Println(count(os.Stdin, *lines))

}

func count(r io.Reader, countLines bool) int {
    scanner := bufio.NewScanner(r)

    if !countLines {
        scanner.Split(bufio.ScanWords)
    }

    wc := 0

    for scanner.Scan() {
        wc++
    }

    return wc
}
