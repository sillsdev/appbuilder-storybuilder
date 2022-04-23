package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sillsdev/appbuilder-storybuilder/src/options"
	OS "github.com/sillsdev/appbuilder-storybuilder/src/os"
	"github.com/sillsdev/appbuilder-storybuilder/src/slideshow"
)

// Main function
func main() {
	// Ask the user for options
	optionFlags := options.ParseFlags()

	// Create a temporary folder to store temporary files
	tempDirectory := OS.CreateDirectory(optionFlags.TemporaryDirectory)

	// Create directory if output directory does not exist
	if optionFlags.OutputDirectory != "" {
		OS.CreateDirectory(optionFlags.OutputDirectory)
	}

	// Search for a template in local folder if no template is provided
	if optionFlags.SlideshowDirectory == "" {
		fmt.Println("No template provided, searching local folder...")
		filepath.WalkDir(".", findTemplate(optionFlags))
	}

	start := time.Now()

	// Parse in the various pieces from the template
	slideshow := slideshow.NewSlideshow(optionFlags.SlideshowDirectory)

	fmt.Println("Scaling images...")
	slideshow.ScaleImages(optionFlags.LowQuality)

	fmt.Println("Creating video...")
	slideshow.CreateVideo(optionFlags.UseOldFade, tempDirectory, optionFlags.OutputDirectory)

	// If user did not specify the -s flag at runtime, delete all the temporary videos
	if !(optionFlags.SaveTemps) {
		OS.DeleteTemporaryDirectory(optionFlags.SaveTemps)
	}

	fmt.Println("Video production completed!")
	duration := time.Since(start)
	fmt.Sprintln(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))

	if optionFlags.OverlayVideoDirectory != "" {
		fmt.Println("Creating overlay video...")
		slideshow.CreateOverlaidVideo(optionFlags.OverlayVideoDirectory, optionFlags.OutputDirectory)
		fmt.Println("Finished creating overlay video")
	}
}

// Function to find the .slideshow template if none provided
func findTemplate(optionFlags options.Options) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		slideRegEx := regexp.MustCompile(`.+(.slideshow)$`) // Regular expression to find the .slideshow file
		if err != nil {
			return err
		}
		if slideRegEx.MatchString(d.Name()) {
			if optionFlags.SlideshowDirectory == "" {
				fmt.Println("Found template: " + path + "\nUsing found template...")
				optionFlags.SetSlideshowDirectory(path)
			}
		}
		return nil
	}
}
