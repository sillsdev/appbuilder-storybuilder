package ffmpeg_pkg

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetVersion() *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-version")

	return cmd
}

func ScaleImage(imagePath string, height string, width string, imageOutputPath string) {
	cmd := exec.Command("ffmpeg", "-i", imagePath,
		"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
		"-y", imageOutputPath)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func TrimLengthOfVideo(duration string, tempPath string) {
	cmd := exec.Command("ffmpeg",
		"-i", tempPath+"/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		tempPath+"/final.mp4",
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func GetVideoLength(inputDirectory string) float64 {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputDirectory,
	)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	//store the video length in an array
	result, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 8)
	return result
}

func CreateTempVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "./"+ImageDirectory,
		"-t", duration+"ms", "-filter_complex", zoom_cmd,
		"-shortest", "-pix_fmt", "yuv420p", "-y", finalOutputDirectory)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

// Function to check CMD error output when running commands
func checkCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}
