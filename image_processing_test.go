package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/beevik/etree"
	"github.com/robfig/imagick/imagick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDataBase      = "testdata"
	testDataInput     = filepath.Join(testDataBase, "input")
	testDataOutput, _ = ioutil.TempDir("", "photo_tests_output_")
	testDataExpected  = filepath.Join(testDataBase, "expected")
)

func init() {
	imagick.Initialize()
}

func TestSupportedFormats(t *testing.T) {
	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	var requiredFormats = []string{
		"BMP", "BMP3", "GIF", "JPEG", "JPG", "PNG", "TGA", "TIFF", "WEBP",
	}
	var actualFormats = wand.QueryFormats("*")

	for _, format := range requiredFormats {
		assert.Contains(t, actualFormats, format)
	}
}

func TestIdentifyImage(t *testing.T) {
	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	var tests = []struct {
		filename  string
		expectErr bool
		expected  IdentifyResult
	}{
		{"Landscape_1.jpg", false, IdentifyResult{formatJpeg, 600, 450, 91}},
		{"Portrait_1.jpg", false, IdentifyResult{formatJpeg, 450, 600, 91}},
		{"fruit.png", false, IdentifyResult{formatPng, 250, 250, 0}},
		{"lichtenstein.gif", false, IdentifyResult{formatGif, 403, 400, 0}},
		{"LENA512.BMP", false, IdentifyResult{formatBmp3, 512, 512, 0}},
		{"jello.tif", false, IdentifyResult{formatTiff, 256, 192, 0}},
		{"LICENSE", true, IdentifyResult{}},
		{"corrupt.BMP", true, IdentifyResult{}},
		{"corrupt-dessert.gif", true, IdentifyResult{}},
		{"corrupt-frowny.gif", true, IdentifyResult{}},
	}
	for _, test := range tests {
		blob, err := ioutil.ReadFile(filepath.Join(testDataInput, test.filename))
		require.Nil(t, err)

		err = wand.ReadImageBlob(blob)
		if test.expectErr {
			assert.NotNil(t, err)
			continue
		}
		if !assert.Nil(t, err) {
			continue
		}
		result := IdentifyResult{
			ContentType: wand.GetImageFormat(),
			Width:       int(wand.GetImageWidth()),
			Height:      int(wand.GetImageHeight()),
			Quality:     int(wand.GetImageCompressionQuality()),
		}
		assert.Equal(t, test.expected, result)
	}
}

func TestThumbnail(t *testing.T) {
	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	tests := []struct {
		filename                                           string
		originalWidth, originalHeight, newWidth, newHeight int
	}{
		{"thumbnail-Landscape_1.jpg", 600, 450, 200, 150},
		{"thumbnail-Portrait_1.jpg", 450, 600, 300, 400},
		{"thumbnail-profile.jpg", 1732, 828, 433, 207},
		{"thumbnail-photo_sphere.jpg", 5376, 2688, 2000, 1000},
		{"thumbnail-photo_sphere2.jpg", 2000, 1000, 1900, 950},
	}

	for _, test := range tests {
		transform := func(inputBytes []byte) ([]byte, error) {
			return ResizeImage(wand, inputBytes, test.originalWidth, test.originalHeight, test.newWidth, test.newHeight)
		}
		checkImageTransform(t, test.filename, test.filename, transform)
	}
	t.Log("Output directory: ", testDataOutput)
}

func TestIsPhotoSphere(t *testing.T) {

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	tests := []struct {
		filename      string
		isPhotoSphere bool
	}{
		{"thumbnail-photo_sphere.jpg", true},
		{"thumbnail-photo_sphere2.jpg", true},
	}

	for _, test := range tests {
		// Read the input photo
		inputFileName := filepath.Join(testDataInput, test.filename)
		inputBytes, err := ioutil.ReadFile(inputFileName)
		if !assert.Nil(t, err) {
			return
		}

		// Get the XMP doc
		err = wand.ReadImageBlob(inputBytes)
		if !assert.Nil(t, err) {
			return
		}
		data := wand.RemoveImageProfile("xmp")
		xmp := etree.NewDocument()
		err = xmp.ReadFromBytes(bytes.TrimRight(data, "\x00"))
		if !assert.Nil(t, err) || !assert.NotNil(t, xmp) {
			return
		}
		wand.Clear()

		assert.Equal(t, test.isPhotoSphere, isPhotoSpherePhoto(xmp))
	}
}

func TestConvertImageToWebp(t *testing.T) {

	tests := []struct {
		filename, outputFilename, origType string
	}{
		{"doll-gray.jpg", "doll-gray.webp", TYPE_JPEG},
		{"fruit.png", "fruit.webp", TYPE_PNG},
		{"jello.tif", "jello.webp", TYPE_TIFF},
		{"LENA512.BMP", "LENA512.webp", TYPE_BMP},
		{"tf_logo.gif", "tf_logo.webp", TYPE_GIF},
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	for _, test := range tests {
		transform := func(inputBytes []byte) ([]byte, error) {
			return ConvertImageToWebp(wand, inputBytes, test.origType)
		}
		checkImageTransform(t, test.filename, test.outputFilename, transform)
	}
	t.Log("Output directory: ", testDataOutput)
}

func checkImageTransform(t *testing.T, inFilename, outFilename string, transform func([]byte) ([]byte, error)) {

	inputFileName := filepath.Join(testDataInput, inFilename)
	outputFileName := filepath.Join(testDataOutput, outFilename)
	expectedFileName := filepath.Join(testDataExpected, outFilename)

	inputBytes, err := ioutil.ReadFile(inputFileName)
	if !assert.Nil(t, err) {
		return
	}

	outputBytes, err := transform(inputBytes)
	if err != nil {
		fmt.Println("TRANSFORM FIALED")
		os.Exit(1)
	}
	if !assert.Nil(t, err) {
		return
	}

	expectedBytes, err := ioutil.ReadFile(expectedFileName)
	if !assert.Nil(t, err) {
		ioutil.WriteFile(outputFileName, outputBytes, os.ModePerm)
		return
	}

	if !bytes.Equal(outputBytes, expectedBytes) {
		t.Error("Equality check failed for output file ", outFilename)
		ioutil.WriteFile(outputFileName, outputBytes, os.ModePerm)
	}
}
