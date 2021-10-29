package main

import (
	"log"
	"os/exec"
)

//location of your repository
//var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
// var basePath = "/Users/hyungyu/Documents/SIL-Video"	//hyungyu
var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video/" // david

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
		"-i", inputImagePath, // input image
		"-i", inputAudioPath, // input audio
		outputPath, // output
	)

	err := cmd.Start() // Start a process on another goroutine
	check(err)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)
}
