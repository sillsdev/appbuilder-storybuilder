package ffmpeg_pkg

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/sillsdev/appbuilder-storybuilder/src/helper"
)

// Function to Check FFmpeg version and choose Xfade or traditional fade accordingly
func CheckVersion() string {
	cmd := CmdGetVersion()
	output, err := cmd.Output()
	CheckCMDError(output, err)

	re := regexp.MustCompile(`version (?P<num>\d+\.\d+(\.\d+)?)`) // Regular expression to fetch the version number, also made last number optional
	match := re.FindSubmatch(output)                              // Returns an array with the matching string, if found
	if match == nil {
		log.Fatal(match)
	}
	version := string(match[1]) // Get the string that holds the version number
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Version is %s\n", version)
	return version
}

// Function to Check FFmpeg version and choose Xfade or traditional fade accordingly
func checkFFmpegVersion(version string) string {
	char := []rune(version) // Convert the string "X.X.X" into a char array [X, ., X, ., X]
	num, _ := strconv.Atoi(string(char[0]))
	if num > 4 { // Version is > 4.x.x
		return "X"
	} else if num == 4 { // Version is 4.x.x
		num, _ = strconv.Atoi(string(char[2]))
		if num >= 3 { // Version is >= 4.3.x
			return "X"
		} else { // Version is < 4.3.x
			return "F"
		}
	} else { // Version is < 4.x.x
		return "F"
	}
}

/* Function to create temporary videos with the corresponding zoom filters for each slide without any audio
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func MakeTempVideosWithoutAudio(Images []string, Transitions []string, TransitionDurations []string, Timings []string, Audios []string, Motions [][][]float64, tempPath string) {
	fmt.Println("Making temporary videos in parallel...")
	totalNumImages := len(Images)

	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		go func(i int) {
			duration := "5000"

			if Timings[i] != "" {
				duration = Timings[i]
			}

			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer wg.Done()
			fmt.Printf("Making temp%d-%d.mp4 video\n", i, totalNumImages)
			zoom_cmd := CreateZoomCommand(Motions[i], helper.ConvertStringToFloat(duration))

			cmd := CmdCreateTempVideo(Images[i], duration, zoom_cmd, fmt.Sprintf(tempPath+"/temp%d-%d.mp4", i, totalNumImages))
			output, err := cmd.CombinedOutput()
			CheckCMDError(output, err)
		}(i)
	}

	// Wait for `wg.Done()` to be exectued the number of times
	//   specified in the `wg.Add()` call.
	// `wg.Done()` should be called the exact number of times
	//   that was specified in `wg.Add()`.
	wg.Wait()
}

/* Function to create the video with all images + transitions
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func CombineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings []string, Audios []string, Motions [][][]float64, tempPath string) {
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
			input_images = append(input_images, "-i", Images[i])
		} else {
			input_images = append(input_images, "-loop", "1", "-ss", "0ms", "-t", Timings[i]+"ms", "-i", Images[i])

			if i == 0 {
				input_filters += fmt.Sprintf(",fade=t=out:st=%sms:d=%sms", Timings[1], TransitionDurations[i])
			} else {
				half_duration, err := strconv.Atoi(TransitionDurations[i])
				helper.Check(err)
				// generate params for ffmpeg zoompan filter
				input_filters += CreateZoomCommand(Motions[i], helper.ConvertStringToFloat(Timings[i]))
				input_filters += fmt.Sprintf(",fade=t=in:st=0:d=%dms,fade=t=out:st=%sms:d=%dms", half_duration/2, Timings[i], half_duration/2)

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
		"-shortest", "-y", tempPath+"/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}

/* Function to merge the temporary videos with transition filters between them
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 */
func MergeTempVideos(Images []string, Transitions []string, TransitionDurations []string, Timings []string, tempPath string) {
	fmt.Println("Merging temporary videos...")
	video_fade_filter := ""
	settb := ""
	last_fade_output := "v0"

	totalNumImages := len(Images)

	video_total_length := 0.0
	video_each_length := make([]float64, totalNumImages)

	input_files := []string{}

	prev_offset := make([]float64, totalNumImages)
	prev_offset[0] = 0.0

	for i := 0; i < totalNumImages; i++ {

		input_files = append(input_files, "-i", fmt.Sprintf(tempPath+"/temp%d-%d.mp4", i, totalNumImages))

	}

	for i := 0; i < totalNumImages-1; i++ {
		transition := Transitions[i]
		transition_duration, err := strconv.ParseFloat(strings.TrimSpace(string(TransitionDurations[i])), 8)
		helper.Check(err)
		transition_duration = transition_duration / 1000

		//add time to the video that is sacrificied to xfade
		settb += fmt.Sprintf("[%d:v]tpad=stop_mode=clone:stop_duration=%f[v%d];", i, transition_duration/2, i)

		//get the current video length in seconds
		cmd := CmdGetVideoLength(fmt.Sprintf(tempPath+"/temp%d-%d.mp4", i, totalNumImages))

		output, err := cmd.CombinedOutput()
		CheckCMDError(output, err)

		//store the video length in an array
		video_each_length[i], err = strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

		//get the total video length of the videos combined thus far in seconds
		video_total_length += video_each_length[i]

		next_fade_output := fmt.Sprintf("v%d%d", i, i+1)

		if i < totalNumImages-2 {
			video_fade_filter += fmt.Sprintf("[%s][v%d]xfade=transition=%s:duration=%f:offset=%f", last_fade_output, i+1,
				transition, transition_duration, video_total_length)
		} else {
			video_fade_filter += fmt.Sprintf("[%s][%d:v]xfade=transition=%s:duration=%f:offset=%f", last_fade_output, i+1,
				transition, transition_duration, video_total_length)
		}

		last_fade_output = next_fade_output

		if i < totalNumImages-2 {
			video_fade_filter += fmt.Sprintf("[%s];", next_fade_output)
		} else {
			video_fade_filter += ",format=yuv420p"
		}

	}

	input_files = append(input_files, "-filter_complex", settb+video_fade_filter, "-y", tempPath+"/video_with_no_audio.mp4")

	cmd := exec.Command("ffmpeg", input_files...)

	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}

