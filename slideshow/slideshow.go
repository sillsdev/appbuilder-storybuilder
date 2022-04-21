package slideshow

import (
	"fmt"
	"strings"

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
}

func NewSlideshow(filePath string) slideshow {
	slideshow_template := readSlideshowXML(filePath)

	Images := []string{}
	Audios := []string{}
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	template_directory := removeFileNameFromDirectory(filePath)

	for _, slide := range slideshow_template.Slide {
		if slide.Audio.Background_Filename.Path != "" {
			Audios = append(Audios, template_directory+slide.Audio.Background_Filename.Path)
		} else {
			if slide.Audio.Filename.Name == "" {
				Audios = append(Audios, "")
			} else {
				Audios = append(Audios, template_directory+slide.Audio.Filename.Name)
			}
		}
		Images = append(Images, template_directory+slide.Image.Name)
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

	slideshow := slideshow{Images, Audios, Transitions, TransitionDurations, Timings, Motions}

	fmt.Println("Parsing completed...")

	return slideshow
}

func (s slideshow) ScaleImages(lowQuality *bool) {
	//Scaling images depending on video quality option
	if *lowQuality {
		FFmpeg.ScaleImages(s.images, "852", "480")
	} else {
		FFmpeg.ScaleImages(s.images, "1280", "720")
	}
}

func (s slideshow) CreateVideo(useOldfade *bool, tempDirectory string, outputDirectory string) {
	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")
	var fadeType string = FFmpeg.CheckVersion()

	useXfade := fadeType == "X" && !*useOldfade

	if useXfade {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.transitions, s.transitionDurations, s.timings, s.audios, s.motions, tempDirectory)
		FFmpeg.MergeTempVideos(s.images, s.transitions, s.transitionDurations, s.timings, tempDirectory)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory)
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.transitions, s.transitionDurations, s.timings, s.audios, s.motions, tempDirectory)
		FFmpeg.MergeTempVideosOldFade(s.images, s.transitionDurations, s.timings, tempDirectory)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory)
	}

	fmt.Println("Finished making video...")
}

func (s slideshow) CreateOverlaidVideo(testVideoDirectory string, finalVideoDirectory string) {
	FFmpeg.CreateOverlaidVideoForTesting(testVideoDirectory, finalVideoDirectory)
}

func removeFileNameFromDirectory(slideshowDirectory string) string {
	template_directory_split := strings.Split(slideshowDirectory, "/")
	template_directory := ""

	if len(template_directory_split) == 1 {
		template_directory = "./"
	} else {
		for i := 0; i < len(template_directory_split)-1; i++ {
			template_directory += template_directory_split[i] + "/"
		}
	}
	return template_directory
}
