package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	OS "github.com/sillsdev/appbuilder-storybuilder/src/os"
	"github.com/sillsdev/appbuilder-storybuilder/src/slideshow"
)

var slideshowDirFlag string
var outputDirFlag string
var tempDirFlag string
var overlayVideoDirFlag string

// Main function
func main() {
	// Ask the user for options
	lowQuality, saveTemps, useOldfade := parseFlags(&slideshowDirFlag, &outputDirFlag, &tempDirFlag, &overlayVideoDirFlag)

	// Create a temporary folder to store temporary files
	tempDirectory := OS.CreateDirectory(tempDirFlag)

	// Create directory if output directory does not exist
	if outputDirFlag != "" {
		OS.CreateDirectory(outputDirFlag)
	}

	// Search for a template in local folder if no template is provided
	if slideshowDirFlag == "" {
		fmt.Println("No template provided, searching local folder...")
		filepath.WalkDir(".", findTemplate)
	}

	start := time.Now()

	// Parse in the various pieces from the template
	slideshow := slideshow.NewSlideshow(slideshowDirFlag)

	fmt.Println("Scaling images...")
	slideshow.ScaleImages(lowQuality)

	fmt.Println("Creating video...")
	slideshow.CreateVideo(useOldfade, tempDirectory, outputDirFlag)

	// If user did not specify the -s flag at runtime, delete all the temporary videos
	if !*saveTemps {
		OS.DeleteTemporaryDirectory(saveTemps)
	}

	fmt.Println("Video production completed!")
	duration := time.Since(start)
	fmt.Sprintln(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))

	if overlayVideoDirFlag != "" {
		fmt.Println("Creating overlay video...")
		slideshow.CreateOverlaidVideo(overlayVideoDirFlag, outputDirFlag)
		fmt.Println("Finished creating overlay video")
	}
}

func parseFlags(templateName *string, outputPath *string, tempPath *string, overlayVideoPath *string) (*bool, *bool, *bool) {
	var lowQuality = flag.Bool("l", false, "(boolean): Low Quality, include to generate a lower quality video (480p instead of 720p)")
	var saveTemps = flag.Bool("s", false, "(boolean): Save Temporaries, include to save temporary files generated during video process)")
	var useOldFade = flag.Bool("f", false, "(boolean): Fadetype, include to use the non-xfade default transitions for video")
	flag.StringVar(templateName, "t", "", "[filepath]: Template Name, specify a template to use (if not included searches current folder for template)")
	flag.StringVar(outputPath, "o", "", "[filepath]: Output Location, specify where to store final result (default is current directory)")
	flag.StringVar(tempPath, "td", "", "[filepath]: Temporary Directory, used to specify a location to store the temporary files used in video production (default is OS' temp folder/storybuilder-*)")
	flag.StringVar(overlayVideoPath, "ov", "", "[filepath]: Overlay Video, specify test video location to create overlay video")
	flag.Parse()

	return lowQuality, saveTemps, useOldFade
}

// Function to find the .slideshow template if none provided
func findTemplate(s string, d fs.DirEntry, err error) error {
	slideRegEx := regexp.MustCompile(`.+(.slideshow)$`) // Regular expression to find the .slideshow file
	if err != nil {
		return err
	}
	if slideRegEx.MatchString(d.Name()) {
		if slideshowDirFlag == "" {
			fmt.Println("Found template: " + s + "\nUsing found template...")
			slideshowDirFlag = s
		}
	}
	return nil
}
