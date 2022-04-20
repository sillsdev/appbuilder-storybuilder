package ffmpeg_pkg

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
)

type ffmpeg struct {
}

func NewFfmpeg() ffmpeg {
	return ffmpeg{}
}

// Function to Check FFmpeg version and choose Xfade or traditional fade accordingly
func (f ffmpeg) CheckFFmpegVersion() string {
	cmd := GetVersion()
	output, err := cmd.Output()
	checkCMDError(output, err)
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
	var result = ""
	char := []rune(version)

	intArr := []int{4, 3, 0} /// 4.3.0 = 4 3 0
	for i := 0; i < len(intArr); i++ {
		var temp = string(char[i])
		if temp == "." {
			break
		}
		num, err := strconv.Atoi(temp) // 4

		if err != nil {
			return err.Error()
		}

		if intArr[i] > num {
			result = "F" // use old fade
			return result
		}
		result = "X" // use new fade
	}
	return result
}

/* Function to scale all the input images to a uniform height/width
 * to prevent issues in the video creation process
 */
func (f ffmpeg) ScaleImages(Images []string, height string, width string) {
	totalNumImages := len(Images)
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		go func(i int) {
			defer wg.Done()
			ScaleImage(Images[i], height, width, Images[i])
		}(i)
	}

	wg.Wait()
}

/* Function to create temporary videos with the corresponding zoom filters for each slide without any audio
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func (f ffmpeg) MakeTempVideosWithoutAudio(Images []string, Transitions []string, TransitionDurations []string, Timings []string, Audios []string, Motions [][][]float64, tempPath string) {
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
			zoom_cmd := createZoomCommand(Motions[i], helper.ConvertStringToFloat(duration))

			CreateTempVideo(Images[i], duration, zoom_cmd, fmt.Sprintf(tempPath+"/temp%d-%d.mp4", i, totalNumImages))
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
func (f ffmpeg) CombineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings []string, Audios []string, Motions [][][]float64, tempPath string) {
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
				input_filters += createZoomCommand(Motions[i], helper.ConvertStringToFloat(Timings[i]))
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
	checkCMDError(output, err)
}

/* Function to merge the temporary videos with transition filters between them
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 */
func (f ffmpeg) MergeTempVideos(Images []string, Transitions []string, TransitionDurations []string, Timings []string, tempPath string) {
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
		video_each_length[i] = GetVideoLength(fmt.Sprintf(tempPath+"/temp%d-%d.mp4", i, totalNumImages))

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
	checkCMDError(output, err)
}

/** Merges the temporary videos using the old fade method with just plain crossfade transitions
 *
 *	Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([]string) - array of timing duration for the audio for each image
 */
func (f ffmpeg) MergeTempVideosOldFade(Images []string, TransitionDurations []string, Timings []string, tempLocation string) {
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
		video_each_length[i] = GetVideoLength(fmt.Sprintf(tempLocation+"/temp%d-%d.mp4", i, totalNumImages))

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
	checkCMDError(output, err)
}

/* Function to add the background and narration audio onto the video_with_no_audio.mp4
 * Parameters:
 *		Timings: ([]string) - array of timing duration for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func (f ffmpeg) AddAudio(Timings []string, Audios []string, tempPath string) {
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
	checkCMDError(output, err)

	trimEnd(tempPath)
}

func (f ffmpeg) CopyFinal(tempPath string, outputFolder string) {
	if len(outputFolder) > 0 {
		copyFile(tempPath+"/final.mp4", outputFolder+"/final.mp4")
	} else { // If -o is not specified, save the final video at the default location
		copyFile(tempPath+"/final.mp4", "./final.mp4")
	}

}

/* Function that creates an overlaid video between created video and testing video to see the differences between the two.
	One video is made half-transparent, changed to its negative image, and overlaid on the other video so that all similarities would cancel out and leave only the differences.
 * Parameters:
 *		trueVideo: (string) - file path to the testing video
*/
func (f ffmpeg) CreateOverlaidVideoForTesting(trueVideo string, destinationLocation string) {
	outputDir := "./overlayVideo.mp4"
	if destinationLocation != "" {
		outputDir = destinationLocation + "/overlayVideo.mp4"
	}
	cmd := exec.Command("ffmpeg",
		"-i", "./final.mp4",
		"-i", trueVideo,
		"-filter_complex", "[1:v]format=yuva444p,lut=c3=128,negate[video2withAlpha],[0:v][video2withAlpha]overlay[out]",
		"-map", "[out]",
		"-y",
		outputDir,
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func createZoomCommand(Motions [][]float64, TimingDuration []float64) string {
	num_frames := int(TimingDuration[0] / (1000.0 / 25.0))

	size_init := Motions[0][3]
	size_change := Motions[1][3] - size_init
	size_incr := size_change / float64(num_frames)

	var x_init float64 = Motions[0][0]
	var x_end float64 = Motions[1][0]
	var x_change float64 = x_end - x_init
	var x_incr float64 = x_change / float64(num_frames)

	var y_init float64 = Motions[0][1]
	var y_end float64 = Motions[1][1]
	var y_change float64 = y_end - y_init
	var y_incr float64 = y_change / float64(num_frames)

	zoom_cmd := fmt.Sprintf("1/((%.3f)%s(%.3f)*on)", size_init-size_incr, checkSign(size_incr), math.Abs(size_incr))
	x_cmd := fmt.Sprintf("%0.3f*iw%s%0.3f*iw*on", x_init-x_incr, checkSign(x_incr), math.Abs(x_incr))
	y_cmd := fmt.Sprintf("%0.3f*ih%s%0.3f*ih*on", y_init-y_incr, checkSign(y_incr), math.Abs(y_incr))
	final_cmd := fmt.Sprintf("scale=8000:-1,zoompan=z='%s':x='%s':y='%s':d=%d:fps=25,scale=1280:720,setsar=1:1", zoom_cmd, x_cmd, y_cmd, num_frames)

	// Test zoompan example from documentation (Zoom in up to 1.5x and pan always at center of picture)
	//final_cmd = "zoompan=z='min(zoom+0.0015,1.5)':d=700:x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)',scale=1500:900,setsar=1:1"

	return final_cmd
}

/* Function to trim the end of the video and remove excess empty audio when the audio file is longer than the video file
 */
func trimEnd(tempPath string) {
	fmt.Println("Trimming video...")

	video_length := GetVideoLength(tempPath + "/video_with_no_audio.mp4")

	//match the video length of the merged video with the true length of the video
	TrimLengthOfVideo(fmt.Sprintf("%f", video_length), tempPath)
}

func copyFile(to string, from string) {
	cmd := exec.Command("ffmpeg", "-i", to, "-y", from)
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func checkSign(num float64) string {
	result := math.Signbit(num)

	if result {
		return "-"
	} else {
		return "+"
	}
}
