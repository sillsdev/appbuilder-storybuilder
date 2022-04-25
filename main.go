package main

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sillsdev/appbuilder-storybuilder/src/helper"
	"github.com/sillsdev/appbuilder-storybuilder/src/options"
	OS "github.com/sillsdev/appbuilder-storybuilder/src/os"
	"github.com/sillsdev/appbuilder-storybuilder/src/slideshow"
)

var filePath string

// Main function
func main() {
	// Ask the user for options
	optionFlags := options.ParseFlags()

	// Create a temporary folder to store temporary files
	tempDirectory, err := OS.CreateDirectory(optionFlags.TemporaryDirectory)
	helper.Check(err)

	// Create directory if output directory does not exist
	if optionFlags.OutputDirectory != "" {
		_, err := OS.CreateDirectory(optionFlags.OutputDirectory)
		helper.Check(err)
	}

	// Search for a template in local folder if no template is provided
	if optionFlags.SlideshowDirectory == "" {
		fmt.Println("No template provided, searching local folder...")

		err := filepath.WalkDir(".", findTemplate(optionFlags.SlideshowDirectory))

		if err.Error() == "FOUND TEMPLATE" {
			optionFlags.SetSlideshowDirectory(filePath)
		}
	}

	start := time.Now()

	// Parse in the various pieces from the template
	slideshow := slideshow.NewSlideshow(optionFlags.SlideshowDirectory)

	fmt.Println("Scaling images...")
	slideshow.ScaleImages(optionFlags.LowQuality)

	fmt.Println("Creating video...")
	slideshow.CreateVideo(optionFlags.UseOldFade, tempDirectory, optionFlags.OutputDirectory)

	fmt.Println("Video production completed!")
	duration := time.Since(start)
	fmt.Sprintln(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))

	if optionFlags.OverlayVideoDirectory != "" {
		fmt.Println("Creating overlay video...")

		finalVideoDirectory := tempDirectory + "/final.mp4"

		slideshow.CreateOverlaidVideo(finalVideoDirectory, optionFlags.OverlayVideoDirectory, optionFlags.OutputDirectory)
		fmt.Println("Finished creating overlay video")
	}

	// If user did not specify the -s flag at runtime, delete all the temporary videos
	if !(optionFlags.SaveTemps) {
		err := OS.DeleteTemporaryDirectory(tempDirectory)
		helper.Check(err)
	}

}

func findTemplate(slideshowDirectory string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, e error) error {
		slideRegEx := regexp.MustCompile(`.+(.slideshow)$`) // Regular expression to find the .slideshow file
		if e != nil {
			return e
		}
		if slideRegEx.MatchString(d.Name()) {
			if slideshowDirectory == "" {
				fmt.Println("Found template: " + path + "\nUsing found template...")

				filePath = path
				return errors.New("FOUND TEMPLATE")
			}
		}

		return nil
	}
}