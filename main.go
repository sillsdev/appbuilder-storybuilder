package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
var basePath = "C:/Users/sehee/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
//var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
//var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	// First we parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	//BackAudioPath := ""
	//BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := [][]string{}
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData()
	for i, slide := range slideshow.Slide {
		if i == 0 {
			//BackAudioPath = slide.Audio.Background_Filename.Path
			//BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else {
			Audios = append(Audios, slide.Audio.Filename.Name)
		}
		Images = append(Images, slide.Image.Name)
		Transitions = append(Transitions, slide.Transition.Type)
		if slide.Transition.Duration == "" {
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		}
		temp := []string{slide.Timing.Start, slide.Timing.Duration}
		Timings = append(Timings, temp)
	}
	fmt.Println("Parsing completed...")
	fmt.Println("Scaling Images...")
	scaleImages(Images, "1500", "900")
	fmt.Println("Creating temporary videos...")
	make_temp_videos(Images, Transitions, TransitionDurations, Timings, Audios)

	fmt.Println("Combining temporary videos with xfade")
	combine_xfade(Images, Transitions, TransitionDurations, Timings)

	//fmt.Println("Combining videos with traditional fade")
	//combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)

	fmt.Println("Finished making video...")
	fmt.Println("Adding audio...")
	addAudio(Timings)

	//fmt.Println("Adding intro music...")
	//addBackgroundMusic(BackAudioPath, BackAudioVolume)
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

func scaleImages(Images []string, height string, width string) {
	for i := 0; i < len(Images); i++ {
		cmd := exec.Command("ffmpeg", "-i", basePath+"/input/"+Images[i],
			"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
			"-y", basePath+"/input/"+Images[i])
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}

func make_temp_videos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	totalNumImages := len(Images)

	for i := 0; i < totalNumImages-1; i++ {
		fmt.Printf("Making temp%d.mp4 video\n", i)
		cmd := exec.Command("ffmpeg", "-loop", "1", "-i", basePath+"/input/"+Images[i],
			"-t", Timings[i][1]+"ms",
			"-shortest", "-pix_fmt", "yuv420p", "-y", fmt.Sprintf("%s/output/temp%d.mp4", basePath, i))

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}

	fmt.Printf("Making temp%d.mp4 video\n", totalNumImages-1)
	cmd := exec.Command("ffmpeg", "-loop", "1", "-t", "2000ms", "-i", basePath+"/input/"+Images[totalNumImages-1],
		"-shortest", "-pix_fmt", "yuv420p",
		"-y", fmt.Sprintf("%s/output/temp%d.mp4", basePath, totalNumImages-1))

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func combine_xfade(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	totalNumImages := len(Images)

	duration, err := strconv.Atoi(Timings[0][1])
	transition_duration, err := strconv.Atoi(TransitionDurations[0])
	check(err)
	prevOffset := duration - transition_duration

	transition := Transitions[0]

	if transition == "crossfade" {
		transition = "fade"
	}

	fmt.Printf("Combining vieos temp%d.mp4 and temp%d.mp4\n", 0, 1)
	cmd := exec.Command("ffmpeg",
		"-i", fmt.Sprintf("%s/output/temp%d.mp4", basePath, 0),
		"-i", fmt.Sprintf("%s/output/temp%d.mp4", basePath, 1),
		"-filter_complex", fmt.Sprintf("[0][1]xfade=transition=%s:duration=%dms:offset=%dms,format=yuv420p", transition, transition_duration, prevOffset),
		"-y", basePath+"/output/merged1.mp4",
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	for i := 1; i < totalNumImages-1; i++ {
		duration, err := strconv.Atoi(Timings[i][1])
		transition_duration, err := strconv.Atoi(TransitionDurations[i])
		check(err)
		offset := duration + prevOffset - transition_duration
		prevOffset = offset

		transition := Transitions[i]

		if transition == "crossfade" {
			transition = "fade"
			fmt.Printf("Combining videos merged%d.mp4 and temp%d.mp4 with %s transition. \n", i, i+1, transition)
			cmd := exec.Command("ffmpeg",
				"-i", fmt.Sprintf("%s/output/merged%d.mp4", basePath, i),
				"-i", fmt.Sprintf("%s/output/temp%d.mp4", basePath, i+1),
				"-filter_complex", fmt.Sprintf("[0][1]xfade=transition=%s:duration=%dms:offset=%dms,format=yuv420p", transition, transition_duration, offset),
				"-y", fmt.Sprintf("%s/output/merged%d.mp4", basePath, i+1),
			)

			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)
		} else if transition == "" {
			fmt.Printf("Combining videos merged%d.mp4 and temp%d.mp4 with no transition. \n", i, i+1)

			//ffmpeg -i input1.mp4 -c copy -bsf:v h264_mp4toannexb -f mpegts intermediate1.ts
			// ffmpeg -i input2.mp4 -c copy -bsf:v h264_mp4toannexb -f mpegts intermediate2.ts
			// ffmpeg -i "concat:intermediate1.ts|intermediate2.ts" -c copy -bsf:a aac_adtstoasc output.mp4

			cmd := exec.Command("ffmpeg",
				"-i", fmt.Sprintf("%s/output/merged%d.mp4", basePath, i),
				"-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", "-y", "intermediate1.ts",
			)

			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)

			cmd = exec.Command("ffmpeg",
				"-i", fmt.Sprintf("%s/output/temp%d.mp4", basePath, i+1),
				"-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", "-y", "intermediate2.ts",
			)

			output, err = cmd.CombinedOutput()
			checkCMDError(output, err)

			cmd = exec.Command("ffmpeg",
				"-i", "concat:intermediate1.ts|intermediate2.ts",
				"-c", "copy", "-bsf:a", "aac_adtstoasc", "-y", fmt.Sprintf("%s/output/merged%d.mp4", basePath, i+1),
			)

			output, err = cmd.CombinedOutput()
			checkCMDError(output, err)
		}
	}
}

func addAudio(Timings [][]string) {
	cmd := exec.Command("ffmpeg", "-i", basePath+"/output/merged4.mp4", "-i", basePath+"/input/narration-001.mp3",
		"-c:v", "copy", "-c:a", "aac", "-y", basePath+"/output/mergedVideo.mp4")

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
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
		"-i", basePath+"/input/"+backgroundAudio,
		"-filter_complex", "[1:0]volume="+fmt.Sprintf("%f", tempVol)+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		"-y", basePath+"/output/finalvideo.mp4",
	)
	output, e := cmd.CombinedOutput()
	checkCMDError(output, e)
}

/** Function to create the video with all images + transitions
*	Parameters:
*		Images: ([]string) - Array of filenames for the images
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
	fmt.Println(input_images)
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}
