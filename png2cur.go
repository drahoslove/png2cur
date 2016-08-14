/**
 * Convert PNG image to microsoft's cursor icon file (CUR)
 * 
 * By drahoslav Bednář
 * github.com/drahoslav7
 *
 * license: WTFPL
 */

package main

import (
	"fmt"
	"image/png"
	"image/color"
	"os"
	"path/filepath"
	"flag"
	"bufio"
	"log"
)


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func main() {

	var (
		inputFileName, outputFileName string
		hotspot struct{
			x int
			y int
		}
		size int
		width, height int
	)


	/**
	 * parse args
	 */

	flag.Usage = func () {
		fmt.Fprintf(os.Stderr, "Usage: ")
		fmt.Fprintf(os.Stderr, "%s input [options]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  input\n")
		fmt.Fprintf(os.Stderr, "        input file name, should be .png\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&outputFileName, "o", "", "output file name (default same as input but .cur)")
	flag.IntVar(&hotspot.x, "x", 0, "horizontal hotspot position from left (default 0)")
	flag.IntVar(&hotspot.y, "y", 0, "vertical hotspot position from top (default 0)")
	flag.Parse()

	inputFileName = flag.Arg(0)
	if inputFileName == "" { 
		flag.Usage()
		return
	}
	if outputFileName == "" {
		outputFileName = inputFileName[0:len(inputFileName)-len(filepath.Ext(inputFileName))] + ".cur"
	}


	/**
	 * first get some info
	 */

	// get file size
	ifStat, err := os.Stat(inputFileName)
	check(err)
	size = int(ifStat.Size())

	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()
	
	// get image dimensions
	imageConfig, err := png.DecodeConfig(inputFile);
	check(err)
	width = imageConfig.Width
	height = imageConfig.Height

	//  print loaded info
	fmt.Printf("%s loaded (%dB)\n", ifStat.Name(), size)
	fmt.Printf("%d×%d px\n", width, height);

	// check format
	if (width > 256 || height > 256) {
		log.Fatal("Dimensions are too big, max 256×256 px")
	}
	if (imageConfig.ColorModel != color.RGBAModel && imageConfig.ColorModel != color.NRGBAModel) {
		log.Fatal("Color model is not RGBA")
	}


	/**
	 * prepare data to write
	 */

	// all little endian
	// sources of cur format specification:
	// https://en.wikipedia.org/wiki/ICO_(file_format)
	// https://www.daubnet.com/en/file-format-cur
	// http://www.blitzbasic.com/Community/post.php?topic=75770&post=847533
	icoHeader := []byte{
		0, 0, // must be 0
		2, 0, // type 2 = CUR (1 for ICO)
		1, 0, // number of icons = 1
	}
	directoryHeader := []byte{
		byte(width), // width of cursor
		byte(height), // height of cursor
		0, // Color count ( 0 = no color palette )
		0, // Reserved
		// note: next 2 lines has different meaning for ICO
		byte(hotspot.x), byte(hotspot.x>>8), // hotspot x from left
		byte(hotspot.y), byte(hotspot.y>>8), // hotspot y from top
		byte(size), byte(size>>8), byte(size>>16), byte(size>>24), // raw data size
		22+40, 0, 0, 0, // offset (of raw data) from start 
	}
	icoDataHeader := []byte{
		40, 0, 0, 0, // Size of this Header
		byte(width), byte(width>>8), byte(width>>16), byte(width>>24), // Width of image
		byte(height), byte(height>>8), byte(height>>16), byte(height>>24), // Height of image
		1, 0, // Planes
		32, 0, // BitCount per pixel, 0 = figure yourself
		5, 0, 0, 0, // Compression 5 = BI_PNG
		byte(size), byte(size>>8), byte(size>>16), byte(size>>24), // raw data size
		0, 0, 0, 0, // XpixelsPerM
		0, 0, 0, 0, // YpixelsPerM
		0, 0, 0, 0, // ColorsUsed
		0, 0, 0, 0, // ColorsImportant
	}


	/**
	 * let's write!
	 */

	inputFile.Seek(0,0)

	outputFile, err := os.Create(outputFileName)
	check(err)
	defer outputFile.Close()
	bufOutputFile := bufio.NewWriter(outputFile)

	bufOutputFile.Write(icoHeader)
	bufOutputFile.Write(directoryHeader)
	bufOutputFile.Write(icoDataHeader)

	bufOutputFile.ReadFrom(inputFile)

	bufOutputFile.Flush()


	/**
	 * Done
	 */

	fmt.Printf("%s created", outputFileName)
	
}
