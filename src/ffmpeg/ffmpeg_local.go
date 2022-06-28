package ffmpeg_pkg

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"path"
)

/* Function to get the ffmpeg version
 *
 * Returns:
		executable "ffmpeg -version" cmd
*/
func CmdGetVersion() *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-version")

	return cmd
}

/* Function to scale an image to specified height and width
 *
 * Parameters:
 *		imagePath - directory of the jpg image location
 *		height - pixel height
 *		width - pixel width
 *		imageOutputPath - directory to save the scaled image
 * Returns:
		exectauble command
*/
func CmdScaleImage(imagePath string, height string, width string, imageOutputPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", imagePath,
		"-vf", fmt.Sprintf("scale=%s:%s", width, height)+",setsar=1:1",
		"-y", imageOutputPath)

	return cmd
}

/* Function to trim the video to a specified duration
 *
 * Parameters:
 *		duration - the length of the video in seconds
 *		tempPath - temporary directory path where all the temp files are saved
 * Returns:
		exectauble command
*/
func CmdTrimLengthOfVideo(duration string, tempPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", path.Join(tempPath, "merged_video.mp4"),
		"-c", "copy", "-t", duration,
		"-y",
		path.Join(tempPath, "final.mp4"),
	)

	return cmd
}

/* Function to get the length (seconds) of a video
 *
 * Parameters:
 *		inputPath - the path of the video to find the length of
 * Returns:
		exectauble command
*/
func CmdGetVideoLength(inputPath string) *exec.Cmd {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)
	// cmd := exec.Command("ffmpeg",
	// 	"-hide_banner",
	// 	"-i", inputPath,
	// 	"-f", "null", "-",
	// )

	return cmd
}

/* Function to generate a single audioless video with the provided image and zoom/pan effects
 *
 * Parameters:
 *		imageDirectory - directory of the image location
 *		duration - duration of the generated video (milliseconds)
 *		zoom_cmd - zoompan filter command
 *		finalOutputDirectory - directory to save the output video
 * Returns:
		exectauble command
*/
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

	zoom_cmd := fmt.Sprintf("1/((%.10f)%s(%.10f)*on)", size_init-size_incr, checkSign(size_incr), math.Abs(size_incr))
	x_cmd := fmt.Sprintf("%0.10f*iw%s%0.10f*iw*on", x_init-x_incr, checkSign(x_incr), math.Abs(x_incr))
	y_cmd := fmt.Sprintf("%0.10f*ih%s%0.10f*ih*on", y_init-y_incr, checkSign(y_incr), math.Abs(y_incr))
	final_cmd := fmt.Sprintf("scale=8000:-1,zoompan=z='%s':x='%s':y='%s':d=%d:fps=25,scale=1280:720,setsar=1:1", zoom_cmd, x_cmd, y_cmd, num_frames)

	return final_cmd
}

/* Function to check CMD error output when running commands
 *
 * Parameters:
 *		output - cmd output
 *		err - error of the cmd result
 */
func CheckCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

/* Function to copy a video from one location to another
 *
 * Parameters:
 *		to - directory of the video
 *		from - directory to move the video
 * Returns:
 *		exectauble ffmpeg cmd
 */
func CmdCopyFile(to string, from string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", to, "-y", from)
	return cmd
}

/* Function to check the sign of a number
 *
 * Parameters:
 *		num: float number to check the sign of
 * Returns:
 *		- or + depending on the sign of the number
 */
func checkSign(num float64) string {
	result := math.Signbit(num)

	if result {
		return "-"
	} else {
		return "+"
	}
}

/* Function to trim the end of the video and remove excess empty audio when the audio file is longer than the video file
 *
 * Parameters:
 *		tempPath - directory of where all the temporary files are saved
 */
func trimEnd(tempPath string) {
	fmt.Println("Trimming end of merged video...")

	video_length := GetVideoLength(tempPath + "/video_with_no_audio.mp4")

	//match the video length of the merged video with the true length of the video
	cmd := CmdTrimLengthOfVideo(fmt.Sprintf("%f", video_length), tempPath)
	output, err := cmd.CombinedOutput()
	CheckCMDError(output, err)
}
