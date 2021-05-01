package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Prints command line arguments.
func echo1() {
	var s, sep string
	// for initialization; condition; post { ... }
	for i := 1; i < len(os.Args); i++ { // os.Args[0] = chap1.go; rest are command line arguments
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Prints command line arguments.
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

// Prints command line arguments; less garbage collecting than echo1() or echo2().
func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

// Prints the text of each line that appears more than once in the standard input, preceded by count.
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
// from a list of named files.
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

// Main functionality of dup2().
func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() { // (ignoring potential errors from input.Err())
		counts[input.Text()]++
	}
}

// Prints the count and text of lines that appear more than once in the named input files, which are
// processed all at once.
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

// Generates GIF animations of random Lissajous figures:
var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

// This contains errors, don't know why!
func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// Fetches the content of a specified URL and prints it as uninterpreted text.
func fetch1(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	fmt.Printf("%s", b)
}

// Fetches URLs in parallel and reports their times and sizes (given os.Args[1:]).
func fetchall(args []string) {
	start := time.Now()
	ch := make(chan string)
	for _, url := range args {
		go fetch2(url, ch) // start a goroutine
	}
	for range args {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// Helper function for fetchall.
func fetch2(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs	%7d		%s", secs, nbytes, url)
}

// Goroutine: concurrent function execution
// Channel: communication mechanism that allows one goroutine to pass values of a specified type to
// another goroutine
// main() runs in a goroutine and the "go" statement creates additional goroutines

// Minimal "echo" server.
func server1() {
	// connects a handler function to incoming URLs whose path begins with / (all URLs)
	http.HandleFunc("/", handler1) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Handler echoes the Path component of the requested URL.
func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// Add some features to the server...
var mu sync.Mutex
var count int

// Minimal "echo" and counter server.
func server2() {
	http.HandleFunc("/", handler2)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}

// Echoes the Path component of the requested URL.
func handler2(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// Echoes the number of calls so far
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

// The above server is running the handler for each incoming request in separate goroutines, but it
// doesn't do this for requests to update count, which may lead to a race condition. mu.Lock() and
// mu.Unlock() ensure that at most one goroutine accesses a given variable at a time.

// Handler could also echo the HTTP request for debugging:
func handler3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil { // shorter and reduces scope of err
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// Puts our Lissajous GIF onto a server.
func lissajousServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	})
	log.Fatal(http.ListenAndServe("localhost:8100", nil))
}

// Switch statements compare result of function called (e.g., coinflip()) to each case:
//	switch coinflip() {
//		case "heads":
//			heads++
//		case "tails":
//			tails++
//		default:
//			fmt.Println("Landed on edge!")
//	}
// Alternatively, "tagless" switch statements don't need an operand; can just list boolean cases:
//	func Signum(x int) int {
//		switch {
//		case x > 0:
//			return +1
//		default:
//			return 0
//		case x < 0:
//			return -1
//		}
//	}
// Essentially, equivalent to "switch true {...}".
// Like for and if statements, switch may include an optional simple statement: short variable
// declaration, increment or assignment statement, or a function call.
// Can use break and continue statements to modify flow control: break causes control to resume at
// the next statement following the innermost for, switch or select statement; continue causes the
// innermost for loop to skip current iteration and go on to next. We can label statements and refer
// to them by name with break or continue, or even use goto statements, though that's not really
// intended for use by humans.

// Name an existing type with a type declaration:
//	type Point struct {
//		X, Y int
//	}

// Pointers are explicitly visible: the & operator yields the address of a variable and the *
// operator retrieves the variable that the pointer refers to, but there is no pointer arithmetic.

// A method is a function associated with a named type; may be attached to almost any named type.
// Interfaces are abstract types that define methods for concrete types.

// Index of standard library packages available at https://golang.org/pkg and community packages at
// https://godoc.org. Use go doc tool to see documentation, e.g.,
//	$ go doc http.ListenAndServe
//	package http // import "net/http"
//
//	func ListenAndServe(addr string, handler Handler) error
//		ListenAndServe listens on the TCP network address addr and then
//		calls Serve with handler to handle requests on incoming connections.
//	...

// A partial driver to demonstrate the above examples
func main() {
	/*fmt.Println("\nFirst implementation of echo:")
	echo1()
	fmt.Println("\nSecond implementation of echo:")
	echo2()
	fmt.Println("\nThird implementation of echo:")
	echo3()
	//lissajous(os.Stdout) // go build ch1/notes.go > ch1/out.gif
	fetch1("http://gopl.io")
	fmt.Println()
	fetch1("http://bad.gopl.io")*/
	//fetchall(os.Args[1:]) // go run ch1/notes.go https://golang.org http://gopl.io https://godoc.org
	lissajousServer() // go run notes.go &
}
