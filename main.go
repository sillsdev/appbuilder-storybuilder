package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
//var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
//var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	//First we parse in the various pieces from the template
	var outputPath = "./output"
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData()
	var titleimg = slideshow.Slide[0].Image.Name
	var img1 = slideshow.Slide[1].Image.Name
	var img2 = slideshow.Slide[2].Image.Name
	var img3 = slideshow.Slide[3].Image.Name
	var introAudio = slideshow.Slide[0].Audio.Background_Filename.Path
	var introVolume = slideshow.Slide[0].Audio.Background_Filename.Volume
	var audio1 = slideshow.Slide[1].Audio.Filename.Name
	var title_start = slideshow.Slide[0].Timing.Start
	var title_duration = slideshow.Slide[0].Timing.Duration
	var img1_start = slideshow.Slide[1].Timing.Start
	var img1_duration = slideshow.Slide[1].Timing.Duration
	var img2_start = slideshow.Slide[2].Timing.Start
	var img2_duration = slideshow.Slide[2].Timing.Duration
	var img3_start = slideshow.Slide[3].Timing.Start
	var img3_duration = slideshow.Slide[3].Timing.Duration

	// //Place them all inside a string slice
	paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration}
	fmt.Println("Finished parsing .slideshow...")

	combineVideos(paths...)
	addBackgroundMusic(introAudio, introVolume)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}

func combineVideos(paths ...string) {
	input_images := []string{}
	input_filters := ""
	totalNumImages := 3
	concatTransitions := ""

	fmt.Println("Getting list of images and filters...")
	for i := 1; i <= totalNumImages; i++ {
		input_images = append(input_images, "-loop", "1", "-ss", paths[9+2*i-2]+"ms", "-t", paths[10+2*i-2]+"ms", "-i", basePath+"/input/"+paths[i+1])
		concatTransitions += fmt.Sprintf("[v%d]", i-1)
		if i == 1 {
			input_filters += fmt.Sprintf("[0:v]crop=trunc(iw/2)*2:trunc(ih/2)*2,fade=t=out:st=%s:d=1000ms[v0];", paths[10])
		} else {
			input_filters += fmt.Sprintf("[%d:v]crop=trunc(iw/2)*2:trunc(ih/2)*2,fade=t=in:st=0:d=1000ms,fade=t=out:st=%sms:d=1000ms[v%d];", i-1, paths[10+2*i-2], i-1)
		}
	}

	concatTransitions += fmt.Sprintf("concat=n=%d:v=1:a=0,format=yuv420p[v]", totalNumImages)
	input_filters += concatTransitions

	input_images = append(input_images, "-i", basePath+"/input/narration-001.mp3", "-filter_complex", input_filters, "-map", "[v]",
		"-map", fmt.Sprintf("%d:a", totalNumImages),
		"-shortest", basePath+"/output/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
}

func addBackgroundMusic(backgroundAudio string, backgroundVolume string) {
	// Convert the background volume to a number between 0 and 1
	var tempVol = 0.0
	if s, err := strconv.ParseFloat(backgroundVolume, 32); err == nil {
		tempVol = s
	} else {
		fmt.Println("Error converting volume to float")
	}
	tempVol = tempVol / 100
	cmd := exec.Command("ffmpeg",
		"-i", basePath+"/output/mergedVideo.mp4",
		"-i", "./input/"+backgroundAudio,
		"-filter_complex", "[1:0]volume="+fmt.Sprintf("%f", tempVol)+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		basePath+"/output/finalvideo.mp4",
	)
	err := cmd.Run()
	check(err)
}
