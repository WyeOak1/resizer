package upload_api

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func resizeGIF(filePath string, w uint, h uint) error {
	const postFix string = "_resized"
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	ext = strings.ToLower(ext)
	imageFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()
	var gifImage *gif.GIF

	gifImage, err = gif.DecodeAll(imageFile)
	if err != nil {
		return err
	}

	tmpFileName := base[0:len(base)-len(ext)] + postFix + ext
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {

		return err

	}
	defer tmpFile.Close()
	if ext == ".gif" {
			for index, frame := range gifImage.Image {
				rect := frame.Bounds()
				tmpImage := frame.SubImage(rect)
				resizedImage := resize.Resize(w, h, tmpImage, resize.Lanczos3)
				//Add colors from original gif image
				var tmpPalette color.Palette
				for x := 1; x <= rect.Dx(); x++ {
					for y := 1; y <= rect.Dy(); y++ {
						if !contains(tmpPalette, gifImage.Image[index].At(x, y)) {
							tmpPalette = append(tmpPalette, gifImage.Image[index].At(x, y))
						}
					}
				}
				//After first image, image may contains only difference
				//bounds may not start from at (0,0)
				resizedBounds := resizedImage.Bounds()
				if index >= 1 {
					marginX := int(math.Floor(float64(rect.Min.X)))
					marginY := int(math.Floor(float64(rect.Min.Y)))
					resizedBounds = image.Rect(marginX, marginY, resizedBounds.Dx()+marginX,
						resizedBounds.Dy()+marginY)
				}
				resizedPalette := image.NewPaletted(resizedBounds, tmpPalette)
				draw.Draw(resizedPalette, resizedBounds, resizedImage, image.ZP, draw.Src)
				gifImage.Image[index] = resizedPalette
			}
			//Set size to resized size
			gifImage.Config.Width = int(w)
			gifImage.Config.Height = int(w)
		_ = gif.EncodeAll(tmpFile, gifImage)
		}
	return nil
}


//Check if color is already in the Palette
func contains(colorPalette color.Palette, c color.Color) bool {
	for _, tmpColor := range colorPalette {
		if tmpColor == c {

			return true

		}
	}
	return false
}
