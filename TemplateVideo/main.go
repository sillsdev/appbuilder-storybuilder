package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	filepath.WalkDir(".", findTemplate)
	var templateName string
	var fadeType string
	flag.StringVar(&templateName, "t", "./eng Visit of the Magi -Mat 2.1-23.slideshow", "Specify template to use.")
	flag.StringVar(&fadeType, "f", "", "Specify transition type (x) for xfade, leave blank for old fade")
	flag.Parse()
	if templateName == "" {
		log.Fatalln("Error, invalid template specified")
	}
	start := time.Now()
	// First we parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	//BackAudioPath := ""
	//BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := [][]string{}
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData(templateName)
	for i, slide := range slideshow.Slide {
		if i == 0 {
			Audios = append(Audios, slide.Audio.Background_Filename.Path)
			//BackAudioPath = slide.Audio.Background_Filename.Path
			//BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else {
			Audios = append(Audios, slide.Audio.Filename.Name)
		}
		Images = append(Images, slide.Image.Name)

		if slide.Transition.Type == "" {
			Transitions = append(Transitions, "fade")
		} else {
			Transitions = append(Transitions, slide.Transition.Type)
		}
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
	//	scaleImages(Images, "1500", "900")
	fmt.Println("Creating video...")

	//if using xfade
	//make_temp_videos(Images, Transitions, TransitionDurations, Timings, Audios)
	if fadeType == "xfade" {
		make_temp_videos_with_audio(Images, Transitions, TransitionDurations, Timings, Audios)
		combine_xfade_with_audio(Images, Transitions, TransitionDurations, Timings)

		//combine_xfade(Images, Transitions, TransitionDurations, Timings)
		//addAudio(Images)
	} else {
		//combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)
	}

	fmt.Println("Finished making video...")

	//fmt.Println("Adding intro music...")
	//addBackgroundMusic(BackAudioPath, BackAudioVolume)
	duration := time.Since(start)
	fmt.Println("Video completed!")
	fmt.Println(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))
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
		cmd := exec.Command("ffmpeg", "-i", "./"+Images[i],
			"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
			"-y", "./"+Images[i])
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}

func findTemplate(s string, d fs.DirEntry, err error) error {
	slideRegEx := regexp.MustCompile(`.+(.slideshow)$`)
	if err != nil {
		return err
	}
	if slideRegEx.MatchString(d.Name()) {
		fmt.Println(s)
	}
	return nil
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
			input_images = append(input_images, "-i", "./"+Images[i])
		} else {
			input_images = append(input_images, "-loop", "1", "-ss", Timings[i][0]+"ms", "-t", Timings[i][1]+"ms", "-i", "./"+Images[i])

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

	input_images = append(input_images, "-i", "./narration-001.mp3",
		"-max_muxing_queue_size", "9999",
		"-filter_complex", input_filters, "-map", "[v]",
		"-map", fmt.Sprintf("%d:a", totalNumImages),
		"-shortest", "-y", "../output/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func addAudio(Images []string) {
	totalNumImages := len(Images)
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("../output/merged%d.mp4", totalNumImages-2), "-i", "./narration-001.mp3",
		"-c:v", "copy", "-c:a", "aac", "-y", "../output/mergedVideo.mp4")

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func addBackgroundMusic(backgroundAudio string, backgroundVolume string) {
	tempVol := 0.0
	// Convert the background volume to a number between 0 and 1, if it exists
	if backgroundVolume != "" {
		if s, err := strconv.ParseFloat(backgroundVolume, 64); err == nil {
			tempVol = s
		} else {
			fmt.Println("Error converting volume to float")
		}
		tempVol = tempVol / 100
	} else {
		tempVol = .5
	}
	cmd := exec.Command("ffmpeg",
		"-i", "../output/mergedVideo.mp4",
		"-i", backgroundAudio,
		"-filter_complex", "[1:0]volume="+fmt.Sprintf("%f", tempVol)+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		"-y", "../output/finalvideo.mp4",
	)
	output, e := cmd.CombinedOutput()
	checkCMDError(output, e)
}

func make_temp_videos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	totalNumImages := len(Images)

	for i := 0; i < totalNumImages-1; i++ {
		fmt.Printf("Making temp%d.mp4 video\n", i)
		cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "./"+Images[i],
			"-t", Timings[i][1]+"ms",
			"-shortest", "-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d.mp4", i))

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}

	fmt.Printf("Making temp%d.mp4 video\n", totalNumImages-1)
	cmd := exec.Command("ffmpeg", "-loop", "1", "-t", "2000ms", "-i", "./"+Images[totalNumImages-1],
		"-shortest", "-pix_fmt", "yuv420p",
		"-y", fmt.Sprintf("../output/temp%d.mp4", totalNumImages-1))

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func make_temp_videos_with_audio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	totalNumImages := len(Images)

	cmd := exec.Command("")

	for i := 0; i < totalNumImages; i++ {
		fmt.Printf("Making temp%d-%d.mp4 video\n", i, totalNumImages)
		if Timings[i][0] == "" {
			cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", "3000ms", "-i", Images[i],
				"-f", "lavfi", "-i", "anullsrc", "-t", "3000ms",
				"-shortest", "-pix_fmt", "yuv420p",
				"-y", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages))
		} else {
			cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", Timings[i][1]+"ms", "-i", "./"+Images[i],
				"-ss", Timings[i][0]+"ms", "-t", Timings[i][1]+"ms", "-i", Audios[i],
				"-shortest", "-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages))
		}

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}

