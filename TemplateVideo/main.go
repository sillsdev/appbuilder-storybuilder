package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var templateName string

// Main function
func main() {
	// Create a temporary folder to store temporary files created when created a video
	os.Mkdir("./temp", 0755)

	// Ask the user for options
	var saveTemps = flag.Bool("s", false, "Include if user wishes to save temporary files created during production")
	flag.StringVar(&templateName, "t", "", "Specify template to use")
	var lowQuality = flag.Bool("l", false, "Include to produce a lower quality video (1280x720 => 852x480)")
	// var outputLocation = flag.Bool("o", false, "Include if the user wants to save the final video to a specific location")
	flag.Parse()

	// Search for a template in local folder if no template is provided
	if templateName == "" {
		fmt.Println("No template provided, searching local folder...")
		filepath.WalkDir(".", findTemplate)
	}

	start := time.Now()

	// Parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := [][]string{}
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData(templateName)
	for _, slide := range slideshow.Slide {
		if slide.Audio.Background_Filename.Path != "" {
			Audios = append(Audios, slide.Audio.Background_Filename.Path)
			BackAudioPath = slide.Audio.Background_Filename.Path
			BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else {
			if slide.Audio.Filename.Name == "" {
				Audios = append(Audios, "")
			} else {
				Audios = append(Audios, slide.Audio.Filename.Name)
			}
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

	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")
	var versionString string = getVersion()
	var fadeType string = checkFFmpegVersion(versionString)

	// Scaling images depending on video quality option
	fmt.Println("Scaling images...")
	if *lowQuality {
		scaleImages(Images, "852", "480")
	} else {
		scaleImages(Images, "1280", "720")
	}

	fmt.Println("Creating video...")

	if fadeType == "X" {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		makeTempVideosWithoutAudio(Images, Transitions, TransitionDurations, Timings, Audios)
		MergeTempVideos(Images, Transitions, TransitionDurations, Timings)
		addAudio(Timings, Audios)
		copyFinal()
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)
		fmt.Println("Adding intro music...")
		addBackgroundMusic(BackAudioPath, BackAudioVolume)
	}

	fmt.Println("Finished making video...")

	// If user did not specify the -s flag at runtime, delete all the temporary videos
	if !*saveTemps {
		fmt.Println("-s not specified, removing temporary videos...")
		err := os.RemoveAll("./temp")
		check(err)
	}

	fmt.Println("Video production completed!")
	duration := time.Since(start)
	fmt.Printf("Time Taken: %.2f seconds\n", duration.Seconds())
}

// Function to check errors from non-CMD output
func check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}

// Function to check CMD error output when running commands
func checkCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

// Function to copy over the final video out of the main directory
func copyFinal() {
	cmd := exec.Command("ffmpeg", "-i", "./temp/final.mp4", "-y", "./final.mp4")
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

/* Function to scale all the input images to a uniform height/width
 * to prevent issues in the video creation process
 */
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

// Function to find the .slideshow template if none provided
func findTemplate(s string, d fs.DirEntry, err error) error {
	slideRegEx := regexp.MustCompile(`.+(.slideshow)$`) // Regular expression to find the .slideshow file
	if err != nil {
		return err
	}
	if slideRegEx.MatchString(d.Name()) {
		if templateName == "" {
			fmt.Println("Found template: " + s + "\nUsing found template...")
			templateName = s
		}
	}
	return nil
}

// Function to get FFmpeg version//
func getVersion() string {
	cmd := exec.Command("ffmpeg", "-version")
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
	return version
}

// Function to Check FFmpeg version and choose Xfade or traditional fade accordingly
func checkFFmpegVersion(version string) string {
	var result = ""
	char := []rune(version)
	intArr := []int{4, 3, 0}

	for i := 0; i < len(intArr); i++ {
		var temp = string(char[i])

		if temp == "." {
			continue
		}
		num, err := strconv.Atoi(temp)
		if err != nil {
			return err.Error()
		}
		//numbers we are comparing
		//fmt.Sprint(string(num) + " compared to: " + string(intArr[i]))
		if i == 0 && intArr[i] < num {
			result = "X" // use new old fade
			return result
		}
		if i == 1 && intArr[i] < num {
			result = "F" // use new old fade
			return result
		}

		if intArr[i] > num {
			result = "F" // use old fade
			return result
		}

		if intArr[i] < num {
			result = "X" // use new old fade
			return result
		}
	}
	return result
}

/* Function to create the video with all images + transitions
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([][]string) - 2-D array of timing data for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
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
		"-shortest", "-y", "./temp/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

// Function to add background music to the intro of the video at the end of the production process
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
		"-i", "./temp/mergedVideo.mp4",
		"-i", backgroundAudio,
		"-filter_complex", "[1:0]volume="+fmt.Sprintf("%f", tempVol)+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		"-y", "../finalvideo.mp4",
	)
	output, e := cmd.CombinedOutput()
	checkCMDError(output, e)
}

/* Function to create temporary videos with the corresponding zoom filters for each slide without any audio
 * Parameters:
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([][]string) - 2-D array of timing data for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func makeTempVideosWithoutAudio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	fmt.Println("Making temporary videos in parallel...")
	totalNumImages := len(Images)

	cmd := exec.Command("")

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

			fmt.Printf("Making temp%d-%d.mp4 video with empty audio\n", i, totalNumImages)
			cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", Timings[i][1]+"ms", "-i", Images[i],
				"-f", "lavfi", "-i", "aevalsrc=0", "-t", Timings[i][1],
				"-shortest", "-pix_fmt", "yuv420p",
				"-y", fmt.Sprintf("./temp/temp%d-%d.mp4", i, totalNumImages))

			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)
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
 *		Images: ([]string) - Array of filenames for the images
 *		Transitions: ([]string) - Array of Xfade transition names to use
 *		TransitionDurations: ([]string) - Array of durations for each transition
 *		Timings: ([][]string) - 2-D array of timing data for the audio for each image
 */
func MergeTempVideos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
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
		input_files = append(input_files, "-i", fmt.Sprintf("./temp/temp%d-%d.mp4", i, totalNumImages))
	}

	for i := 0; i < totalNumImages-1; i++ {
		transition := Transitions[i]
		transition_duration, err := strconv.ParseFloat(strings.TrimSpace(string(TransitionDurations[i])), 8)
		transition_duration = transition_duration / 1000

		//add time to the video that is sacrificied to xfade
		settb += fmt.Sprintf("[%d:v]tpad=stop_mode=clone:stop_duration=%f[v%d];", i, transition_duration/2, i)

		//get the current video length in seconds
		cmd := exec.Command("ffprobe",
			"-v", "error",
			"-show_entries", "format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1",
			fmt.Sprintf("./temp/temp%d-%d.mp4", i, totalNumImages),
		)
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

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
	input_files = append(input_files, "-filter_complex", settb+video_fade_filter, "-y", "./temp/video_with_no_audio.mp4")

	cmd := exec.Command("ffmpeg", input_files...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

/* Function to add the background and narration audio onto the video_with_no_audio.mp4
 * Parameters:
 *		Timings: ([][]string) - 2-D array of timing data for the audio for each image
 *		Audios: ([]string) - Array of filenames for the audios to be used
 */
func addAudio(Timings [][]string, Audios []string) {
	fmt.Println("Adding audio...")
	audio_inputs := []string{}

	audio_filter := ""
	audio_last_filter := ""
	audio_inputs = append(audio_inputs, "-y", "-i", "./temp/video_with_no_audio.mp4")

	for i := 0; i < len(Audios); i++ {
		if Audios[i] != "" {
			audio_inputs = append(audio_inputs, "-i", Audios[i])
			totalDuration := 0.0

			for j := 0; j < i; j++ {
				if Audios[i] == Audios[j] {
					transition_duration, err := strconv.ParseFloat(strings.TrimSpace(Timings[j][1]), 8)
					check(err)
					transition_duration = transition_duration / 1000
					totalDuration += transition_duration
				}
			}

			//place the audio at the start of each slide
			audio_filter += fmt.Sprintf("[%d:a]atrim=start=%f:duration=%sms,asetpts=expr=PTS-STARTPTS[a%d];", i+1, totalDuration, strings.TrimSpace(Timings[i][1]), i+1)
			audio_last_filter += fmt.Sprintf("[a%d]", i+1)
		}
	}

	audio_last_filter += fmt.Sprintf("concat=n=%d:v=0:a=1[a]", len(Audios))
	audio_filter += audio_last_filter

	audio_inputs = append(audio_inputs, "-filter_complex", audio_filter, "-map", "0:v", "-map", "[a]", "-codec:v", "copy", "-codec:a", "libmp3lame", "-shortest", "./temp/merged_video.mp4")

	cmd := exec.Command("ffmpeg", audio_inputs...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	TrimEnd()
}

/* Function to trim the end of the video and remove excess empty audio when the audio file is longer than the video file
 */
func TrimEnd() {
	fmt.Println("Trimming video...")
	//get the true length of the video
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"./temp/video_with_no_audio.mp4",
	)
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	video_length, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

	//match the video length of the merged video with the true length of the video
	cmd = exec.Command("ffmpeg",
		"-i", "./temp/merged_video.mp4",
		"-c", "copy", "-t", fmt.Sprintf("%f", video_length),
		"-y",
		"./temp/final.mp4",
	)

	output, err = cmd.CombinedOutput()
	checkCMDError(output, err)
}
