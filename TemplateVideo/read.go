package main

import (
	"encoding/xml"
	"io/ioutil"
)

type Slideshow struct {
	Slide []Slide `xml:"slide"`
}

type Slide struct {
	Audio      Audio      `xml:"audio"`
	Image      Image      `xml:"image"`
	Motion     Motion     `xml:"motion"`
	Timing     Timing     `xml:"timing"`
	Transition Transition `xml:"transition"`
}

type Audio struct {
	Background          string              `xml:"background,attr"`
	Background_Filename Background_Filename `xml:"background-filename"`
	Filename            Filename            `xml:"filename"`
}

type Background_Filename struct {
	Volume string `xml:"volume,attr"`
	Path   string `xml:",chardata"`
}

type Filename struct {
	Name string `xml:",chardata"`
}

type Image struct {
	Name string `xml:",chardata"`
}

type Motion struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
}

type Timing struct {
	Start    string `xml:"start,attr"`
	End      string `xml:"end,attr"`
	Duration string `xml:"duration,attr"`
}

type Transition struct {
	Duration string `xml:"duration,attr"`
	Type     string `xml:",chardata"`
}

/* Function to parse xml data from the .slideshow file provided
 */
func readData(filePath string) *Slideshow {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	slideshow := &Slideshow{}

	_ = xml.Unmarshal([]byte(data), &slideshow)
	return slideshow
}
