package slideshow

import (
	"fmt"
	"strings"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
	"github.com/gordon-cs/SIL-Video/Compiler/xml"
)

type slideshow struct {
	images              []string
	audios              []string
	backAudioPath       string
	backAudioVolume     string
	transitions         []string
	transitionDurations []string
	timings             []string
	motions             [][][]float64
}

func NewSlideshow(filePath string) slideshow {
	slideshow_template := xml.ReadSlideshowXML(filePath)

	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	template_directory := removeFileNameFromDirectory(filePath)

	for _, slide := range slideshow_template.Slide {
		if slide.Audio.Background_Filename.Path != "" {
			Audios = append(Audios, template_directory+slide.Audio.Background_Filename.Path)
			BackAudioPath = slide.Audio.Background_Filename.Path
			BackAudioVolume = slide.Audio.Background_Filename.Volume
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

	slideshow := slideshow{Images, Audios, BackAudioPath, BackAudioVolume, Transitions, TransitionDurations, Timings, Motions}

	fmt.Println("Parsing completed...")

	return slideshow
}

func (s slideshow) GetImages() []string {
	return s.images
}

func (s slideshow) GetAudios() []string {
	return s.audios
}

func (s slideshow) GetTransitions() []string {
	return s.transitions
}

func (s slideshow) GetTransitionDurations() []string {
	return s.transitionDurations
}

func (s slideshow) GetTimings() []string {
	return s.timings
}

func (s slideshow) GetMotions() [][][]float64 {
	return s.motions
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