func combine_xfade_with_audio_faster(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	totalNumImages := len(Images)

	for totalNumImages != 1 {
		for i := 0; i < totalNumImages; i += 2 {
			transition_duration, err := strconv.Atoi(TransitionDurations[i])
			transition_duration_half := float32(transition_duration) * 0.75
			transition := Transitions[i]

			cmd := exec.Command("ffprobe", "-i",
				fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
				"-v", "quiet",
				"-show_entries", "format=duration",
				"-hide_banner", "-of", "default=noprint_wrappers=1:nokey=1")

			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)

			actual_duration, error := strconv.ParseFloat(strings.TrimSpace(string(output)), 32)
			check(error)

			length_of_video := 0

			for j := i * (len(Images) / totalNumImages); j < (i+1)*(len(Images)/totalNumImages); j++ {
				fmt.Println(totalNumImages, j)
				duration, err := strconv.Atoi(Timings[j][1])
				check(err)
				length_of_video += duration
			}

			offset := 0

			if int(transition_duration_half)*(len(Images)-totalNumImages) == 0 {

				offset = length_of_video - transition_duration
			} else {
				offset = length_of_video - transition_duration*(len(Images)/totalNumImages)
			}

			fmt.Println("offset: ", offset, "calculated length: ", length_of_video, "actual duration: ", actual_duration*1000)

			fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4 with %s transition to temp%d-%d.mp4. \n", i, totalNumImages, i+1, totalNumImages, transition, i/2, totalNumImages/2)

			if i == totalNumImages-2 && totalNumImages == len(Images) {
				cmd = exec.Command("ffmpeg",
					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
					"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%dms", transition, transition_duration, offset),
					"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i/2, totalNumImages/2),
				)
			} else {
				cmd = exec.Command("ffmpeg",
					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
					"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%dms;acrossfade=d=%d:o=0:c1=tri:c2=tri", transition, transition_duration, offset, transition_duration/1000),
					"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i/2, totalNumImages/2),
				)
			}

			output, err = cmd.CombinedOutput()
			checkCMDError(output, err)
		}
		totalNumImages /= 2
	}
}

