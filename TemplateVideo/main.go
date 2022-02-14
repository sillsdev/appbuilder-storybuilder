package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func main() {
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
	//fmt.Println("Scaling Images...")
	//scaleImages(Images, "1500", "900")
	fmt.Println("Creating video...")

	//if using xfade
	if fadeType == "xfade" {
		//allImages := make_temp_videos_with_audio(Images, Transitions, TransitionDurations, Timings, Audios)
		allImages := []int{0, 1, 2, 3, 4, 5, 6, 7}
		mergeSort(allImages, Images, Transitions, TransitionDurations, Timings, 0)
		//combine_xfade_with_audio(Images, Transitions, TransitionDurations, Timings)
		//combine_xfade_with_audio_faster(Images, Transitions, TransitionDurations, Timings)
		//combine_xfade(Images, Transitions, TransitionDurations, Timings)
		//addAudio(Images)
	} else {
		//make_temp_videos(Images, Transitions, TransitionDurations, Timings, Audios)
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
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(len(Images))

	for i := 0; i < len(Images); i++ {
		go func(i int) {
			defer wg.Done()
			cmd := exec.Command("ffmpeg", "-i", "./"+Images[i],
				"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
				"-y", "./"+Images[i])
			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)
		}(i)
	}

	wg.Wait()
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

func make_temp_videos_with_audio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) []int {
	totalNumImages := len(Images)

	cmd := exec.Command("")

	allImages := []int{}

	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		// Spawn a thread for each iteration in the loop.
		// Pass 'i' into the goroutine's function
		//   in order to make sure each goroutine
		//   uses a different value for 'i'.
		go func(i int) {
			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer wg.Done()

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
		}(i)

		allImages = append(allImages, i)
	}

	// Wait for `wg.Done()` to be exectued the number of times
	//   specified in the `wg.Add()` call.
	// `wg.Done()` should be called the exact number of times
	//   that was specified in `wg.Add()`.
	// If the numbers do not match, `wg.Wait()` will either
	//   hang infinitely or throw a panic error.
	wg.Wait()
	fmt.Println(allImages)
	return allImages
}

// func combine_xfade_with_audio_faster(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
// 	totalNumImages := len(Images)

// 	for totalNumImages != 1 {
// 		for i := 0; i < totalNumImages; i += 2 {
// 			transition_duration, err := strconv.Atoi(TransitionDurations[i])
// 			transition_duration_half := float32(transition_duration) * 0.75
// 			transition := Transitions[i]

// 			cmd := exec.Command("ffprobe", "-i",
// 				fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
// 				"-v", "quiet",
// 				"-show_entries", "format=duration",
// 				"-hide_banner", "-of", "default=noprint_wrappers=1:nokey=1")

// 			output, err := cmd.CombinedOutput()
// 			checkCMDError(output, err)

// 			actual_duration, error := strconv.ParseFloat(strings.TrimSpace(string(output)), 32)
// 			check(error)

// 			length_of_video := 0

// 			for j := i * (len(Images) / totalNumImages); j < (i+1)*(len(Images)/totalNumImages); j++ {
// 				fmt.Println(totalNumImages, j)
// 				duration, err := strconv.Atoi(Timings[j][1])
// 				check(err)
// 				length_of_video += duration
// 			}

// 			offset := 0

// 			if int(transition_duration_half)*(len(Images)-totalNumImages) == 0 {

// 				offset = length_of_video - transition_duration
// 			} else {
// 				offset = length_of_video - transition_duration*(len(Images)/totalNumImages)
// 			}

// 			fmt.Println("offset: ", offset, "calculated length: ", length_of_video, "actual duration: ", actual_duration*1000)

// 			fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4 with %s transition to temp%d-%d.mp4. \n", i, totalNumImages, i+1, totalNumImages, transition, i/2, totalNumImages/2)

// 			if i == totalNumImages-2 && totalNumImages == len(Images) {
// 				cmd = exec.Command("ffmpeg",
// 					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
// 					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
// 					"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%dms", transition, transition_duration, offset),
// 					"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i/2, totalNumImages/2),
// 				)
// 			} else {
// 				cmd = exec.Command("ffmpeg",
// 					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
// 					"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
// 					"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%dms:offset=%dms;acrossfade=d=%d:o=0:c1=tri:c2=tri", transition, transition_duration, offset, transition_duration/1000),
// 					"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i/2, totalNumImages/2),
// 				)
// 			}

// 			output, err = cmd.CombinedOutput()
// 			checkCMDError(output, err)
// 		}
// 		totalNumImages /= 2
// 	}
// }

