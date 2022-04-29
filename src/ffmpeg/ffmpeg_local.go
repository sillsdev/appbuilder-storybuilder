package ffmpeg_pkg

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// ffmpeg command to get the version number
func CmdGetVersion() *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-version")

	return cmd
}

// ffmpeg command to scale an image by specified height and width
func CmdScaleImage(imagePath string, height string, width string, imageOutputPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", imagePath,
		"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
		"-y", imageOutputPath)

	return cmd
}

// ffmpeg command to trim the video to a specific duration
func CmdTrimLengthOfVideo(duration string, tempPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", tempPath+"/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		tempPath+"/final.mp4",
	)

	return cmd
}

// ffmpeg command to get the length (in seconds) of a video
func CmdGetVideoLength(inputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputDirectory,
	)

	return cmd
}

// ffmpeg command to generate a single audioless video with the provided image and zoom/pan effects
func CmdCreateTempVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", ImageDirectory,
		"-t", duration+"ms", "-filter_complex", zoom_cmd,
		"-shortest", "-pix_fmt", "yuv420p", "-y", finalOutputDirectory)
	return cmd
}

/* Function to generate a proper ffmpeg filter to apply the zoom/pan effects
 *
 * Parameters:
 *		Motions - array of motion data containing start and end rectangles
 *		TimingDuration - duration of the zoom/pan effect
 * Returns:
 *		final_cmd - the finalized zoom/pan command for a single video
 */
func CreateZoomCommand(Motions [][]float64, TimingDuration float64) string {
	num_frames := int(TimingDuration / (1000.0 / 25.0))

	size_init := Motions[0][3]
	size_change := Motions[1][3] - size_init
	size_incr := size_change / float64(num_frames)

	var x_init float64 = Motions[0][0]
	var x_end float64 = Motions[1][0]
	var x_change float64 = x_end - x_init
	var x_incr float64 = x_change / float64(num_frames)

	var y_init float64 = Motions[0][1]
	var y_end float64 = Motions[1][1]
	var y_change float64 = y_end - y_init
	var y_incr float64 = y_change / float64(num_frames)

	zoom_cmd := fmt.Sprintf("1/((%.3f)%s(%.3f)*on)", size_init-size_incr, checkSign(size_incr), math.Abs(size_incr))
	x_cmd := fmt.Sprintf("%0.3f*iw%s%0.3f*iw*on", x_init-x_incr, checkSign(x_incr), math.Abs(x_incr))
	y_cmd := fmt.Sprintf("%0.3f*ih%s%0.3f*ih*on", y_init-y_incr, checkSign(y_incr), math.Abs(y_incr))
	final_cmd := fmt.Sprintf("scale=8000:-1,zoompan=z='%s':x='%s':y='%s':d=%d:fps=25,scale=1280:720,setsar=1:1", zoom_cmd, x_cmd, y_cmd, num_frames)

	return final_cmd
}

// Function to check CMD error output when running commands
func CheckCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

// ffmpeg command to copy a video from one location to another
func CmdCopyFile(to string, from string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", to, "-y", from)
	return cmd
}

// ffmpeg command to check the sign of a number
func checkSign(num float64) string {
	result := math.Signbit(num)

	if result {
		return "-"
	} else {
		return "+"
	}
}

/* Function to trim the end of the video and remove excess empty audio when the audio file is longer than the video file
 */
func trimEnd(tempPath string) {
	fmt.Println("Trimming end of merged video...")

	cmd := CmdGetVideoLength(tempPath + "/video_with_no_audio.mp4")
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)

	video_length, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

	//match the video length of the merged video with the true length of the video
	cmd = CmdTrimLengthOfVideo(fmt.Sprintf("%f", video_length), tempPath)
	output, err = cmd.CombinedOutput()
	CheckCMDError(output, err)
}
