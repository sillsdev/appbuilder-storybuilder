package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
//var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
//var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
var basePath = "C:/Users/Davideo/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	// First we parse in the various pieces from the template
	// Each piece is parsed as a string from the .slideshow xml and placed in a string array
	Images := []string{}
	Audios := []string{}
	BackAudioPath := "" // These two pieces are specific to the backgroundAudio that plays on top of the narration audio
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	// The timings are two part, a start and duration so they are placed in a 2D array
	Timings := [][]string{}
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData()
	for i, slide := range slideshow.Slide {
		if i == 0 { // The very first slide has background audio data
			BackAudioPath = slide.Audio.Background_Filename.Path
			BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else { // Every other slide has normal audio data
			Audios = append(Audios, slide.Audio.Filename.Name)
		}
		Images = append(Images, slide.Image.Name)
		Transitions = append(Transitions, slide.Transition.Type)
		if slide.Transition.Duration == "" { // If the slide doesn't have a transition, put a default
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		}
		temp := []string{slide.Timing.Start, slide.Timing.Duration}
		Timings = append(Timings, temp)
	}
	fmt.Println("Parsing completed...")
	fmt.Println("Scaling Images...")
	scaleImages(Images, "1500", "900") // Scales to a default 1500x900 but this can easily be made customizable by a command line option
	fmt.Println("Creating video...")
	combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)
	fmt.Println("Finished making video...")
	fmt.Println("Adding intro music...")
	addBackgroundMusic(BackAudioPath, BackAudioVolume)
	fmt.Println("Video completed!")
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

/** Function to scale all the images to a uniform height, width and aspect ratio to prevent errors in video creation
*	Parameters:
*			Images ([]string): An array of filepath strings to the images being scaled
*			height (string): The height (in pixels) to scale the images to
*			width (string): The width (in pixels) to scale the images to
*	Return:
*			None
 */
func scaleImages(Images []string, height string, width string) {
	for i := 0; i < len(Images); i++ {
		cmd := exec.Command("ffmpeg", "-i", basePath+"/input/"+Images[i],
			"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
			"-y", basePath+"/input/"+Images[i])
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}

/** Function to create the video with all images + transitions
*	Parameters:
*		Images ([]string): An array of filepath strings to the images being used
*		Transitions ([]string): An array of strings specifying the transitions to use for each image
*		TransitionDurations ([]string): An array of strings specifying the duration for each transition
*		Timings ([][]string): A 2D array of strings specifying the start and durations to use for pan/zoom effects
*		Audios ([]string): An array of filepath strings specifying the audio to be played for each image
*	Return:
*			None
 */
func combineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	input_images := []string{}
	input_filters := ""
	totalNumImages := len(Images)
	concatTransitions := ""

	fmt.Println("Getting list of images and filters...")
	for i := 0; i < totalNumImages; i++ {
		// Everything needs to be concatenated so always add the image to concatTransitions
		concatTransitions += fmt.Sprintf("[v%d]", i)
		// Everything needs to be cropped so add the crop filter to every image
		input_filters += fmt.Sprintf("[%d:v]crop=trunc(iw/2)*2:trunc(ih/2)*2", i)
		if i == totalNumImages-1 { // Credits image has no timings/transitions
			input_images = append(input_images, "-i", basePath+"/input/"+Images[i])
		} else {
			input_images = append(input_images, "-loop", "1", "-ss", Timings[i][0]+"ms", "-t", Timings[i][1]+"ms", "-i", basePath+"/input/"+Images[i])

			if i == 0 {
				input_filters += fmt.Sprintf(",fade=t=out:st=%sms:d=%sms", Timings[i][1], TransitionDurations[i])
			} else {
				half_duration, err := strconv.Atoi(TransitionDurations[i])
				check(err)
				input_filters += fmt.Sprintf(",fade=t=in:st=0:d=%dms,fade=t=out:st=%sms:d=%dms", half_duration/2, Timings[i][1], half_duration/2)
			}
		}
		input_filters += fmt.Sprintf("[v%d];", i)

	}

	concatTransitions += fmt.Sprintf("concat=n=%d:v=1:a=0,format=yuv420p[v]", totalNumImages)
	input_filters += concatTransitions

	input_images = append(input_images, "-i", basePath+"/input/narration-001.mp3",
		"-max_muxing_queue_size", "9999",
		"-filter_complex", input_filters, "-map", "[v]",
		"-map", fmt.Sprintf("%d:a", totalNumImages),
		"-shortest", "-y", basePath+"/output/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

/** Function to add background music once a video is finished
*	Parameters:
*		backgroundAudio (string): The filepath to the audio being used
*		backgroundVolume (string): The volume level to use for the background audio
*	Return:
*			None
 */
func addBackgroundMusic(backgroundAudio string, backgroundVolume string) {
	/*The volume value is originally a string and is a value between 0 and 100,
	* but ffmpeg needs a value between 0 and 1, so we use strconv.ParseFloat() to turn the string into
	* a float, divide by 100, then convert back to a string using fmt.Sprintf()
	 */
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
