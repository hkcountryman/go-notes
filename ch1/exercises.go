package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

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

// Change the Lissajous program's color palette to green on black, for added authenticity. To create
// the web color #RRGGBB, use color.RGBA{0xRR, 0xGG, 0xBB, 0xff}.
func ex5(out io.Writer) {
	palette := []color.Color{color.RGBA{0x00, 0x80, 0x00, 0xff}, color.Black}
	const (
		whiteIndex = 0 // first color in palette
		blackIndex = 1 // next color in palette
	)
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

// Modify the Lissajous program to produce images in multiple colors by adding more values to
// palette and then displaying them by changing the third argument of SetColorIndex.
// Doesn't really work as envisioned but this is a time sink I don't need to know how to do.
func ex6(out io.Writer) {
	red := color.RGBA{0x80, 0x00, 0x00, 0xff}
	orange := color.RGBA{0xff, 0x8c, 0x00, 0xff}
	yellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
	green := color.RGBA{0x00, 0x80, 0x00, 0xff}
	blue := color.RGBA{0x00, 0x00, 0x80, 0xff}
	purple := color.RGBA{0x80, 0x00, 0x80, 0xff}
	palette := []color.Color{color.White, red, orange, yellow, green, blue, purple, color.Black}
	const (
		whiteIndex  = 0
		redIndex    = 1
		orangeIndex = 2
		yellowIndex = 3
		greenIndex  = 4
		blueIndex   = 5
		purpleIndex = 6
		blackIndex  = 7
	)
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
	rand.Seed(time.Now().UnixNano())
	max := 7
	min := 1
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(rand.Intn(max-min+1)+min))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// The function call io.Copy(dst, src) reads from the src and writes to dst. Use it instead of
// ioutil.ReadAll to copy the response body to os.Stdout without requiring a buffer large enough to
// hold the entire stream. Be sure to check the error result of io.Copy.
func ex7(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	return resp
}

// Modify fetch to add the prefix http:// to each argument URL if it is missing. You might want to
// use strings.HasPrefix.
func ex8(url string) *http.Response {
	var resp *http.Response
	pre := "http://"
	if strings.HasPrefix(url, pre) {
		resp = ex7(url)
	} else {
		resp = ex7(pre + url)
	}
	return resp
}

// Modify fetch to also print the HTTP status code, found in resp.Status
func ex9(url string) {
	resp := ex8(url)
	fmt.Println("\nStatus code: " + resp.Status)
}

// Puts our Lissajous GIF onto a server.
func lissajousServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ex6(w) //ex5(w)
	})
	log.Fatal(http.ListenAndServe("localhost:8010", nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	/*fmt.Println()
	ex9("http://gopl.io")
	fmt.Println()*/
	lissajousServer()
}
