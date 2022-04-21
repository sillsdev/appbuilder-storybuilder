package slideshow

import (
	"encoding/xml"
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

func readSlideshowXML(filePath string) *slideshow_template {
	data, err := ioutil.ReadFile(filePath)
	helper.Check(err)

	slideshow_template := &slideshow_template{}
	_ = xml.Unmarshal([]byte(data), &slideshow_template)

	return slideshow_template
}
