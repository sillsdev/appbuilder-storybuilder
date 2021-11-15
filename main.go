package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
// var basePath = "/Users/hyungyu/Documents/SIL-Video"	//hyungyu
// var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david

//location of where you downloaded FFmpeg
var baseFFmpegPath = "C:/FFmpeg" //windows
// var baseFFmpegPath = "/usr/local/"	//mac

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"

var inputImagePath string
var inputAudioPath string
var inputFilePath = basePath + "/inputs.json"
var outputPath string

type InputTest struct {
	AudioLocation  string
	ImageLocation  string
	OutputLocation string
}

func main() {
	// First we read in the input file and parse the json
	data, err := ioutil.ReadFile(inputFilePath)
	check(err)
	var inputConfig InputTest
	err = json.Unmarshal(data, &inputConfig)
	check(err)
	// Set all the path vars
	inputAudioPath = inputConfig.AudioLocation
	inputImagePath = inputConfig.ImageLocation
	outputPath = inputConfig.OutputLocation
	// convertToVideo()
	readData()
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
