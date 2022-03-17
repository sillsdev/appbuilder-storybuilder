package main

import (
	"errors"
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
var location string

// Main function
func main() {
	// Create a temporary folder to store temporary files created when created a video
	os.Mkdir("./temp", 0755)

	// Ask the user for options
	var saveTemps = flag.Bool("s", false, "Include if user wishes to save temporary files created during production")
	flag.StringVar(&templateName, "t", "", "Specify template to use")
	var lowQuality = flag.Bool("l", false, "Include to produce a lower quality video (1280x720 => 852x480)")
	var changeOutput = flag.Bool("o", false, "Include if the user wants to save the final video to a specific location")
	flag.Parse()

	if *changeOutput {
		mydir, _ := os.Getwd()
		fmt.Println("Current working directory: " + mydir)
		fmt.Println("Enter output location: ")
		fmt.Scanln(&location)
		// https://freshman.tech/snippets/go/create-directory-if-not-exist/
		if _, err := os.Stat(location); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(location, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}

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
	var fadeType string = checkFFmpegVersion()

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
		allImages := makeTempVideosWithAudio(Images, Transitions, TransitionDurations, Timings, Audios)
		mergeVideos(allImages, Images, Transitions, TransitionDurations, Timings, 0)
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
	fmt.Printf(fmt.Sprintf("Time Taken: %.2f seconds\n", duration.Seconds()))
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
	if len(location) > 0 {
		cmd := exec.Command("ffmpeg", "-i", "./temp/merged0-0.mp4", "-y", location+"/final.mp4")
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	} else {
		cmd := exec.Command("ffmpeg", "-i", "./temp/merged0-0.mp4", "-y", "./final.mp4")
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}
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

// Function to Check FFmpeg version and choose Xfade or traditional fade accordingly
func checkFFmpegVersion() string {
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
	var result = ""
	char := []rune(version)

	intArr := []int{4, 3, 0}
	for i := 0; i < len(intArr); i++ {
		var temp = string(char[i])
		if temp == "." {
			break
		}
		num, version := strconv.Atoi(temp)

		if version != nil {
			return version.Error()
		}

		if intArr[i] > num {
			result = "F" // use old fade
			return result
		}
		result = "X" // use new fade
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

// Function to add audio to the temporary videos
func makeTempVideosWithAudio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) []int {
	totalNumImages := len(Images)

	cmd := exec.Command("")

	allImages := []int{}

	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines that are about to run concurrently
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		// Spawn a thread for each iteration in the loop
		// Pass 'i' into the goroutine's function in order to make sure each goroutine uses a different value for 'i'
		go func(i int) {
			// At the end of the goroutine, tell the WaitGroup that another thread has completed
			defer wg.Done()

			if Timings[i][0] == "" || Audios[i] == "" {
				fmt.Printf("Making temp%d-%d.mp4 video with empty audio\n", i, totalNumImages)
				cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", "3000ms", "-i", Images[i],
					"-f", "lavfi", "-i", "aevalsrc=0", "-t", "3000ms",
					"-shortest", "-pix_fmt", "yuv420p",
					"-y", fmt.Sprintf("./temp/temp%d-%d.mp4", i, totalNumImages))
			} else {
				fmt.Printf("Making temp%d-%d.mp4 video\n", i, totalNumImages)
				cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", Timings[i][1]+"ms", "-i", "./"+Images[i],
					"-ss", Timings[i][0]+"ms", "-t", Timings[i][1]+"ms", "-i", Audios[i],
					"-shortest", "-pix_fmt", "yuv420p", "-y", fmt.Sprintf("./temp/temp%d-%d.mp4", i, totalNumImages))

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
	wg.Wait()
	return allImages
}

// Function to merge all temporary videos together
// Product final video using merge() function
func mergeVideos(items []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {
	if len(items) < 2 {
		return items
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	first := []int{}

	go func() {
		defer wg.Done()
		first = mergeVideos(items[:len(items)/2], Images, Transitions, TransitionDurations, Timings, depth+1)
	}()

	second := mergeVideos(items[len(items)/2:], Images, Transitions, TransitionDurations, Timings, depth+1)

	wg.Wait()

	return merge(first, second, Images, Transitions, TransitionDurations, Timings, depth)
}

func merge(a []int, b []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {
	final := []int{}
	i := 0
	j := 0

	if len(a) == 1 && len(b) == 1 {
		//combine the individual temporary videos into merged files
		totalNumImages := len(Images)
		transition := Transitions[a[0]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[0]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[0]][1])
		check(err)
		offset := (float64(duration) - float64(transition_duration)) / 1000

		//check if video has full or partial audio
		cmd := exec.Command("ffprobe",
			"-i", fmt.Sprintf("./temp/temp%d-%d.mp4", a[0], totalNumImages),
			"-v", "error", "-of", "flat=s_",
			"-select_streams", "1", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1")

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

		audio_duration_offset, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		check(err)

		if offset-audio_duration_offset > 1 {
			//the calculated offset is more than 1 seconds longer than the true duration of the video
			offset = offset - float64(transition_duration/1000)
		}

		fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], totalNumImages, b[0], totalNumImages, transition, a[0], depth)

		cmd = exec.Command("ffmpeg",
			"-i", fmt.Sprintf("./temp/temp%d-%d.mp4", a[0], totalNumImages),
			"-i", fmt.Sprintf("./temp/temp%d-%d.mp4", b[0], totalNumImages),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("./temp/merged%d-%d.mp4", a[0], depth),
		)

		output, err = cmd.CombinedOutput()
		checkCMDError(output, err)
	} else if len(a) == 1 && len(b) == 2 {
		//if odd number of things to merge, then it merges a temporary video with merged file
		totalNumImages := len(Images)
		index := len(a) - 1
		newDepth := depth

		newDepth++

		if depth == newDepth {
			newDepth = a[index]
		}

		transition := Transitions[a[index]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[0]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[0]][1])
		offset := (float64(duration) - float64(transition_duration)) / 1000

		fmt.Printf("Combining videos temp%d-%d.mp4 and merged%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], totalNumImages, b[0], newDepth, transition, a[0], depth)

		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("./temp/temp%d-%d.mp4", a[0], totalNumImages),
			"-i", fmt.Sprintf("./temp/merged%d-%d.mp4", b[0], newDepth),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("./temp/merged%d-%d.mp4", a[0], depth),
		)
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	} else {
		//merging two merged files
		index := len(a) - 1
		newDepth := depth / 2

		if newDepth == 1 {
			newDepth = 2
		}
		if newDepth == 0 {
			newDepth = 1
		}

		depth++

		if depth == newDepth {
			newDepth = depth - 1
		}

		transition := Transitions[a[index]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[index]])
		check(err)

		transition_duration_float := float64(transition_duration) / 1000

		duration := 0

		for i := 0; i < len(a); i++ {
			duration_temp, err := strconv.Atoi(Timings[a[i]][1])
			check(err)

			duration += duration_temp
		}

		offset := (float64(duration) - float64(transition_duration*len(a))) / 1000

		fmt.Printf("Combining videos merged%d-%d.mp4 and merged%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], depth, b[0], depth, transition, a[0], newDepth)

		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("./temp/merged%d-%d.mp4", a[0], depth),
			"-i", fmt.Sprintf("./temp/merged%d-%d.mp4", b[0], depth),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("./temp/merged%d-%d.mp4", a[0], newDepth),
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
