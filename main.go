package main

import (
	"log"

	ffmpeg "github.com/floostack/transcoder/ffmpeg"
)

//location of your repository
//var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video"	//windows
var basePath = "/Users/hyungyu/Documents/SIL-Video"	//mac

//image name
var imageName = "VB-John 4v43-44.jpg"

//audio name
var audioName = "inputs_mp3_44-JHNgul-01.mp3"

//video name
var videoName = "video.mp4"

//location of where you downloaded FFmpeg
//var baseFFmpegPath = "C:/FFmpeg"	//windows
var baseFFmpegPath = "/usr/local/"	//mac

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"	//both windows, mac
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"	//both windows, mac

//var inputPath = basePath + "/input/"
var inputImagePath = basePath + "/input/image/" + imageName
var inputAudioPath = basePath + "/input/audio/" + audioName
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
		//Input(inputPath).
		Input(inputImagePath).
		Output(tempPath).
		WithOptions(opts).
		Start(opts)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range progress {
		log.Printf("%+v", msg)
	}

}
