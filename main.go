package main

import (
	"fmt"
	"log"
	"os/exec"
	//"os/exec"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
//var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
// var basePath = "/Users/hyungyu/Documents/SIL-Video"	//hyungyu
//var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

//location of where you downloaded FFmpeg
var baseFFmpegPath = "C:/FFmpeg" //windows
// var baseFFmpegPath = "/usr/local/"	//mac

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"

func main() {
	// First we read in the input file and parse the json
	//convertToVideo()

	// First we parse in the various pieces from the template
	var outputPath = "./output"
	var slideshow = readData()
	var titleimg = slideshow.Slide[0].Image.Name

	var img1 = slideshow.Slide[1].Image.Name

	var img2 = slideshow.Slide[2].Image.Name

	var img3 = slideshow.Slide[3].Image.Name

	var introAudio = slideshow.Slide[0].Audio.Background_Filename.Path

	var audio1 = slideshow.Slide[1].Audio.Filename.Name

	// Place them all inside a string slice
	paths := []string{outputPath, titleimg, img1, img2, img3, introAudio}
	// Using append, this can made variable for slides of any length/size
	paths = append(paths, audio1)
	// Pass our paths parameter to the convert function
	convertToVideo(paths...)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func convertToVideo(paths ...string) {
	// Here we can parse an individual element from paths
	fmt.Println(paths[0])
	// Here we can iterate through each element and access it
	for index, value := range paths {
		fmt.Println(index)
		fmt.Println(value)
	}

	cmd := exec.Command("ffmpeg",
		"-framerate", "1", // frame  to define how fast the pictures are read in, in this case, 1 picture per second
		"-i", "/Users/roddy/Desktop/SeniorProject/SIL-Video/image-%d.jpg", // input image
		"-r", "30", // the framerate of the output video
		"-i", "/Users/roddy/Desktop/SeniorProject/SIL-Video/narration-001.mp3", // input audio
		"/Users/roddy/Desktop/SeniorProject/SIL-Video/output/output.mp4", // output
	)

	err := cmd.Start() // Start a process on another goroutine
	check(err)

	err = cmd.Wait() // wait until ffmpeg finishg
	check(err)
}
