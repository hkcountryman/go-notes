package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Prints command line arguments
func echo1() {
	var s, sep string
	// for initialization; condition; post { ... }
	for i := 1; i < len(os.Args); i++ { // os.Args[0] = chap1.go; rest are command line arguments
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Prints command line arguments
func echo2() {
	s, sep := "", ""
	// range returns index, value (index is anonymous)
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// Equivalent ways to declare a variable:
//		s := ""				(only for use within a function)
//		var s string		(only for initializing to default zero value, i.e., "")
//		var s = ""			(rarely used except to declare multiple variables)
//		var s string = ""	(explicit about type)

// Prints command line arguments; less garbage collecting than echo1() or echo2()
func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

// Prints the text of each line that appears more than once in the standard input, preceded by count
func dup1() {
	counts := make(map[string]int) // key is any type that can be compared with ==; value any type
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++ // the ++ and -- postfix notation is a statement, not an expression
	}
	for line, n := range counts { // (ignoring potential errors from input.Err())
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// "Verbs" for use in Printf():
//		%d				decimal integer
//		%x, %o, %b		hexadecimal, octal, binary integer
//		%f, %g, %e		float: 3.141593, 3.141592653589793, 3.141593.+00
//		%t				boolean: true or false
//		%c				rune (Unicode code point)
//		%s				string
//		%q				quoted string "abc" or rune 'c'
//		%v				any value in a natural format
//		%T				type of any value
//		%%				literal percent sign (no operand)

// Prints the count and text of lines that appear more than once in the input, either from stdin or
// from a list of named files
func dup2() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts) // pass by reference to counts: countLines CAN mutate the map
	} else {
		for _, arg := range files {
			// os.Open() returns open file (*os.File) and value of built-in error type, ideally nil
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// Main functionality of dup2()
func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() { // (ignoring potential errors from input.Err())
		counts[input.Text()]++
	}
}

// Prints the count and text of lines that appear more than once in the named input files, which are
// processed all at once
func dup3() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		// ReadFile() returns a byte slice that must be cast to a string to use strings.Split()
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// bufio.Scanner, ioutil.ReadFile(), and ioutil.WriteFile() use the Read() and Write() methods of
// *os.File, but are easier to use than those lower-level routines.

////////////////////////////////////////////////////////////////////////////////////////////////////

// A partial driver to demonstrate the above examples
func main() {
	fmt.Println("\nFirst implementation of echo:")
	echo1()
	fmt.Println("\nSecond implementation of echo:")
	echo2()
	fmt.Println("\nThird implementation of echo:")
	echo3()
	fmt.Println()
}
