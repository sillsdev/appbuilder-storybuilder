package main

import (
	"os/exec"
)

func cmdIndividualVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "./"+ImageDirectory,
		"-t", duration+"ms", "-filter_complex", zoom_cmd,
		"-shortest", "-pix_fmt", "yuv420p", "-y", finalOutputDirectory)

	return cmd
}

func cmdVideoLength(inputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputDirectory,
	)
	return cmd
}

func cmdTrimLengthOfVideo(duration string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", "./temp/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		"./temp/final.mp4",
	)
	return cmd
}
