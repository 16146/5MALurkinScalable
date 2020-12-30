package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"sync"

	gim "github.com/ozankasikci/go-image-merge"
)

//Create a struct to deal with pixel

type Pixel struct {
	Point image.Point
	Color color.Color
}

var resp1 *http.Response
var err1 error
var resp2 *http.Response
var err2 error

// Keep it DRY so don't have to repeat opening file and decode
func OpenAndDecode(filepath string) (image.Image, string, error) {
	imgFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()
	img, format, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	return img, format, nil
}

// Decode image.Image's pixel data into []*Pixel
func DecodePixelsFromImage(img image.Image, offsetX, offsetY int) []*Pixel {
	pixels := []*Pixel{}
	for y := 0; y <= img.Bounds().Max.Y; y++ {
		for x := 0; x <= img.Bounds().Max.X; x++ {
			p := &Pixel{
				Point: image.Point{x + offsetX, y + offsetY},
				Color: img.At(x, y),
			}
			pixels = append(pixels, p)
		}
	}
	return pixels
}

//interactions with the slaves
func get(port int, wg *sync.WaitGroup) {
	if port == 8092 {
		resp1, err1 = http.Get("http://localhost:8092/get_mbrot")
		defer wg.Done()
	} else {
		resp2, err2 = http.Get("http://localhost:8093/get_mbrot")
		defer wg.Done()
	}
}

func main() {
	//go routine for the slaves
	var wg sync.WaitGroup
	wg.Add(1)
	go get(8092, &wg)
	wg.Add(1)
	go get(8093, &wg)
	wg.Wait()

	// accepts *Grid instances, grid unit count x, grid unit count y
	// returns an *image.RGBA object
	grids := []*gim.Grid{
		{ImageFilePath: "slave1/mandelbrot.png"},
		{ImageFilePath: "slave2/mandelbrot.png"},
	}
	rgba, err := gim.New(grids, 1, 2).Merge()

	// save the output to jpg or png
	file, err := os.Create("final_mandelbrot.png")
	err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 80})
	err = png.Encode(file, rgba)
	fmt.Printf(" %s\n", err)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}

}
