package main

import (
	"log"
	"os/exec"

	"github.com/xfrr/goffmpeg/ffmpeg"
	"github.com/xfrr/goffmpeg/transcoder"
)

//location of your repository
var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
// var basePath = "/Users/hyungyu/Documents/SIL-Video"	//hyungyu

//image name
var imageName = "VB-John 4v43-44.jpg"

//audio name
var audioName = "inputs_mp3_44-JHNgul-01.mp3"

//video name
var videoName = "video.mp4"

//location of where you downloaded FFmpeg
var baseFFmpegPath = "C:/FFmpeg" //windows
// var baseFFmpegPath = "/usr/local/"	//mac

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"

var inputImagePath = basePath + "/input/image/" + imageName
var inputAudioPath = basePath + "/input/audio/" + audioName
var outputPath = basePath + "/output/" + videoName

func main() {
	convertToVideo()
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func convertToVideo() {
	cmd := exec.Command("ffmpeg",
		"-i", inputImagePath, // take stdin as input
		"-i", inputAudioPath, // strip out all (mostly) metadata
		outputPath,
	)

	err := cmd.Start() // Start a process on another goroutine
	check(err)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)
}

func imageToVideo() {
	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)
	trans.SetConfiguration(ffmpeg.Configuration{
		FfmpegBin:  FfmpegBinPath,
		FfprobeBin: FfprobeBinPath,
	})

	err := trans.Initialize(inputImagePath, outputPath)
	log.Println(err)

	trans.MediaFile().SetFrameRate(1)
	trans.MediaFile().SetVideoCodec("libx264")
	trans.MediaFile().SetOutputFormat("mp4")

	trans.Run(true)

	progress := trans.Output()

	// Example of printing transcoding progress
	for msg := range progress {
		log.Println(msg)
	}
}