func combine_xfade_with_audio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	totalNumImages := len(Images)
	totalDuration := 0

	duration, err := strconv.Atoi(Timings[0][1])
	next_duration, err := strconv.Atoi(Timings[1][1])
	totalDuration += duration
	transition_duration, err := strconv.Atoi(TransitionDurations[0])
	check(err)

	transition := Transitions[0]
	//offset := duration/1000 - transition_duration/1000

	fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4\n", 0, totalNumImages, 1, totalNumImages)

	cmd := exec.Command("ffmpeg",
		"-i", fmt.Sprintf("../output/temp%d-%d.mp4", 0, totalNumImages),
		"-i", fmt.Sprintf("../output/temp%d-%d.mp4", 1, totalNumImages),
		"-filter_complex",
		fmt.Sprintf("xfade=transition=%s:duration=%d:offset=%d,format=yuv420p[video];[0:a]aresample=async=1:first_pts=0,apad,atrim=0:%d[A0];[1:a]aresample=async=1:first_pts=0,apad,atrim=0:%d[A1];[A0][A1]acrossfade=d=1:c1=tri:c2=tri[audio]", transition, transition_duration/1000, (duration-transition_duration)/1000, duration/1000, next_duration/1000),
		"-map", "[video]",
		"-map", "[audio]",
		"-movflags",
		"+faststart",
		"-y", "../output/merged1.mp4",
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	for i := 1; i < totalNumImages-1; i++ {
		duration, err = strconv.Atoi(Timings[i][1])
		next_duration, err = strconv.Atoi(Timings[i+1][1])
		totalDuration += duration
		transition_duration, err := strconv.Atoi(TransitionDurations[i])
		transition := Transitions[i]

		//C:\Users\sehee\scoop\shims\ffmpeg.exe -i ../output/temp0-9.mp4 -i ../output/temp1-9.mp4 -filter_complex xfade=transition=fade:duration=1:offset=0,format=yuv420p[video];[0:a]aresample=async=1:first_pts=0,apad,atrim=0:7[A0];[1:a]aresample=async=1:first_pts=0,apad,atrim=0:10[A1];[A0][A1]acrossfade=d=1:c1=tri:c2=tri[audio] -map [video] -map [audio] -movflags +faststart -y ../output/merged1.mp4
		//C:\Users\sehee\scoop\shims\ffmpeg.exe -i ../output/merged1.mp4 -i ../output/temp2-9.mp4 -filter_complex xfade=transition=wipeleft:duration=1:offset=16,format=yuv420p[video];[0:a]aresample=async=1:first_pts=0,apad,atrim=0:10[A0];[1:a]aresample=async=1:first_pts=0,apad,atrim=0:8[A1];[A0][A1]acrossfade=d=1:c1=tri:c2=tri[audio], -map [video] -map [audio] -movflags +faststart -pix_fmt yuv420p -y ../output/merged2.mp4
		cmd := exec.Command("ffprobe", "-i",
			fmt.Sprintf("../output/merged%d.mp4", i),
			"-v", "quiet",
			"-show_entries", "format=duration",
			"-hide_banner", "-of", "default=noprint_wrappers=1:nokey=1")

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

		//actual_duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 32)
		check(err)

		offset := (float64(totalDuration) - float64(transition_duration)) / 1000

		fmt.Printf("Combining videos merged%d.mp4 and temp%d-%d.mp4 with %s transition. \n", i, i+1, totalNumImages, transition)
		cmd = exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/merged%d.mp4", i),
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", i+1, totalNumImages),
			"-filter_complex", fmt.Sprintf("xfade=transition=%s:duration=%d:offset=%d,format=yuv420p[video];[0:a]aresample=async=1:first_pts=0,apad,atrim=0:%d[A0];[1:a]aresample=async=1:first_pts=0,apad,atrim=0:%d[A1];[A0][A1]acrossfade=d=1:c1=tri:c2=tri[audio]", transition, transition_duration/1000, int(offset), duration/1000, next_duration/1000),
			"-map", "[video]",
			"-map", "[audio]",
			"-movflags",
			"+faststart",
			"-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/merged%d.mp4", i+1),
		)

		output, err = cmd.CombinedOutput()
		checkCMDError(output, err)
	}
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

func mergeSort(items []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {
	if len(items) < 2 {
		return items
	}
	first := mergeSort(items[:len(items)/2], Images, Transitions, TransitionDurations, Timings, depth+1)
	second := mergeSort(items[len(items)/2:], Images, Transitions, TransitionDurations, Timings, depth+1)

	return merge(first, second, Images, Transitions, TransitionDurations, Timings, depth)
}

func merge(a []int, b []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {

	final := []int{}
	i := 0
	j := 0

	if len(a) == 1 && len(b) == 1 {

		totalNumImages := len(Images)
		transition := Transitions[a[0]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[0]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[0]][1])
		offset := (float64(duration) - float64(transition_duration)) / 1000

		fmt.Println(transition, transition_duration_float, offset)

		fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], totalNumImages, b[0], totalNumImages, transition, a[0], a[len(a)-1])
		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", a[0], totalNumImages),
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", b[0], totalNumImages),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=0.25:o=0[outa]", transition, transition_duration_float, offset),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("../output/merged%d.mp4", a[0]),
		)

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	} else {
		fmt.Println(a, b)

		index := len(a) - 1
		transition := Transitions[a[index]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[index]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[index]][1])
		offset := (float64(duration) - float64(transition_duration)) / 1000

		fmt.Printf("Combining videos merged%d.mp4 and merged%d.mp4 with fade transition to merged%d-%d. \n", a[0], b[0], a[0], index)

		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/merged%d.mp4", a[0]),
			"-i", fmt.Sprintf("../output/merged%d.mp4", b[0]),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=0.25:o=0[outa]", transition, transition_duration_float, offset),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("../output/merged%d.mp4", a[0]),
		)

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

	}

	for i < len(a) && j < len(b) {

		final = append(final, a[i])
		i++

	}

	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}
