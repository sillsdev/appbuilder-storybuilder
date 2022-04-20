package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gordon-cs/SIL-Video/Compiler/ffmpeg_pkg"
	opSys "github.com/gordon-cs/SIL-Video/Compiler/os"
	"github.com/gordon-cs/SIL-Video/Compiler/slideshow"
)

var slideshowDirectory string
var outputLocation string
var tempLocation string
var overlayVideoPath string

var FFmpeg = ffmpeg_pkg.NewFfmpeg()
var OS = opSys.NewOS()

// Main function
func main() {
	// Ask the user for options
	lowQuality, helpFlag, saveTemps, useOldfade := parseFlags(&slideshowDirectory, &outputLocation, &tempLocation, &overlayVideoPath)
	if *helpFlag {
		displayHelpMessage()
		return
	}

	// Create a temporary folder to store temporary files created when created a video
	tempLocation = OS.CreateTemporaryFolder(tempLocation)

	// Create directory if output directory does not exist
	if outputLocation != "" {
		OS.CreateDirectory(outputLocation)
	}

	// Search for a template in local folder if no template is provided
	if slideshowDirectory == "" {
		fmt.Println("No template provided, searching local folder...")
		filepath.WalkDir(".", findTemplate)
	}

	start := time.Now()

	// Parse in the various pieces from the template

	slideshow := slideshow.NewSlideshow(slideshowDirectory)

	Images := slideshow.GetImages()
	Transitions := slideshow.GetTransitions()
	TransitionDurations := slideshow.GetTransitionDurations()
	Timings := slideshow.GetTimings()
	Audios := slideshow.GetAudios()
	Motions := slideshow.GetMotions()

	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")

	var fadeType string = FFmpeg.CheckFFmpegVersion()

	//Scaling images depending on video quality option
	fmt.Println("Scaling images...")
	if *lowQuality {
		FFmpeg.ScaleImages(Images, "852", "480")
	} else {
		FFmpeg.ScaleImages(Images, "1280", "720")
	}

	fmt.Println("Creating video...")

	if fadeType == "X" && !*useOldfade {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(Images, Transitions, TransitionDurations, Timings, Audios, Motions, tempLocation)
		FFmpeg.MergeTempVideos(Images, Transitions, TransitionDurations, Timings, tempLocation)
		FFmpeg.AddAudio(Timings, Audios, tempLocation)
		FFmpeg.CopyFinal(tempLocation, outputLocation)
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(Images, Transitions, TransitionDurations, Timings, Audios, Motions, tempLocation)
		FFmpeg.MergeTempVideosOldFade(Images, TransitionDurations, Timings, tempLocation)
		FFmpeg.AddAudio(Timings, Audios, tempLocation)
		FFmpeg.CopyFinal(tempLocation, outputLocation)
	}

	fmt.Println("Finished making video...")

	// If user did not specify the -s flag at runtime, delete all the temporary videos
	if !*saveTemps {
		OS.DeleteTemporaryVideos(saveTemps)
	}

	fmt.Println("Video production completed!")
	duration := time.Since(start)
	fmt.Sprintln(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))

	if overlayVideoPath != "" {
		fmt.Println("Creating overlay video...")
		FFmpeg.CreateOverlaidVideoForTesting(overlayVideoPath, outputLocation)
		fmt.Println("Finished creating overlay video")
	}
}

func parseFlags(templateName *string, outputPath *string, tempPath *string, overlayVideoPath *string) (*bool, *bool, *bool, *bool) {
	var lowQuality = flag.Bool("l", false, "Include to produce a lower quality video (1280x720 => 852x480)")
	var help = flag.Bool("h", false, "Include option flag to display list of possible flags and their uses")
	var saveTemps = flag.Bool("s", false, "Include to save the temporary files after production")
	var useOldFade = flag.Bool("f", false, "Include to use traditional ffmpeg fade")
	flag.StringVar(templateName, "t", "", "Specify template to use")
	flag.StringVar(outputPath, "o", "", "Specify output location")
	flag.StringVar(tempPath, "td", "", "Specify temp directory location (If user wishes to save temporary files created during production)")
	flag.StringVar(overlayVideoPath, "ov", "", "Specify test video location to create overlay video")
	flag.Parse()

	return lowQuality, help, saveTemps, useOldFade
}

func removeFileNameFromDirectory(slideshowDirectory string) string {
	var template_directory_split []string
	if runtime.GOOS == "windows" { // Windows uses '\' for filepaths
		template_directory_split = strings.Split(slideshowDirectory, "\\")
	} else {
		template_directory_split = strings.Split(slideshowDirectory, "/")
	}
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

// Function to find the .slideshow template if none provided
func findTemplate(s string, d fs.DirEntry, err error) error {
	slideRegEx := regexp.MustCompile(`.+(.slideshow)$`) // Regular expression to find the .slideshow file
	if err != nil {
		return err
	}
	if slideRegEx.MatchString(d.Name()) {
		if slideshowDirectory == "" {
			fmt.Println("Found template: " + s + "\nUsing found template...")
			slideshowDirectory = s
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
