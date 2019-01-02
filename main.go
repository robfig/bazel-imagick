// bazel-imagick is for verifying that the imagick ImageMagick wrapper library
// is functional. The main method identifies an image given to it, and it
// defines some basic functionality that is exercised by a test suite.
package main // import "github.com/robfig/bazel-imagick"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beevik/etree"
	"github.com/robfig/imagick/imagick"
)

var identify = flag.String("identify", "", "path to an image to identify")

type IdentifyResult struct {
	ContentType string
	Width       int
	Height      int
	Quality     int
}

func main() {
	flag.Parse()
	if *identify == "" {
		flag.Usage()
		return
	}
	imagick.Initialize()
	wand := imagick.NewMagickWand()

	b, err := ioutil.ReadFile(*identify)
	fatalIf(err)

	err = wand.ReadImageBlob(b)
	fatalIf(err)

	result := IdentifyResult{
		ContentType: wand.GetImageFormat(),
		Width:       int(wand.GetImageWidth()),
		Height:      int(wand.GetImageHeight()),
		Quality:     int(wand.GetImageCompressionQuality()),
	}

	wand.Destroy()
	imagick.Terminate()

	fmt.Printf("%#v\n", result)
}

func fatalIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isPhotoSpherePhoto(xmp *etree.Document) bool {
	return findXMPElement(xmp, "ProjectionType") != nil ||
		findXMPAttr(xmp, "ProjectionType") != nil
}

func ResizeImage(wand *imagick.MagickWand, blob []byte, originalWidth, originalHeight, newWidth, newHeight int) ([]byte, error) {
	defer wand.Clear()

	err := wand.ReadImageBlob(blob)
	if err != nil {
		return nil, err
	}

	err = wand.ResizeImage(uint(newWidth), uint(newHeight), imagick.FILTER_LANCZOS, 1.0)
	if err != nil {
		return nil, err
	}

	for _, name := range wand.GetImageProfiles("*") {
		if name == "xmp" {
			data := wand.RemoveImageProfile(name)
			wand.SetImageProfile(name, data)
		}
	}

	return wand.GetImageBlob(), nil
}
