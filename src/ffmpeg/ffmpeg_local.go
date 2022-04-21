package ffmpeg_pkg

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func getVersion() *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-version")

	return cmd
}

func scaleImage(imagePath string, height string, width string, imageOutputPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", imagePath,
		"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
		"-y", imageOutputPath)

	return cmd
}

func trimLengthOfVideo(duration string, tempPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", tempPath+"/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		tempPath+"/final.mp4",
	)

	return cmd
}

func getVideoLength(inputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputDirectory,
	)

	return cmd
}

func createTempVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "./"+ImageDirectory,
		"-t", duration+"ms", "-filter_complex", zoom_cmd,
		"-shortest", "-pix_fmt", "yuv420p", "-y", finalOutputDirectory)

	return cmd
}

// Function to check CMD error output when running commands
func checkCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

/* Function to trim the end of the video and remove excess empty audio when the audio file is longer than the video file
 */
func trimEnd(tempPath string) {
	fmt.Println("Trimming video...")

	cmd := getVideoLength(tempPath + "/video_with_no_audio.mp4")
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	video_length, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

	//match the video length of the merged video with the true length of the video
	cmd = trimLengthOfVideo(fmt.Sprintf("%f", video_length), tempPath)
	output, err = cmd.CombinedOutput()
	checkCMDError(output, err)
}

func copyFile(to string, from string) {
	cmd := exec.Command("ffmpeg", "-i", to, "-y", from)
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func checkSign(num float64) string {
	result := math.Signbit(num)

	if result {
		return "-"
	} else {
		return "+"
	}
}
