package main

import (
	"log"

	ffmpeg "github.com/floostack/transcoder/ffmpeg"
)

//location of your repository
var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video"

//image name
var imageName = "VB-John 4v43-44.jpg"

//video name
var videoName = "video.mp4"

//location of where you downloaded FFmpeg
var baseFFmpegPath = "C:/FFmpeg"

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"

var inputPath = basePath + "/input/image/" + imageName
var outputPath = basePath + "/output/" + videoName

func main() {

	format := "mp4"
	overwrite := true

	opts := ffmpeg.Options{
		OutputFormat: &format,
		Overwrite:    &overwrite,
	}

	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   FfmpegBinPath,
		FfprobeBinPath:  FfprobeBinPath,
		ProgressEnabled: true,
	}

	progress, err := ffmpeg.
		New(ffmpegConf).
		Input(inputPath).
		Output(outputPath).
		WithOptions(opts).
		Start(opts)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range progress {
		log.Printf("%+v", msg)
	}

}
