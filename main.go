package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
//var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
// var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	// First we parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Motions := [][][]float64{}
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
		// Parse the zoom/pan data, which only exists for slides that are not the title and credits
		if slide.Motion.Start != "" {
			temp := [][]float64{convertStringToFloat(slide.Motion.Start), convertStringToFloat(slide.Motion.End)}
			Motions = append(Motions, temp)
			// fmt.Println(Motions)
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
	fmt.Println("Creating video...")
	combineVideos(Images, Transitions, TransitionDurations, Timings, Audios, Motions)
	fmt.Println("Finished making video...")
	fmt.Println("Adding intro music...")
	addBackgroundMusic(BackAudioPath, BackAudioVolume)
	fmt.Println("Video completed!")
}

/* Function to split the motion data into 4 pieces and convert them all to floats
 *  Parameters:
 *			stringData (string): The string that contains the four numerical values separated by spaces
 *  Returns:
 *			A float64 array with the four converted values
 */
func convertStringToFloat(stringData string) []float64 {
	floatData := []float64{}
	slicedStrings := strings.Split(stringData, " ")
	for _, str := range slicedStrings {
		flt, err := strconv.ParseFloat(str, 64)
		check(err)
		floatData = append(floatData, flt)
	}
	return floatData
}

func checkSign(num float64) float64 {

	//return true for negative
	//return false for positive
	result := math.Signbit(num)

	if result == true {
		num = -1
	} else {
		num = 1
	}

	return num
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

/** Function to create the video with all images + transitions
*	Parameters:
*		Images: ([]string) - Array of filenames for the images
 */
func combineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string, Motions [][][]float64) {
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

				// generate params for ffmpeg zoompan filter

				// in story buider, this is int variable(not float64).
				var num_frames float64 = 22.92 // HARD CODED JUST FOR NOW (540sec / 24FPS)

				var size_init float64 = Motions[i-1][0][3]
				var size_change float64 = Motions[i-1][1][3] - size_init
				var size_incr float64 = size_change / num_frames

				// var zoom_init float64 = 1.0 / Motions[i-1][0][3]
				// var zoom_change float64 = 1.0/Motions[i-1][1][3] - zoom_init
				// var zoom_incr = zoom_change / num_frames

				var x_init float64 = Motions[i-1][0][0]
				var x_end float64 = Motions[i-1][1][0]
				var x_change float64 = x_end - x_init
				var x_incr float64 = x_change / num_frames

				var y_init float64 = Motions[i-1][0][1]
				var y_end float64 = Motions[i-1][1][1]
				var y_change float64 = y_end - y_init
				var y_incr float64 = y_change / num_frames

				var zoom_cmd string = ""
				var x_cmd string = ""
				var y_cmd string = ""
				zoom_cmd += fmt.Sprintf("1/((%.1f)*%0.1f*(%.1f)*on)", size_init-size_incr, checkSign(size_incr), math.Abs(size_incr))
				x_cmd += fmt.Sprintf("%0.1f*iw*%0.1f*%0.1f*iw*on", x_init-x_incr, checkSign(x_incr), math.Abs(x_incr))
				y_cmd += fmt.Sprintf("%0.1f*ih*%0.1f*%0.1f*ih*on", y_init-y_incr, checkSign(y_incr), math.Abs(y_incr))

				// fmt.Println(Motions)
				// fmt.Println(size_init)
				// fmt.Println(size_change)
				// fmt.Println(size_incr)
				// fmt.Println(zoom_init)
				// fmt.Println(zoom_change)
				// fmt.Println(zoom_incr)
				// fmt.Println(zoom_cmd)
				// fmt.Println(x_cmd)
				// fmt.Println(y_cmd)

				input_filters += fmt.Sprintf(",zoompan=z='%s':x='%s':y='%s':d=%f:fps=24,fade=t=in:st=0:d=%dms,fade=t=out:st=%sms:d=%dms", zoom_cmd, x_cmd, y_cmd, num_frames, half_duration/2, Timings[i][1], half_duration/2)
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