/** Merges the temporary videos using the old fade method with just plain crossfade transitions
 *
 *	Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 */
func MergeTempVideosOldFade(Images []string, TransitionDurations []string, Timings []string, tempLocation string) {
	fmt.Println("Merging temporary videos with traditional fade...")
	video_fade_filter := ""
	last_fade_output := ""
	settb := ""

	totalNumImages := len(Images)
	video_total_duration := 0.0
	video_total_length_minus_fade_transition := 0.0

	video_each_length := make([]float64, totalNumImages)

	input_files := []string{}

	prev_offset := make([]float64, totalNumImages)
	prev_offset[0] = 0.0

	for i := 0; i < totalNumImages; i++ {
		input_files = append(input_files, "-i", fmt.Sprintf(tempLocation+"/temp%d-%d.mp4", i, totalNumImages))
	}

	for i := 0; i < totalNumImages; i++ {
		transition_duration, err := strconv.ParseFloat(strings.TrimSpace(string(TransitionDurations[i])), 8)
		helper.Check(err)
		transition_duration = transition_duration / 1000

		//get the current video length in seconds
		cmd := CmdGetVideoLength(fmt.Sprintf(tempLocation+"/temp%d-%d.mp4", i, totalNumImages))
		output, err := cmd.CombinedOutput()
		CheckCMDError(output, err)

		video_each_length[i], err = strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

		//get the total video length of the videos combined thus far in seconds
		video_total_duration += video_each_length[i]

		video_total_length_minus_fade_transition = video_total_duration - transition_duration

		//add time to the video that is sacrificied to xfade
		settb += fmt.Sprintf("[%d:v]tpad=stop_mode=clone:stop_duration=%f[v%d];", i, transition_duration/2, i)

		if i == 0 {
			video_fade_filter += "[v0]setpts=PTS-STARTPTS[v_0];"
			last_fade_output += "[base][v_0]overlay[tmp1];"
		} else if i != totalNumImages-1 {
			video_fade_filter += fmt.Sprintf("[v%d]fade=in:st=0:d=%f:alpha=1,setpts=PTS-STARTPTS+((%f)/TB)[v_%d];",
				i, transition_duration, video_total_duration-video_each_length[i], i)

			last_fade_output += fmt.Sprintf("[tmp%d][v_%d]overlay[tmp%d];", i, i, i+1)
		} else {
			video_fade_filter += fmt.Sprintf("[v%d]fade=in:st=0:d=%f:alpha=1,setpts=PTS-STARTPTS+((%f)/TB)[v_%d];",
				i, transition_duration, video_total_duration-video_each_length[i], i)

			last_fade_output += fmt.Sprintf("[tmp%d][v_%d]overlay,format=yuv420p[fv]", i, i)
		}
	}

	setDimensions := fmt.Sprintf("color=black:%dx%d:d=%f[base];", 1280, 720, video_total_length_minus_fade_transition)

	input_files = append(input_files, "-filter_complex", setDimensions+settb+video_fade_filter+last_fade_output, "-map", "[fv]", "-y", tempLocation+"/video_with_no_audio.mp4")

	cmd := exec.Command("ffmpeg", input_files...)

	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}

