package slideshow

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"

	FFmpeg "github.com/sillsdev/appbuilder-storybuilder/ffmpeg"
	"github.com/sillsdev/appbuilder-storybuilder/helper"
)

type slideshow struct {
	images              []string
	audios              []string
	transitions         []string
	transitionDurations []string
	timings             []string
	motions             [][][]float64
	templateName        string
}

func NewSlideshow(slideshowDirectory string) slideshow {
	slideshow_template := readSlideshowXML(slideshowDirectory)

	Images := []string{}
	Audios := []string{}
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	templateDir, template_name := splitFileNameFromDirectory(slideshowDirectory)

	for _, slide := range slideshow_template.Slide {
		if slide.Audio.Background_Filename.Path != "" {
			Audios = append(Audios, templateDir+slide.Audio.Background_Filename.Path)
		} else {
			if slide.Audio.Filename.Name == "" {
				Audios = append(Audios, "")
			} else {
				Audios = append(Audios, templateDir+slide.Audio.Filename.Name)
			}
		}
		Images = append(Images, templateDir+slide.Image.Name)
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
		var motions = [][]float64{}
		if slide.Motion.Start == "" {
			motions = [][]float64{{0, 0, 1, 1}, {0, 0, 1, 1}}
		} else {
			motions = [][]float64{helper.ConvertStringToFloat(slide.Motion.Start), helper.ConvertStringToFloat(slide.Motion.End)}
		}

		Motions = append(Motions, motions)
		Timings = append(Timings, slide.Timing.Duration)
	}

	slideshow := slideshow{Images, Audios, Transitions, TransitionDurations, Timings, Motions, template_name}

	fmt.Println("Parsing completed...")

	return slideshow
}

/* Function to scale all the input images depending on video quality
 * option to a uniform height/width to prevent issues in the video creation process.
 */

func (s slideshow) ScaleImages(lowQuality *bool) {
	width := "1280"
	height := "720"

	if *lowQuality {
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

func (s slideshow) CreateVideo(useOldfade *bool, tempDirectory string, outputDirectory string) {
	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")
	var fadeType string = FFmpeg.CheckVersion()
	useXfade := fadeType == "X" && !*useOldfade

	final_template_name := strings.TrimSuffix(s.templateName, ".slideshow")

	if useXfade {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.transitions, s.transitionDurations, s.timings, s.audios, s.motions, tempDirectory)
		FFmpeg.MergeTempVideos(s.images, s.transitions, s.transitionDurations, s.timings, tempDirectory)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.transitions, s.transitionDurations, s.timings, s.audios, s.motions, tempDirectory)
		FFmpeg.MergeTempVideosOldFade(s.images, s.transitionDurations, s.timings, tempDirectory)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	}

	fmt.Println("Finished making video...")
}

func (s slideshow) CreateOverlaidVideo(testVideoDirectory string, finalVideoDirectory string) {
	FFmpeg.CreateOverlaidVideoForTesting(testVideoDirectory, finalVideoDirectory)
}

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
