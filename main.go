package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
//var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
//var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	// First we parse in the various pieces from the template
	var outputPath = "./output"
	var slideshow = readData()

	var titleimg = slideshow.Slide[0].Image.Name
	var img1 = slideshow.Slide[1].Image.Name
	var img2 = slideshow.Slide[2].Image.Name
	var img3 = slideshow.Slide[3].Image.Name
	var introAudio = slideshow.Slide[0].Audio.Background_Filename.Path
	var audio1 = slideshow.Slide[1].Audio.Filename.Name
	var title_start = slideshow.Slide[0].Timing.Start
	var title_duration = slideshow.Slide[0].Timing.Duration
	var img1_start = slideshow.Slide[1].Timing.Start
	var img1_duration = slideshow.Slide[1].Timing.Duration
	var img2_start = slideshow.Slide[2].Timing.Start
	var img2_duration = slideshow.Slide[2].Timing.Duration
	var img3_start = slideshow.Slide[3].Timing.Start
	var img3_duration = slideshow.Slide[3].Timing.Duration

	// Place them all inside a string slice
	paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration}

	createTempVideos(paths...)
	findVideos()
	combineVideos()
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func createTempVideos(paths ...string) {
	fmt.Println(paths)
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("ffmpeg",
			// "-i", fmt.Sprintf("%s/input/image-%d.jpg", basePath, i), // input image
			"-i", basePath+"/input/"+paths[i + 1],
			"-r", "30", // the framerate of the output video
			"-ss", paths[9+2*i-2]+"ms",
			"-t", paths[10+2*i-2]+"ms",
			"-i", basePath+"/input/narration-001.mp3", // input audio
			fmt.Sprintf("%s/output/output%d.mp4", basePath, i), // output
		)

		err := cmd.Start() // Start a process on another goroutine
		check(err)

		err = cmd.Wait() // wait until ffmpeg finishg
		check(err)
	}
}

func findVideos() {
	textfile, err := os.Create(basePath + "/output/text.txt")
	check(err)

	defer textfile.Close()

	files, err := ioutil.ReadDir(basePath + "/output")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mp4") {
			textfile.WriteString("file ")
			textfile.WriteString(file.Name())
			textfile.WriteString("\n")
		}
	}

	textfile.Sync()
}

func combineVideos() {
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", basePath+"/output/text.txt",
		basePath+"/output/mergedVideo.mp4",
	)

	err := cmd.Run() // Start a process on another goroutine
	check(err)
}
