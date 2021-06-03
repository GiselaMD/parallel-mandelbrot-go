package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	posX = -2
	posY = -1.2
	height = 2.5

	imgWidth     = 1024
	imgHeight    = 1024
	maxIter      = 2000
	samples      = 100

	showProgress = true
)

const (
	ratio = float64(imgWidth) / float64(imgHeight)
)

func main() {
	log.Println("Initial processing...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	log.Println("Rendering...")
	start := time.Now()
	render(img)
	end := time.Now()

	log.Println("Done rendering in", end.Sub(start))

	log.Println("Saving image...")
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func render(img *image.RGBA) {

	jobs := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func () {
			for y := range jobs {
				for x := 0; x < imgWidth; x++ {
					var colorR, colorG, colorB int
					for i := 0; i < samples; i++ {
						a := height * ratio * ((float64(x) + RandFloat64()) / float64(imgWidth)) + posX
						b := height * ((float64(y) + RandFloat64()) / float64(imgHeight)) + posY
						c := paint(mandelbrotIteraction(a, b, maxIter))
						colorR += int(c.R)
						colorG += int(c.G)
						colorB += int(c.B)
					}
					var cr, cg, cb uint8
						cr = uint8(float64(colorR) / float64(samples))
						cg = uint8(float64(colorG) / float64(samples))
						cb = uint8(float64(colorB) / float64(samples))
					img.SetRGBA(x, y, color.RGBA{ R: cr, G: cg, B: cb, A: 255 })
				}
			}
		}()
	}

	for y := 0; y < imgHeight; y++ {
		jobs <- y
		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", y, imgHeight, int(100*(float64(y) / float64(imgHeight))))
		}
	}
	if showProgress {
		fmt.Printf("\r%d/%[1]d (100%%)\n", imgHeight)
	}
}

func paint(r float64, iter int) color.RGBA {
	insideSet := color.RGBA{ R: 0, G: 0, B: 0, A: 255 }

	// validar se está dentro do conjunto
	// https://pt.wikipedia.org/wiki/Conjunto_de_Mandelbrot
	if r > 4 {
		return hslToRGB(float64(0.70) - float64(iter) / 3500 * r, 1, 0.5)
	}

	return insideSet
}

func mandelbrotIteraction(a, b float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x * x, y * y, x * y
		if xx + yy > 4 {
			return xx + yy, i
		}
		// xn+1 = x^2 - y^2 + a
		x = xx - yy + a
		// yn+1 = 2xy + b
		y = 2 * xy + b
	}

	return xx + yy, maxIter
}

