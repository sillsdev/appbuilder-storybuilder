package slideshow

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
)

type slideshow_template struct {
	Slide []slide `xml:"slide"`
}

type slide struct {
	Audio      audio      `xml:"audio"`
	Image      image      `xml:"image"`
	Motion     motion     `xml:"motion"`
	Timing     timing     `xml:"timing"`
	Transition transition `xml:"transition"`
}

type audio struct {
	Background          string              `xml:"background,attr"`
	Background_Filename background_Filename `xml:"background-filename"`
	Filename            filename            `xml:"filename"`
}

type background_Filename struct {
	Volume string `xml:"volume,attr"`
	Path   string `xml:",chardata"`
}

type filename struct {
	Name string `xml:",chardata"`
}

type image struct {
	Name string `xml:",chardata"`
}

type motion struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
}

type timing struct {
	Start    string `xml:"start,attr"`
	End      string `xml:"end,attr"`
	Duration string `xml:"duration,attr"`
}

type transition struct {
	Duration string `xml:"duration,attr"`
	Type     string `xml:",chardata"`
}

type slideshow struct {
	Images              []string
	Audios              []string
	BackAudioPath       string
	BackAudioVolume     string
	Transitions         []string
	TransitionDurations []string
	Timings             []string
	Motions             [][][]float64
}

func NewSlideshow(filePath string) slideshow {
	slideshow_template := readData(filePath)

	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	template_directory := helper.RemoveFileNameFromDirectory(filePath)

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

func readData(filePath string) *slideshow_template {
	data, err := ioutil.ReadFile(filePath)
	helper.Check(err)

	slideshow_template := &slideshow_template{}
	_ = xml.Unmarshal([]byte(data), &slideshow_template)
	return slideshow_template
}