/* Function to add the background and narration audio onto the video_with_no_audio.mp4
 * Parameters:
 *		Timings: ([]string) - array of timing duration for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func AddAudio(Timings []string, Audios []string, tempPath string) {
	fmt.Println("Adding audio...")
	audio_inputs := []string{}

	audio_filter := ""
	audio_last_filter := ""

	audio_inputs = append(audio_inputs, "-y", "-i", tempPath+"/video_with_no_audio.mp4")

	for i := 0; i < len(Audios); i++ {
		if Audios[i] != "" {
			audio_inputs = append(audio_inputs, "-i", Audios[i])
			totalDuration := 0.0

			for j := 0; j < i; j++ {
				if Audios[i] == Audios[j] {
					transition_duration, err := strconv.ParseFloat(strings.TrimSpace(Timings[j]), 8)
					helper.Check(err)
					transition_duration = transition_duration / 1000
					totalDuration += transition_duration
				}
			}

			//place the audio at the start of each slide
			audio_filter += fmt.Sprintf("[%d:a]atrim=start=%f:duration=%sms,asetpts=expr=PTS-STARTPTS[a%d];", i+1, totalDuration, strings.TrimSpace(Timings[i]), i+1)
			audio_last_filter += fmt.Sprintf("[a%d]", i+1)
		}
	}

	audio_last_filter += fmt.Sprintf("concat=n=%d:v=0:a=1[a]", len(Audios)-1)
	audio_filter += audio_last_filter

	audio_inputs = append(audio_inputs, "-filter_complex", audio_filter, "-map", "0:v", "-map", "[a]", "-codec:v", "copy", "-codec:a", "libmp3lame", tempPath+"/merged_video.mp4")

	cmd := exec.Command("ffmpeg", audio_inputs...)
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)

	trimEnd(tempPath)
}

func CopyFinal(tempPath string, outputFolder string, name string) {
	// If -o is specified, save the final video at the specified location

	var cmd *exec.Cmd

	if len(outputFolder) > 0 {
		cmd = CmdCopyFile(tempPath+"/final.mp4", outputFolder+"/"+name+".mp4")
	} else { // If -o is not specified, save the final video at the default location
		cmd = CmdCopyFile(tempPath+"/final.mp4", name+".mp4")
	}

	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}

/* Function that creates an overlaid video between created video and testing video to see the differences between the two.
	One video is made half-transparent, changed to its negative image, and overlaid on the other video so that all similarities would cancel out and leave only the differences.
 * Parameters:
 *		trueVideo: (string) - file path to the testing video
*/
func CreateOverlaidVideoForTesting(finalVideoDirectory string, trueVideo string, destinationLocation string) {
	outputDir := "./overlayVideo.mp4"
	if destinationLocation != "" {
		outputDir = destinationLocation + "/overlayVideo.mp4"
	}
	cmd := exec.Command("ffmpeg",
		"-i", finalVideoDirectory,
		"-i", trueVideo,
		"-filter_complex", "[1:v]format=yuva444p,lut=c3=128,negate[video2withAlpha],[0:v][video2withAlpha]overlay[out]",
		"-map", "[out]",
		"-y",
		outputDir,
	)

	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}