package main

import (
	"github.com/robfig/imagick/imagick"
)

const (
	formatBmp  string = "BMP"
	formatBmp3 string = "BMP3"
	formatGif  string = "GIF"
	formatJpeg string = "JPEG"
	formatPng  string = "PNG"
	formatTiff string = "TIFF"
	formatWebp string = "WEBP"

	TYPE_BMP  string = "image/bmp"
	TYPE_GIF  string = "image/gif"
	TYPE_JPEG string = "image/jpeg"
	TYPE_PNG  string = "image/png"
	TYPE_TIFF string = "image/tiff"
	TYPE_WEBP string = "image/webp"

	neutralGamma = 0.454545
)

func ConvertImageToWebp(wand *imagick.MagickWand, blob []byte, origType string) ([]byte, error) {
	defer wand.Clear()

	err := wand.ReadImageBlob(blob)
	if err != nil {
		return nil, err
	}
	err = wand.SetFormat(formatWebp)
	if err != nil {
		return nil, err
	}
	err = wand.SetImageCompressionQuality(75.0)
	if err != nil {
		return nil, err
	}
	err = wand.SetOption("webp:lossless", isWebpLossless(origType))
	if err != nil {
		return nil, err
	}
	return wand.GetImageBlob(), nil
}

func isWebpLossless(origType string) string {
	switch origType {
	case TYPE_BMP, TYPE_JPEG, TYPE_TIFF:
		return "false"
	case TYPE_GIF, TYPE_PNG:
		return "true"
	}
	return "true"
}
