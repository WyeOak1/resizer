package upload_api

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	parseError := r.ParseMultipartForm(20 * 1024 * 1024)
	if parseError != nil {
		fmt.Println(parseError)
	}
	file, handler, err := r.FormFile("MyFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	height := r.FormValue("height")
	width := r.FormValue("width")
	fmt.Println("height:  ", height )
	fmt.Println("width: ", width)

	contentType := handler.Header.Get("Content-Type")
	fmt.Println("File Info")
	fmt.Println("File Name: ", handler.Filename)
	fmt.Println("File Size: ", handler.Size)
	fmt.Println("Content-Type: ", contentType)

	fmt.Println(contentType)
	if contentType == "image/png"{

		tempFile, err := ioutil.TempFile("uploads", "upload-*.png")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		fileBytes, err1 := ioutil.ReadAll(file)
		if err1 != nil {
			fmt.Println(err1)
		}

		_, _ = tempFile.Write(fileBytes)
		fmt.Println("tempFile.Name:              ",tempFile.Name())
		fmt.Println("done")

		files, err := ioutil.ReadDir("./uploads")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			fmt.Println(f.Name())
		}
		resizePNGFile(tempFile.Name(),parsing(height),parsing(width))

	}
	if contentType == "image/jpg" || contentType == "image/jpeg"{
		tempFile, err := ioutil.TempFile("uploads", "upload-*.jpg")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		fileBytes, err1 := ioutil.ReadAll(file)
		if err1 != nil {
			fmt.Println(err1)
		}
		_, _ = tempFile.Write(fileBytes)
		fmt.Println("tempFile.Name:              ",tempFile.Name())
		fmt.Println("done")

		files, err := ioutil.ReadDir("./uploads")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			fmt.Println(f.Name())

		}
		resizeJPGFile(tempFile.Name(),parsing(height),parsing(width))
	}
	if contentType == "image/gif" {
		tempFile, err := ioutil.TempFile("uploads", "upload-*.gif")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		fileBytes, err1 := ioutil.ReadAll(file)
		if err1 != nil {
			fmt.Println(err1)
		}
		_, _ = tempFile.Write(fileBytes)
		fmt.Println("tempFile.Name:              ",tempFile.Name())
		fmt.Println("done")
		err11 := resizeGIF(tempFile.Name(), parsing(width), parsing(height))
		if err11 != nil {
			fmt.Println(err11)
		}
	}
}
func resizePNGFile(name string, height uint, width uint)  {
	src := openImage(name)
	out, err := os.Create("resize.png")
	m := resize.Resize(width, height, src, resize.Lanczos3)

	if err != nil {
		fmt.Println(err)
	}

	_ = png.Encode(out, m)
	out.Close()
}

func resizeJPGFile(name string, height uint, width uint)  {
	img := openImage(name)
	m := resize.Resize(width, height, img, resize.Lanczos3)
	out, err := os.Create("resize.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_ = jpeg.Encode(out, m, nil)
}

func parsing(num string)  uint {
	u64, err := strconv.ParseUint(num, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	number := uint(u64)
	return number
}


func openImage(pathFile string) image.Image {
	fl, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err) 
	}
	defer fl.Close()
	img, _, err := image.Decode(fl)
	if err != nil {
		log.Fatal(err)
	}
	return img
}