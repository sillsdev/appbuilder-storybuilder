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
	Images := []string{};
	Audios := []string{};
	BackAudioPath := "";
	BackAudioVolume := "";
	Transitions := []string{};
	TransitionDurations := []string{};
	Timings := [][]string{};

	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData();
	for i, slide := range slideshow.Slide {
		if (i == 0){
			BackAudioPath = slide.Audio.Background_Filename.Path;
			BackAudioVolume = slide.Audio.Background_Filename.Volume;
		} else {
			Audios = append(Audios, slide.Audio.Filename.Name);
		}
		Images = append(Images, slide.Image.Name);
		Transitions = append(Transitions, slide.Transition.Type);
		TransitionDurations = append(TransitionDurations, slide.Transition.Duration);
		temp := []string{slide.Timing.Start, slide.Timing.Duration};
		Timings = append(Timings, temp);
		fmt.Println(Timings[0][0]);
	}
	fmt.Println("Finished parsing .slideshow...")
	fmt.Println("Creating temporary videos...")
	createTempVideos(Images, Audios, Transitions, TransitionDurations, Timings);
	fmt.Println("Finished creating temporary videos...")
	fmt.Println("Fetching temporary video paths...")
	findVideos()
	fmt.Println("Finished fetching temporary video paths...")
	fmt.Println("Combining temporary videos into single video...")
	combineVideos()
	fmt.Println("Finished combining temporary videos...")
	addBackgroundMusic(BackAudioPath, BackAudioVolume)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}
func checkCMDError(output []byte, err error){
	if (err != nil) {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

func combineVideos() {
//   input_images := []string{}
// 	input_filters := ""
// 	totalNumImages := 3
// 	concatTransitions := ""
  
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", basePath+"/output/text.txt",
		"-y", basePath+"/output/mergedVideo.mp4",
	)

	output, e := cmd.CombinedOutput()
	checkCMDError(output,e)
}

func addBackgroundMusic(backgroundAudio string, backgroundVolume string) {
	// Convert the background volume to a number between 0 and 1
	tempVol := 0.0
	if s, err := strconv.ParseFloat(backgroundVolume, 64); err == nil {
        tempVol = s;
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
