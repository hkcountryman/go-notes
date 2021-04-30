package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// Modify the echo program to also print os.Args[0], the name of the command that invoked it
func ex1() {
	fmt.Println(strings.Join(os.Args[:], " "))
}

// Modify the echo program to print the index and value of each of its arguments, one per line
func ex2() {
	fmt.Println("Index\tArgument")
	for idx, arg := range os.Args[1:] {
		fmt.Printf("%d\t%s\n", idx, arg)
	}
}

// Modify dup2 to print the names of all files in which each duplicated line occurs
func ex4() {
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, "")
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, arg)
			f.Close()
		}
	}
}

func countLines(f *os.File, name string) {
	anyRepeats := false
	counts := make(map[string]int)
	input := bufio.NewScanner(f)
	for input.Scan() { // (ignoring potential errors from input.Err())
		counts[input.Text()]++
		if counts[input.Text()] > 1 {
			anyRepeats = true
		}
	}
	if name != "" && anyRepeats {
		fmt.Println("Repeated lines in file \"" + name + "\":")
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

//
func main() {
	fmt.Println()
	ex4()
	fmt.Println()
}
