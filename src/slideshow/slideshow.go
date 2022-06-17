package slideshow

import (
	"fmt"
	Image "image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"

	FFmpeg "github.com/sillsdev/appbuilder-storybuilder/src/ffmpeg"
	"github.com/sillsdev/appbuilder-storybuilder/src/helper"
)

/* Structure of a .slideshow
 * 	images: filepath strings to the images to be used for each slide
 *	audios: filepath strings to the narration audios to be used for each slide
 *	transitions: strings describing which (Xfade only) transition to be used in between each slide
 *	transitionDurations: strings describing the time (in milliseconds) for each transition to last
 *	timings: strings describing the time (in milliseconds) for each slide to last, also used for motions
 *	motions: arrays of floats describing the dimensions and positions for the start and end rectangles for zoom/pan effects
 *	templateName: string parsed from the .slideshow filename to be used for the final video product
 */
type slideshow struct {
	images              []string
	audios              []string
	transitions         []string
	transitionDurations []string
	timings             []string
	motions             [][][]float64
	templateName        string
	tempPath            string
}

/* Function to create a new slideshow from a .slideshow template. The code parses the pieces out
 * and stores them in the slideshow struct
 *
 * Parameters:
 *			slideshowDirectory - the filepath to the .slideshow to be parsed
 *			v - verbose flag to determine what feedback to print
 * Returns:
 *			slideshow - the filled slideshow structure, containing all the data parsed
 */
