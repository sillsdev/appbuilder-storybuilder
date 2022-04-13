package slideshow

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
)

type slideshow_template struct {
	slide []slide `xml:"slide"`
}

type slide struct {
	audio      audio      `xml:"audio"`
	image      image      `xml:"image"`
	motion     motion     `xml:"motion"`
	timing     timing     `xml:"timing"`
	transition transition `xml:"transition"`
}

type audio struct {
	background          string              `xml:"background,attr"`
	background_Filename background_Filename `xml:"background-filename"`
	filename            filename            `xml:"filename"`
}

type background_Filename struct {
	volume string `xml:"volume,attr"`
	path   string `xml:",chardata"`
}

type filename struct {
	name string `xml:",chardata"`
}

type image struct {
	name string `xml:",chardata"`
}

type motion struct {
	start string `xml:"start,attr"`
	end   string `xml:"end,attr"`
}

type timing struct {
	start    string `xml:"start,attr"`
	end      string `xml:"end,attr"`
	duration string `xml:"duration,attr"`
}

type transition struct {
	duration string `xml:"duration,attr"`
	ttype    string `xml:",chardata"`
}

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

	for _, slide := range slideshow_template.slide {
		if slide.audio.background_Filename.path != "" {
			Audios = append(Audios, template_directory+slide.audio.background_Filename.path)
			BackAudioPath = slide.audio.background_Filename.path
			BackAudioVolume = slide.audio.background_Filename.volume
		} else {
			if slide.audio.filename.name == "" {
				Audios = append(Audios, "")
			} else {
				Audios = append(Audios, template_directory+slide.audio.filename.name)
			}
		}
		Images = append(Images, template_directory+slide.image.name)
		if slide.transition.ttype == "" {
			Transitions = append(Transitions, "fade")
		} else {
			Transitions = append(Transitions, slide.transition.ttype)
		}
		if slide.transition.duration == "" {
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.transition.duration)
		}
		var motions = [][]float64{}
		if slide.motion.start == "" {
			motions = [][]float64{{0, 0, 1, 1}, {0, 0, 1, 1}}
		} else {
			motions = [][]float64{helper.ConvertStringToFloat(slide.motion.start), helper.ConvertStringToFloat(slide.motion.end)}
		}

		Motions = append(Motions, motions)
		Timings = append(Timings, slide.timing.duration)
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

func (s slideshow) GetImages() []string {
	return s.images
}

func (s slideshow) GetAudios() []string {
	return s.audios
}