func combine_xfade_with_audio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	totalNumImages := len(Images)
	totalDuration := 0

	duration, err := strconv.Atoi(Timings[0][1])
	totalDuration += duration
	transition_duration, err := strconv.Atoi(TransitionDurations[0])
	check(err)

	transition := Transitions[0]
	offset := duration/1000 - transition_duration/1000

	fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4\n", 0, totalNumImages, 1, totalNumImages)
	cmd := exec.Command("ffmpeg",
		"-i", fmt.Sprintf("../output/temp%d-%d.mp4", 0, totalNumImages),
		"-i", fmt.Sprintf("../output/temp%d-%d.mp4", 1, totalNumImages),
		"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%d[video];acrossfade=d=%d:o=0:c1=tri:c2=tri[audio]", transition, transition_duration, offset, transition_duration/1000),
		"-vsync", "0", "-map", "[video]",
		"-map", "[audio]",
		"-pix_fmt", "yuv420p", "-y", "../output/merged1.mp4",
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	for i := 1; i < totalNumImages-1; i++ {
		duration, err := strconv.Atoi(Timings[i][1])
		totalDuration += duration
		transition_duration, err := strconv.Atoi(TransitionDurations[i])
		transition := Transitions[i]

		cmd := exec.Command("ffprobe", "-i",
			fmt.Sprintf("../output/merged%d.mp4", i),
			"-v", "quiet",
			"-show_entries", "format=duration",
			"-hide_banner", "-of", "default=noprint_wrappers=1:nokey=1")

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

		actual_duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 32)
		check(err)

		offset := Round(actual_duration, 0.5) - float64(transition_duration/1000)
		offset = offset

		fmt.Println(actual_duration*1000, offset)

		fmt.Printf("Combining videos merged%d.mp4 and temp%d-%d.mp4 with %s transition. \n", i, i+1, totalNumImages, transition)
		cmd = exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/merged%d.mp4", i),
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
			"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%f[video];acrossfade=d=%d:o=0:c1=tri:c2=tri[audio]", transition, transition_duration, offset, transition_duration/1000),
			"-vsync", "0", "-map", "[video]",
			"-map", "[audio]",
			"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/merged%d.mp4", i+1),
		)

		output, err = cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

func combine_xfade(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	totalNumImages := len(Images)
	//totalDurations := 0

	duration, err := strconv.Atoi(Timings[0][1])
	transition_duration, err := strconv.Atoi(TransitionDurations[0])
	transition_duration_half := transition_duration / 2
	check(err)

	transition := Transitions[0]
	prevOffset := duration - transition_duration_half

	fmt.Printf("Combining vieos temp%d.mp4 and temp%d.mp4\n", 0, 1)
	cmd := exec.Command("ffmpeg",
		"-i", fmt.Sprintf("../output/temp%d.mp4", 0),
		"-i", fmt.Sprintf("../output/temp%d.mp4", 1),
		"-filter_complex", fmt.Sprintf("[0][1]xfade=transition=%s:duration=%dms:offset=%dms,format=yuv420p", transition, transition_duration, prevOffset),
		"-y", "../output/merged1.mp4",
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	for i := 1; i < totalNumImages-1; i++ {
		duration, err := strconv.Atoi(Timings[i][1])
		//start, err := strconv.Atoi(Timings[i][0])
		transition_duration, err := strconv.Atoi(TransitionDurations[i])
		transition_duration_half := transition_duration / 2
		transition := Transitions[i]

		check(err)
		offset := duration + prevOffset - transition_duration_half
		//fmt.Println(prevOffset)
		prevOffset = offset

		//fmt.Println(duration, offset, transition_duration)

		fmt.Printf("Combining videos merged%d.mp4 and temp%d.mp4 with %s transition. \n", i, i+1, transition)
		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/merged%d.mp4", i),
			"-i", fmt.Sprintf("../output/temp%d.mp4", i+1),
			"-filter_complex", fmt.Sprintf("[0][1]xfade=transition=%s:duration=%dms:offset=%dms,format=yuv420p", transition, transition_duration, offset),
			"-y", fmt.Sprintf("../output/merged%d.mp4", i+1),
		)

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
}
