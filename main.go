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
	// First we parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := [][]string{}

	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData()
	for i, slide := range slideshow.Slide {
		if i == 0 {
			BackAudioPath = slide.Audio.Background_Filename.Path
			BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else {
			Audios = append(Audios, slide.Audio.Filename.Name)
		}
		Images = append(Images, slide.Image.Name)
		Transitions = append(Transitions, slide.Transition.Type)
		TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		temp := []string{slide.Timing.Start, slide.Timing.Duration}
		Timings = append(Timings, temp)
		fmt.Println(Timings[0][0])
	}
	fmt.Println("Combining temporary videos into single video...")
	combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)
	fmt.Println("Finished combining temporary videos...")
	addBackgroundMusic(BackAudioPath, BackAudioVolume)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}
func checkCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

func combineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	input_images := []string{}
	input_filters := ""
	totalNumImages := len(Images)
	concatTransitions := ""

	fmt.Println("Getting list of images and filters...")
	for i := 1; i <= totalNumImages; i++ {
		input_images = append(input_images, "-loop", "1", "-ss", Timings[i-1][0]+"ms", "-t", Timings[i-1][1]+"ms", "-i", basePath+"/input/"+Images[i-1])
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
	tempVol := 0.0
	if s, err := strconv.ParseFloat(backgroundVolume, 64); err == nil {
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
		"-y", basePath+"/output/finalvideo.mp4",
	)
	output, e := cmd.CombinedOutput()
	checkCMDError(output, e)
}
