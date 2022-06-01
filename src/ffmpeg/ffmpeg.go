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

/* Function to parse FFmpeg version from string and choose Xfade or traditional fade accordingly
 *
 * Returns:
 * 		The string returned from checkFFmpegVersion, either "X" or "F"
 */
func ParseVersion() string {
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
	return compareVersion(version)
}

/* Private function to check the ffmpeg version and choose Xfade or traditional fade accordingly
 *
 * Parameters:
 *		version - the version number to be checked
 * Returns:
		"X"(xfade) if version > 4.3.0, "F" (default fade) otherwise
*/
func compareVersion(version string) string {
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
 *		Images - Array of filenames for the images
 *		Timings - Array of timing duration for the audio for each image
 *		Motions - Array of start and end rectangles to use for the zoom/pan effects
 *		tempPath - Filepath to the temporary directory to store each temp video
 *		v - verbose flag to determine what feedback to print
 */
func MakeTempVideosWithoutAudio(Images []string, Timings []string, Motions [][][]float64, tempPath string, v bool) {
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
			zoom_cmd := CreateZoomCommand(Motions[i], helper.ConvertStringToFloat(duration)[0])
			if v {
				fmt.Println(fmt.Sprintf("Making temp%d-%d.mp4 video with:\n	Image: %s\n	Duration: %s ms\n	Start Rectangle (left, top, width, height): %f\n	End Rectangle (left, top, width, height): %f\n	Zoom Cmd: %s\n",
					i+1, totalNumImages, Images[i], duration, Motions[i][0], Motions[i][1], zoom_cmd))
			} else {
				fmt.Println(fmt.Sprintf("Making temp%d-%d.mp4 video", i+1, totalNumImages))
			}

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

/* Function to merge the temporary videos with transition filters between them
 * Parameters:
 *		Images - Array of filenames for the images
 *		Transitions - Array of Xfade transition names to use
 *		TransitionDurations - Array of durations for each transition
 *		Timings - array of timing duration for the audio for each image
 *		tempPath - path to the temp folder where the videos are stored
 *		v - verbose flag to determine what feedback to print
 */
func MergeTempVideos(Images []string, Transitions []string, TransitionDurations []string, Timings []string, tempPath string, v bool) {
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

		if v {
			fmt.Printf("%dth merge has transition %s and duration %f\n", i, transition, transition_duration)
		}
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
	CheckCMDError(output, err)
}

/** Merges the temporary videos using the old fade method with just plain crossfade transitions
 *
 *	Parameters:
 *		Images - Array of filenames for the images
 *		TransitionDurations - Array of durations for each transition
 *		Timings - array of timing duration for the audio for each image
 *		tempLocation - path to the temp folder where the videos are stored
 *		v - verbose flag to determine what feedback to print
 */
func MergeTempVideosOldFade(Images []string, TransitionDurations []string, Timings []string, tempLocation string, v bool) {
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

		if v {
			fmt.Printf("%dth merge has default fade transition and duration %f\n", i, transition_duration)
		}

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
	CheckCMDError(output, err)
}

/* Function to add the background and narration audio onto the video_with_no_audio.mp4
 * Parameters:
 *		Timings - array of timing duration for the audio for each image
 *		Audios - Array of filenames for the audios to be used
 *		tempPath - path to the temp folder where the audioless video is stored
 *		v - verbose flag to determine what feedback to print
 */
func AddAudio(Timings []string, Audios map[string]*AudioTrack, tempPath string, v bool) {
	fmt.Println("Adding audio...")
	audio_inputs := []string{}

	if v {
		fmt.Printf("Timings: %#v\n", Timings)
		fmt.Printf("Audios: %#v\n", Audios)
		fmt.Printf("TempPath: %s\n", tempPath)
	}

	audio_filter := ""
	audio_inputs = append(audio_inputs, "-y", "-i", tempPath+"/video_with_no_audio.mp4")

	i := 1
	for _, audioTrack := range Audios {
		audio_inputs = append(audio_inputs, "-i", audioTrack.Filename)
		//i var time_start int64
		time_start := 0.0
		// Sum up the durations from the start to the FrameStart to know when to start this audio
		for j := 0; j < audioTrack.FrameStart; j++ {
			//i transition_duration, err := strconv.ParseInt(strings.TrimSpace(Timings[j]), 10, 64)
			transition_duration, err := strconv.ParseFloat(strings.TrimSpace(Timings[j]), 64)
			helper.Check(err)
			time_start += transition_duration
		}
		var time_duration int64
		for k := 0; k < audioTrack.FrameCount; k++ {
			transition_duration, err := strconv.ParseInt(strings.TrimSpace(Timings[audioTrack.FrameStart+k]), 10, 64)
			helper.Check(err)
			time_duration += transition_duration
		}

		// apply filter for this audio input
		time_start = time_start / 1000.0
		if len(audio_filter) > 0 {
			audio_filter += ";"
		}
		audio_filter += fmt.Sprintf("[%d:a]atrim=start=0:duration=%dms,asetpts=expr=PTS+%f", i, time_duration, time_start)
		i++
	}

	audio_inputs = append(audio_inputs, "-filter_complex", audio_filter, "-map", "0:v", "-codec:v", "copy", "-codec:a", "libmp3lame", tempPath+"/merged_video.mp4")

	if v {
		println("Adding compiled audio to merged video and generating final result...")
		fmt.Printf("Command: %#v\n", audio_inputs)
	}
	cmd := exec.Command("ffmpeg", audio_inputs...)
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)

	trimEnd(tempPath)
}

/* Function to copy the final video from the temp folder to the output location specified
 * and change the filename
 *
 * Parameters:
 *		tempPath - path to the temp folder
 *		outputFolder - path to the folder to store the final result
 *		name - name to label the final video
 */
func CopyFinal(tempPath string, outputFolder string, name string) {
	// If -o is specified, save the final video at the specified location
	// Else save it to the folder of the executable
	var outputName string
	if len(outputFolder) > 0 {
		outputName = outputFolder + "/" + name + ".mp4"
	} else { // If -o is not specified, save the final video at the default location
		outputName = name + ".mp4"
	}

	fmt.Printf("Copying final video from temp folder to %s...\n", outputName)
	cmd := CmdCopyFile(tempPath+"/final.mp4", outputName)
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}

/* Function that creates an overlaid video between created video and testing video to see the differences between the two.
 *	One video is made half-transparent, changed to its negative image, and overlaid on the other video so that all similarities would cancel out and leave only the differences.
 * Parameters:
 *		finalVideoDirectory - folder where the final video produced is held
 *		trueVideo - file path to the comparison video
 *		destinationLocation - filepath to the folder to store the overlaid video
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

func ParseVideoLength(output string) float64 {
	re := regexp.MustCompile(`Duration: (?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2}\.\d{2})`)
	match := re.FindStringSubmatch(output)
	if match == nil {
		log.Fatal(match)
	}

	hour, err := strconv.ParseInt(string(match[1]), 10, 8)
	if err != nil {
		log.Fatal(err)
	}

	minute, err := strconv.ParseInt(string(match[2]), 10, 8)
	if err != nil {
		log.Fatal(err)
	}

	second, err := strconv.ParseFloat(string(match[3]), 8)
	if err != nil {
		log.Fatal(err)
	}

	return second + float64(60*minute) + float64(3600*hour)
}

func GetVideoLength(inputPath string) float64 {
	fmt.Println("File: " + inputPath)
	cmd := CmdGetVideoLength(inputPath)
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)

	return ParseVideoLength(string(output))
}