func NewSlideshow(slideshowDirectory string, v bool, tempPath string) slideshow {
	slideshow_template := readSlideshowXML(slideshowDirectory)

	Images := []string{}
	Audios := []string{}
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := []string{}
	Motions := [][][]float64{}

	fmt.Println("Parsing .slideshow file...")

	templateDir, template_name := splitFileNameFromDirectory(slideshowDirectory)

	for i, slide := range slideshow_template.Slide {
		Timings = append(Timings, slide.Timing.Duration)
		if slide.Audio.Background_Filename.Path != "" { // Intro music is stored differently in the xml
			if slide.Audio.Filename.Name != "" {
				FFmpeg.MergeAudios(templateDir+slide.Audio.Background_Filename.Path, templateDir+slide.Audio.Filename.Name, Timings[i-1], Timings[i], tempPath)
				Audios = append(Audios, path.Join(tempPath, "mergedAudio.mp3"))
			} else {
				Audios = append(Audios, templateDir+slide.Audio.Background_Filename.Path)
			}
		} else {
			if slide.Audio.Filename.Name != "" {
				Audios = append(Audios, templateDir+slide.Audio.Filename.Name)
			} else {
				Audios = append(Audios, "")
			}
		}
		Images = append(Images, templateDir+slide.Image.Name)
		if slide.Transition.Type == "" { // Default to a basic crossfade if no transition provided
			Transitions = append(Transitions, "fade")
		} else {
			Transitions = append(Transitions, slide.Transition.Type)
		}
		if slide.Transition.Duration == "" { // Default to 1000ms transition if none provided
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		}
		var motions = [][]float64{}
		if slide.Motion.Start == "" { // If no motion specified, default to a static "zoom/pan" effect
			motions = [][]float64{{0, 0, 1, 1}, {0, 0, 1, 1}}
		} else {
			motions = [][]float64{helper.ConvertStringToFloat(slide.Motion.Start), helper.ConvertStringToFloat(slide.Motion.End)}
		}
		Motions = append(Motions, motions)
	}

	if v {
		fmt.Printf("Parsed %d images, %d audios, %d transitions, %d transition durations, %d timings, and %d motions, from %s\n",
			len(Images), len(Audios), len(Transitions), len(TransitionDurations), len(Timings), len(Motions), template_name)
	}
	slideshow := slideshow{Images, Audios, Transitions, TransitionDurations, Timings, Motions, template_name, tempPath}

	fmt.Println("Parsing completed...")

	return slideshow
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ternaryInt(condition bool, x int, y int) int {
	if condition {
		return x
	}
	return y
}

func readImage(name string) (Image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, err := jpeg.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func percentToPixel(percent float64, whole int) int {
	return int(math.Floor(percent*float64(whole) + 0.5))
}

func pixelToPercent(pixel int, whole int) float64 {
	return math.Floor(float64(pixel)/float64(whole)*100.0+0.5) / 100.0
}

func enlargeBoundsToAspectRatio(input Image.Rectangle, targetRatio float64) Image.Rectangle {
	curRatio := float64(input.Dx()) / float64(input.Dy())
	width := input.Dx()
	height := input.Dy()
	if math.Abs(curRatio-targetRatio) < 0.01 {
		// Current rect is fine
		return input
	} else if curRatio < targetRatio {
		// Change width
		width = int(math.Round(float64(height) * targetRatio))
	} else {
		// Change height
		height = int(math.Round(float64(width) / targetRatio))
	}
	return Image.Rect(input.Min.X, input.Min.Y, width, height)
}

func (s slideshow) CropImage(i int, v bool) (string, error) {
	img, err := readImage(s.images[i])
	if err != nil {
		return "", err
	}

	imgBounds := img.Bounds()
	heightImg := imgBounds.Dy()
	heightHd := (imgBounds.Dx() * 9) / 16
	if v {
		fmt.Printf("Crop: [%d] width=%d, heightImg=%d, heightHd=%d\n", i, imgBounds.Dx(), heightImg, heightHd)
	}
	if Abs(heightImg-heightHd) < 5 {
		// It is close enough to 16/9 aspect ratio so do nothing
		if v {
			fmt.Printf("Crop: [%d] close enough. using: %s\n\n", i, s.images[i])
		}
		return s.images[i], nil
	} else {
		// Find 16x9 bounding box
		// 1. union the bounds
		// 2. find 16x9 that encloses it (might be larger in one dimension than the image)
		// 3. create a new image the size of the enlarged bounds (without the x, y offsets)
		// 4. copy the contents of the union bounds to the new images
		// 5. adjust the motion for the new image
		if v {
			fmt.Printf("Crop: [%d] startMotion=[%f %f %f %f] endMotion=[%f %f %f %f]\n", i,
				s.motions[i][0][0], s.motions[i][0][1], s.motions[i][0][2], s.motions[i][0][3],
				s.motions[i][1][0], s.motions[i][1][1], s.motions[i][1][2], s.motions[i][1][3])
		}

		startBounds := Image.Rect(
			percentToPixel(s.motions[i][0][0], imgBounds.Dx()),
			percentToPixel(s.motions[i][0][1], imgBounds.Dy()),
			percentToPixel(s.motions[i][0][0]+s.motions[i][0][2], imgBounds.Dx()),
			percentToPixel(s.motions[i][0][1]+s.motions[i][0][3], imgBounds.Dy()))
		endBounds := Image.Rect(
			percentToPixel(s.motions[i][1][0], imgBounds.Dx()),
			percentToPixel(s.motions[i][1][1], imgBounds.Dy()),
			percentToPixel(s.motions[i][1][0]+s.motions[i][1][2], imgBounds.Dx()),
			percentToPixel(s.motions[i][1][1]+s.motions[i][1][3], imgBounds.Dy()))
		if v {
			fmt.Printf("Crop: [%d] startPixels=[%d %d %d %d] endPixels=[%d %d %d %d]\n", i,
				startBounds.Min.X, startBounds.Min.Y, startBounds.Dx(), startBounds.Dy(),
				endBounds.Min.X, endBounds.Min.Y, endBounds.Dx(), endBounds.Dy())
		}

		unionBounds := startBounds.Union(endBounds)
		if v {
			fmt.Printf("Crop: [%d] union=[%d %d %d %d]\n", i, unionBounds.Min.X, unionBounds.Min.Y, unionBounds.Dx(), unionBounds.Dy())
		}
		unionBoundSize := Image.Rect(0, 0, unionBounds.Dx(), unionBounds.Dy())
		enlargedBounds := enlargeBoundsToAspectRatio(unionBounds, 16.0/9.0)
		newImageBounds := Image.Rect(0, 0, enlargedBounds.Dx(), enlargedBounds.Dy())
		if v {
			fmt.Printf("Crop: [%d] enlarge=[%d %d %d %d]\n", i, enlargedBounds.Min.X, enlargedBounds.Min.Y, enlargedBounds.Dx(), enlargedBounds.Dy())
		}

		newImg := Image.NewRGBA(newImageBounds)
		draw.Draw(newImg, newImageBounds, &Image.Uniform{color.RGBA{255, 255, 255, 255}}, Image.ZP, draw.Src)
		draw.Draw(newImg, unionBoundSize, img, unionBounds.Min, draw.Src)

		startResult := Image.Rect(startBounds.Min.X-enlargedBounds.Min.X, startBounds.Min.Y-enlargedBounds.Min.Y, startBounds.Dx(), startBounds.Dy())
		endResult := Image.Rect(endBounds.Min.X-enlargedBounds.Min.X, endBounds.Min.Y-enlargedBounds.Min.Y, endBounds.Dx(), endBounds.Dy())
		if v {
			fmt.Printf("Crop: [%d] startResPixels=[%d %d %d %d] endResPixels=[%d %d %d %d]\n", i,
				startResult.Min.X, startResult.Min.Y, startResult.Dx(), startResult.Dy(),
				endResult.Min.X, endResult.Min.Y, endResult.Dx(), endResult.Dy())
		}

		s.motions[i][0][0] = pixelToPercent(startBounds.Min.X-enlargedBounds.Min.X, newImageBounds.Dx())
		s.motions[i][0][1] = pixelToPercent(startBounds.Min.Y-enlargedBounds.Min.Y, newImageBounds.Dy())
		s.motions[i][0][2] = pixelToPercent(startBounds.Dx(), newImageBounds.Dx())
		s.motions[i][0][3] = pixelToPercent(startBounds.Dy(), newImageBounds.Dy())
		s.motions[i][1][0] = pixelToPercent(endBounds.Min.X-enlargedBounds.Min.X, newImageBounds.Dx())
		s.motions[i][1][1] = pixelToPercent(endBounds.Min.Y-enlargedBounds.Min.Y, newImageBounds.Dy())
		s.motions[i][1][2] = pixelToPercent(endBounds.Dx(), newImageBounds.Dx())
		s.motions[i][1][3] = pixelToPercent(endBounds.Dy(), newImageBounds.Dy())
		if v {
			fmt.Printf("Crop: [%d] startMotion=[%f %f %f %f] endMotion=[%f %f %f %f]\n", i,
				s.motions[i][0][0], s.motions[i][0][1], s.motions[i][0][2], s.motions[i][0][3],
				s.motions[i][1][0], s.motions[i][1][1], s.motions[i][1][2], s.motions[i][1][3])
		}

		outputImage := path.Join(s.tempPath, path.Base(s.images[i]))
		if v {
			fmt.Printf("Crop: [%d] outputing: %s\n\n", i, outputImage)
		}
		fd, err := os.Create(outputImage)
		if err != nil {
			return "", err
		}
		jpeg.Encode(fd, newImg, nil)
		return outputImage, nil
	}
}

/* Function to scale all the input images depending on video quality
 * option to a uniform height/width to prevent issues in the video creation process.
 *
 * Parameters:
 *			lowQuality - specifies whether to generate a lower quality video by scaling the images to a smaller dimension
 *			v - verbose flag to determine what feedback to print
 */
func (s slideshow) ScaleImages(lowQuality bool, v bool) {
	width := "1280"
	height := "720"

	if lowQuality {
		println("-l specified, producing lower quality video")
		width = "852"
		height = "480"
	}

	totalNumImages := len(s.images)
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		go func(i int) {
			defer wg.Done()
			inputImage, err := s.CropImage(i, v)
			if err != nil {
				log.Fatal(err)
			}
			outputImage := path.Join(s.tempPath, path.Base(s.images[i]))
			cmd := FFmpeg.CmdScaleImage(inputImage, height, width, outputImage)
			s.images[i] = outputImage
			output, err := cmd.CombinedOutput()
			FFmpeg.CheckCMDError(output, err)
		}(i)
	}

	wg.Wait()
}

/* Function to create a video with all the data parsed from the .slideshow
 *
 * Parameters:
 *			useOldFade - specifies whether to use the old fade style instead of XFade, if desired
 *			tempDirectory - filepath to the temp folder to store the temporary videos created
 *			outputDirectory - filepath to the location to store the final completed video
 *			v - verbose flag to determine what feedback to print
 */
func (s slideshow) CreateVideo(useOldfade bool, tempDirectory string, outputDirectory string, v bool) {
	if v {
		fmt.Println("Temp Directory: " + tempDirectory)
		fmt.Println("Output Directory: " + outputDirectory)
	}
	// Checking FFmpeg version to use Xfade
	fmt.Println("Checking FFmpeg version...")
	var fadeType string = FFmpeg.ParseVersion()
	useXfade := fadeType == "X" && !useOldfade

	final_template_name := strings.TrimSuffix(s.templateName, ".slideshow")

	if useXfade {
		fmt.Println("FFmpeg version is bigger than 4.3.0, using Xfade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.timings, s.audios, s.motions, tempDirectory, v)
		FFmpeg.MergeTempVideos(s.images, s.transitions, s.transitionDurations, s.timings, tempDirectory, v)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory, v)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	} else {
		fmt.Println("FFmpeg version is smaller than 4.3.0, using old fade transition method...")
		FFmpeg.MakeTempVideosWithoutAudio(s.images, s.timings, s.audios, s.motions, tempDirectory, v)
		FFmpeg.MergeTempVideosOldFade(s.images, s.transitionDurations, s.timings, tempDirectory, v)
		FFmpeg.AddAudio(s.timings, s.audios, tempDirectory, v)
		FFmpeg.CopyFinal(tempDirectory, outputDirectory, final_template_name)
	}

	fmt.Println("Finished making video...")
}

// Helper function to generate an overlaid video of the software's result and a comparison video
func (s slideshow) CreateOverlaidVideo(finalVideoDirectory string, testVideoDirectory string, overlaidVideoDirectory string) {
	FFmpeg.CreateOverlaidVideoForTesting(finalVideoDirectory, testVideoDirectory, overlaidVideoDirectory)
}

/* Function to separate the .slideshow filename from the directory path
 *
 * Parameters:
 *		slideshowDirectory - path to the .slideshow file
 * Returns:
 *		template_directory - folder path leading up to the .slideshow file
 *		template_name - name of the .slideshow file
 */
func splitFileNameFromDirectory(slideshowDirectory string) (string, string) {
	var template_directory_split []string

	template_directory_split = regexp.MustCompile("[\\/\\\\]+").Split(slideshowDirectory, -1)

	template_directory := ""
	template_name := template_directory_split[len(template_directory_split)-1]

	if len(template_directory_split) == 1 {
		if runtime.GOOS != "windows" {
			template_directory = ""
		}
	} else {
		for i := 0; i < len(template_directory_split)-1; i++ {
			template_directory += template_directory_split[i] + "/"
		}
	}

	return template_directory, template_name
}
