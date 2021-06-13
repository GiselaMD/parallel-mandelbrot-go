package main

import (
	"image"
	"image/color"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Pix struct {
	x     int
	y     int
	color color.RGBA
}

const (
	posX   = -2
	posY   = -1.2
	height = 2.5

	imgWidth  = 1024
	imgHeight = 1024
	maxIter   = 3000
	samples   = 200

	ratio = float64(imgWidth) / float64(imgHeight)
)

var (
	img *image.RGBA
	win *pixelgl.Window
	err error
)

// func main() {
// 	log.Println("Initial processing...")
// 	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

// 	log.Println("Rendering...")

// 	render()

// 	log.Printf("Opening window size: [%v,%v]\n", imgWidth, imgHeight)

// 	if err := ebiten.Run(update, imgWidth, imgHeight, 1, "Test PAD"); err != nil {
// 		log.Println(err)
// 	}

// 	log.Println("Done!")
// }

func run() {
	log.Println("Initial processing...")
	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
	}

	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Rendering...")

	drawBuffer := make(chan Pix)
	render(drawBuffer)
	draw(drawBuffer)

	log.Println("Done!")

	for !win.Closed() {
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func draw(drawBuffer chan Pix) {
	for i := range drawBuffer {
		img.SetRGBA(i.x, i.x, i.color)
		// 	// img.SetRGBA(i.x, i.y, i.color)
		// 	// pic := pixel.PictureDataFromImage(img)
		// 	// sprite := pixel.NewSprite(pic, pic.Bounds())
		// 	// sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		// 	// win.Update()
		// 	// fmt.Println(i.color)
	}
}

func render(drawBuffer chan Pix) {
	//TODO: dividir em seçoes a partir da quantidade de threads passada por parametro
	for x := 0; x < imgWidth; x++ {
		// para cada coluna: 1 go routine (thread - eu acho que o go se gerencia para nao rodar mais threads do que o SO permite)
		// para cada ponto x,y processa a cor pelo mandelbrot
		go func(x int, drawBuffer chan Pix) {
			for y := 0; y < imgHeight; y++ {
				var colorR, colorG, colorB int
				for i := 0; i < samples; i++ {
					a := height*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + posX
					b := height*((float64(y)+RandFloat64())/float64(imgHeight)) + posY
					c := paint(mandelbrotIteraction(a, b, maxIter))
					colorR += int(c.R)
					colorG += int(c.G)
					colorB += int(c.B)
				}
				var cr, cg, cb uint8
				cr = uint8(float64(colorR) / float64(samples))
				cg = uint8(float64(colorG) / float64(samples))
				cb = uint8(float64(colorB) / float64(samples))

				// desenha o pixel processado na imagem
				// img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
				drawBuffer <- Pix{
					x: x, y: y, color: color.RGBA{R: cr, G: cg, B: cb, A: 255},
				}
			}

		}(x, drawBuffer)

	}
	// close(drawBuffer)
}

func paint(r float64, iter int) color.RGBA {
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
