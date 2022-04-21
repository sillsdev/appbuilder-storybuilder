package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	OS "github.com/sillsdev/appbuilder-storybuilder/os"
	"github.com/sillsdev/appbuilder-storybuilder/slideshow"
)

var slideshowDirFlag string
var outputDirFlag string
var tempDirFlag string
var overlayVideoDirFlag string

// Main function
func main() {
	// Ask the user for options
	lowQuality, helpFlag, saveTemps, useOldfade := parseFlags(&slideshowDirFlag, &outputDirFlag, &tempDirFlag, &overlayVideoDirFlag)
	if *helpFlag {
		displayHelpMessage()
		return
	}

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

func parseFlags(templateName *string, outputPath *string, tempPath *string, overlayVideoDirFlag *string) (*bool, *bool, *bool, *bool) {
	var lowQuality = flag.Bool("l", false, "Include to produce a lower quality video (1280x720 => 852x480)")
	var help = flag.Bool("h", false, "Include option flag to display list of possible flags and their uses")
	var saveTemps = flag.Bool("s", false, "Include to save the temporary files after production")
	var useOldFade = flag.Bool("f", false, "Include to use traditional ffmpeg fade")
	flag.StringVar(templateName, "t", "", "Specify template to use")
	flag.StringVar(outputPath, "o", "", "Specify output location")
	flag.StringVar(tempPath, "td", "", "Specify temp directory location (If user wishes to save temporary files created during production)")
	flag.StringVar(overlayVideoDirFlag, "ov", "", "Specify test video location to create overlay video")
	flag.Parse()

	return lowQuality, help, saveTemps, useOldFade
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

func displayHelpMessage() {
	println("Usage: program-name [OPTIONS]\n")
	println("Options list:\n")
	println("            -t [filepath]: Template Name, specify a template to use (if not included searches current folder for template)\n")
	println("            -s (boolean): Save Temporaries, include to save temporary files generated during video process)\n")
	println("            -td [filepath]: Temporary Directory, used to specify a location to store the temporary files used in video production (default is current-directory/temp)\n")
	println("            -o [filepath]: Output Location, specify where to store final result (default is current directory)\n")
	println("            -l (boolean): Low Quality, include to generate a lower quality video (480p instead of 720p)\n")
	println("            -v (boolean): Verbosity, include to increase the verbosity of the status messages printed during video process\n")
	println("            -h (boolean): Help, include to display this help message and quit\n")
}
