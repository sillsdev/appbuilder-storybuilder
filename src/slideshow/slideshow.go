package slideshow

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"

	FFmpeg "github.com/sillsdev/appbuilder-storybuilder/src/ffmpeg"
	"github.com/sillsdev/appbuilder-storybuilder/src/helper"
)

/* Structure of a .slideshow
 * 	images: filepath strings to the images to be used for each slide
 *	audios: filepath strings to the narration audios to be used for each slide
 *	transitions: strings describing which (Xfade only) transition to be used in between each slide
 *	transitionDurations: strings describing the time (in milliseconds) for each transition to last
 *	timings: strings describing the time (in milliseconds) for each slide to last, also used for motions
 *	motions: arrays of floats describing the dimensions and positions for the start and end rectangles for zoom/pan effects
 *	templateName: string parsed from the .slideshow filename to be used for the final video product
 */

type slideshow struct {
	images              []string
	audioTracks         map[string]*FFmpeg.AudioTrack
	transitions         []string
	transitionDurations []string
	timings             []string
	motions             [][][]float64
	templateName        string
}

/* Function to create a new slideshow from a .slideshow template. The code parses the pieces out
 * and stores them in the slideshow struct
 *
 * Parameters:
 *			slideshowDirectory - the filepath to the .slideshow to be parsed
 *			v - verbose flag to determine what feedback to print
 * Returns:
 *			slideshow - the filled slideshow structure, containing all the data parsed
 */
func NewSlideshow(slideshowDirectory string, v bool) slideshow {
	slideshow_template := readSlideshowXML(slideshowDirectory)

	Images := []string{}
	Audios := make(map[string]*FFmpeg.AudioTrack)
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	templateDir, template_name := splitFileNameFromDirectory(slideshowDirectory)

	for i, slide := range slideshow_template.Slide {
		if v {
			fmt.Printf("slide[%d] = %+v\n", i, slide)
		}
		if slide.Audio.Background_Filename.Path != "" {
			// Intro music is stored differently in the xml
			filename := slide.Audio.Background_Filename.Path
			if v {
				fmt.Println("-- Background Audio: " + filename)
			}
			value, ok := Audios[filename]
			if ok {
				value.FrameCount += 1
			} else {
				value = &FFmpeg.AudioTrack{Filename: templateDir + filename, FrameStart: i, FrameCount: 1}
				Audios[filename] = value
			}
		}
		if slide.Audio.Filename.Name != "" {
			// Narration audio
			filename := slide.Audio.Filename.Name
			if v {
				fmt.Println("-- Narration Audio: " + filename)
			}
			value, ok := Audios[filename]
			if ok {
				value.FrameCount += 1
			} else {
				value = &FFmpeg.AudioTrack{Filename: templateDir + filename, FrameStart: i, FrameCount: 1}
				Audios[filename] = value
			}
		}
		Images = append(Images, templateDir+slide.Image.Name)
		if slide.Transition.Type == "" { // Default to a basic crossfade if no transition provided
			Transitions = append(Transitions, "fade")
		} else {
			Transitions = append(Transitions, slide.Transition.Type)
		}
		if slide.Transition.Duration == "" { // Default to 1000ms transition if none provided
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		}
		var motions = [][]float64{}
		if slide.Motion.Start == "" { // If no motion specified, default to a static "zoom/pan" effect
			motions = [][]float64{{0, 0, 1, 1}, {0, 0, 1, 1}}
		} else {
			motions = [][]float64{helper.ConvertStringToFloat(slide.Motion.Start), helper.ConvertStringToFloat(slide.Motion.End)}
		}
		Motions = append(Motions, motions)
		Timings = append(Timings, slide.Timing.Duration)
	}

	if v {
		fmt.Printf("Parsed %d images, %d audios, %d transitions, %d transition durations, %d timings, and %d motions, from %s\n",
			len(Images), len(Audios), len(Transitions), len(TransitionDurations), len(Timings), len(Motions), template_name)
	}
	slideshow := slideshow{Images, Audios, Transitions, TransitionDurations, Timings, Motions, template_name}

	fmt.Println("Parsing completed...")

	return slideshow
}

/* Function to scale all the input images depending on video quality
 * option to a uniform height/width to prevent issues in the video creation process.
 *
 * Parameters:
 *			lowQuality - specifies whether to generate a lower quality video by scaling the images to a smaller dimension
 */
func (s slideshow) ScaleImages(lowQuality bool) {
	width := "1280"
	height := "720"

	if lowQuality {
		println("-l specified, producing lower quality video")
		width = "852"
		height = "480"
	}

	totalNumImages := len(s.images)
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		go func(i int) {
			defer wg.Done()
			cmd := FFmpeg.CmdScaleImage(s.images[i], height, width, s.images[i])
			output, err := cmd.CombinedOutput()
			FFmpeg.CheckCMDError(output, err)
		}(i)
	}

	wg.Wait()
}

/* Function to create a video with all the data parsed from the .slideshow
 *
 * Parameters:
 *			useOldFade - specifies whether to use the old fade style instead of XFade, if desired
 *			tempDirectory - filepath to the temp folder to store the temporary videos created
 *			outputDirectory - filepath to the location to store the final completed video
 *			v - verbose flag to determine what feedback to print
 */
func (s slideshow) CreateVideo(useOldfade bool, tempDirectory string, outputDirectory string, v bool) {
	if v {
		fmt.Println("Temp Directory: " + tempDirectory)
		fmt.Println("Output Directory: " + outputDirectory)
	}
	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")
	var fadeType string = FFmpeg.ParseVersion()
	useXfade := fadeType == "X" && !useOldfade

	final_template_name := strings.TrimSuffix(s.templateName, ".slideshow")

	if useXfade {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		fmt.Println("Skipping MakeTemp and MergeTemp")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.timings, s.motions, tempDirectory, v)
		FFmpeg.MergeTempVideos(s.images, s.transitions, s.transitionDurations, s.timings, tempDirectory, v)
		FFmpeg.AddAudio(s.timings, s.audioTracks, tempDirectory, v)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.timings, s.motions, tempDirectory, v)
		FFmpeg.MergeTempVideosOldFade(s.images, s.transitionDurations, s.timings, tempDirectory, v)
		FFmpeg.AddAudio(s.timings, s.audioTracks, tempDirectory, v)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	}

	fmt.Println("Finished making video...")
}

// Helper function to generate an overlaid video of the software's result and a comparison video
func (s slideshow) CreateOverlaidVideo(finalVideoDirectory string, testVideoDirectory string, overlaidVideoDirectory string) {
	FFmpeg.CreateOverlaidVideoForTesting(finalVideoDirectory, testVideoDirectory, overlaidVideoDirectory)
}

/* Function to separate the .slideshow filename from the directory path
 *
 * Parameters:
 *		slideshowDirectory - path to the .slideshow file
 * Returns:
 *		template_directory - folder path leading up to the .slideshow file
 *		template_name - name of the .slideshow file
 */
func splitFileNameFromDirectory(slideshowDirectory string) (string, string) {
	var template_directory_split []string

	template_directory_split = regexp.MustCompile("[\\/\\\\]+").Split(slideshowDirectory, -1)

	template_directory := ""
	template_name := template_directory_split[len(template_directory_split)-1]

	if len(template_directory_split) == 1 {
		if runtime.GOOS != "windows" {
			template_directory = ""
		}
	} else {
		for i := 0; i < len(template_directory_split)-1; i++ {
			template_directory += template_directory_split[i] + "/"
		}
	}

	return template_directory, template_name
}
