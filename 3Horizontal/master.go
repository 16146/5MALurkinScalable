package main

import (
	"bytes"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	gim "github.com/ozankasikci/go-image-merge"
)

var slave1 *http.Response
var slave2 *http.Response

func poke(slave int, wg *sync.WaitGroup) {
	if slave == 1 {
		slave1, _ = http.Get("http://localhost:8081/poke")
		defer wg.Done()
	}
	if slave == 2 {
		slave2, _ = http.Get("http://localhost:8082/poke")
		defer wg.Done()
	}
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(2)
	go poke(1, &wg)
	go poke(2, &wg)
	wg.Wait()

	data1, _ := ioutil.ReadAll(slave1.Body)
	data2, _ := ioutil.ReadAll(slave2.Body)

	img1, _ := jpeg.Decode(bytes.NewReader(data1))
	img2, _ := jpeg.Decode(bytes.NewReader(data2))

	grids := []*gim.Grid{
		{Image: &img1},
		{Image: &img2},
	}
	rgba, err := gim.New(grids, 1, 2).Merge()

	output, err := os.Create("./final_output.jpg")

	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(output, rgba, &jpeg.Options{Quality: 80})

	elapsed := time.Since(start)
	log.Printf("Process mandelbrot took %s", elapsed)
	time.Sleep(time.Second)

}
