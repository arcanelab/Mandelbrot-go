package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	width   = 1920
	height  = 1080
	maxIter = 20

	rMin = -3.0
	rMax = 0.5
	iMin = -1
	iMax = 1
)

type GradientTable []struct {
	color    colorful.Color
	position float64
}

var gradientData = GradientTable{
	{unpackColor("#000000"), 0.0},
	{unpackColor("#FF3A3A"), 0.25},
	{unpackColor("#8FFF42"), 0.5},
	{unpackColor("#ACB7FF"), 0.75},
	{unpackColor("#000000"), 1.0},
}

func gradient(t float64) colorful.Color {
	i := 0
	left := 0
	right := 1
	for ; i < len(gradientData)-1; i++ { // find between which colors t falls
		x0 := gradientData[i].position
		x1 := gradientData[i+1].position
		if x0 <= t && x1 >= t {
			break
		}
		left = i
		right = i + 1
	}

	// normalize t to the range between the gradient color points
	T := (t - gradientData[left].position) / (gradientData[right].position - gradientData[left].position)
	return gradientData[left].color.BlendLab(gradientData[right].color, T).Clamped()
}

func main() {
	fmt.Println("")

	output := image.NewNRGBA(image.Rect(0, 0, width, height))
	scale := width / (rMax - rMin)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mandout := mandebrot(complex(float64(x)/scale+rMin, float64(y)/scale+iMin))
			r, g, b := gradient(mandout).RGB255()
			output.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	savePng("mandelbrot.png", output)
}

func mandebrot(c complex128) float64 {
	z := 0 + 0i
	z += c
	var i uint = 0
	for ; i < maxIter; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 4 {
			break
		}
	}
	if i == maxIter {
		return cmplx.Abs(z) / 4
	}
	return float64(i) / maxIter
}

func savePng(path string, image image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	return png.Encode(file, image)
}

func unpackColor(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		fmt.Errorf("Error while creating color gradient")
	}
	return c
}
