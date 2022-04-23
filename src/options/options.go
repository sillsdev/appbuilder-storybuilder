package options

import (
	"flag"
)

type options struct {
	SlideshowDirectory    string
	OutputDirectory       string
	TemporaryDirectory    string
	OverlayVideoDirectory string
	LowQuality            bool
	SaveTemps             bool
	UseOldFade            bool
}

func ParseFlags() options {
	var slideshowDirectory string
	var outputDirectory string
	var temporaryDirectory string
	var overlayVideoDirectory string
	var lowQuality bool
	var saveTemps bool
	var useOldFade bool

	flag.BoolVar(&lowQuality, "l", false, "(boolean): Low Quality, include to generate a lower quality video (480p instead of 720p)")
	flag.BoolVar(&saveTemps, "s", false, "(boolean): Save Temporaries, include to save temporary files generated during video process)")
	flag.BoolVar(&useOldFade, "f", false, "(boolean): Fadetype, include to use the non-xfade default transitions for video")

	flag.StringVar(&slideshowDirectory, "t", "", "[filepath]: Template Name, specify a template to use (if not included searches current folder for template)")
	flag.StringVar(&outputDirectory, "o", "", "[filepath]: Output Location, specify where to store final result (default is current directory)")
	flag.StringVar(&temporaryDirectory, "td", "", "[filepath]: Temporary Directory, used to specify a location to store the temporary files used in video production (default is OS' temp folder/storybuilder-*)")
	flag.StringVar(&overlayVideoDirectory, "ov", "", "[filepath]: Overlay Video, specify test video location to create overlay video")
	flag.Parse()

	options := options{slideshowDirectory, outputDirectory, temporaryDirectory, overlayVideoDirectory, lowQuality, saveTemps, useOldFade}

	return options

}
