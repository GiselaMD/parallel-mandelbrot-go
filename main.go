package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type PixelColor struct {
	x  int
	y  int
	cr uint8
	cg uint8
	cb uint8
}

type Pix struct {
	x  int
	y  int
	cr uint8
	cg uint8
	cb uint8
}

const (
	posX   = -2
	posY   = -1.2
	height = 2.5

	imgWidth  = 1024
	imgHeight = 1024
	maxIter   = 1000
	samples   = 100

	numBlocks = 1936

	ratio = float64(imgWidth) / float64(imgHeight)
)

var (
	img      *image.RGBA
	mutex    sync.Mutex
	pixColor PixelColor
)

func run() {
	log.Println("Initial processing...")
	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	cfg := pixelgl.WindowConfig{
		Title:  "Parallel Mandelbrot - PAD",
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Rendering...")

	drawBuffer := make(chan Pix)
	// done := make(chan bool, 1)
	render(drawBuffer)
	go draw(drawBuffer, win)

	for !win.Closed() {
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}

	log.Println("Done!")

}

func main() {
	pixelgl.Run(run)
}

func draw(drawBuffer chan Pix, win *pixelgl.Window) {
	for i := range drawBuffer {
		img.SetRGBA(i.x, i.y, color.RGBA{R: i.cr, G: i.cg, B: i.cb, A: 255})

		// mutex.Lock()
		// // Lock para apenas uma goroutine por vez acessar o data
		// img.SetRGBA(pixColor.x, pixColor.y, color.RGBA{R: pixColor.cr, G: pixColor.cg, B: pixColor.cb, A: 255})
		// mutex.Unlock()
	}
}

func render(drawBuffer chan Pix) {
	//TODO: dividir em seçoes a partir da quantidade de threads passada por parametro
	var sqrt = int(math.Sqrt(numBlocks))
	for i := sqrt; i >= 0; i-- {
		for j := 0; j <= sqrt; j++ {

			// para cada coluna: 1 go routine (thread - eu acho que o go se gerencia para nao rodar mais threads do que o SO permite)
			// para cada ponto x,y processa a cor pelo mandelbrot

			// TODO: x inicial e x final, y inicial e y final pra mandar pras go routines

			go func(i, j int, drawBuffer chan Pix) {

				// matrix width
				// 1024 / 4

				// 0 - 256
				// 256 - 512
				// 512 - 768
				// 768 - 1023

				// --|--|--
				// --|--|--
				// --|--|--

				for x := i * (imgWidth / sqrt); x < (i+1)*(imgWidth/sqrt); x++ {
					for y := j * (imgHeight / sqrt); y < (j+1)*(imgHeight/sqrt); y++ {
						var colorR, colorG, colorB int
						for k := 0; k < samples; k++ {
							a := height*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + posX
							b := height*((float64(y)+RandFloat64())/float64(imgHeight)) + posY
							c := pixelColor(mandelbrotIteraction(a, b, maxIter))
							colorR += int(c.R)
							colorG += int(c.G)
							colorB += int(c.B)
						}
						var cr, cg, cb uint8
						cr = uint8(float64(colorR) / float64(samples))
						cg = uint8(float64(colorG) / float64(samples))
						cb = uint8(float64(colorB) / float64(samples))

						// mutex.Lock()
						// Lock para apenas uma goroutine por vez acessar o data
						// TODO: Criar matrix de resultado por pixel
						pixColor = PixelColor{x, y, cr, cg, cb}
						// mutex.Unlock()
						// done <- true

						drawBuffer <- Pix{
							x, y, cr, cg, cb,
						}
					}
				}
			}(i, j, drawBuffer)
		}
	}
}

func pixelColor(r float64, iter int) color.RGBA {
	insideSet := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// validar se está dentro do conjunto
	// https://pt.wikipedia.org/wiki/Conjunto_de_Mandelbrot
	if r > 4 {
		return hslToRGB(float64(0.70)-float64(iter)/3500*r, 1, 0.5)
	}

	return insideSet
}

func mandelbrotIteraction(a, b float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		// xn+1 = x^2 - y^2 + a
		x = xx - yy + a
		// yn+1 = 2xy + b
		y = 2*xy + b
	}

	return xx + yy, maxIter
}
