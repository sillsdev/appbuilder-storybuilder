package main

import (
	"fmt"
	"os/exec"
)

func cmdCreateTempVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", ImageDirectory,
		"-t", duration+"ms", "-filter_complex", zoom_cmd,
		"-shortest", "-pix_fmt", "yuv420p", "-y", finalOutputDirectory)

	return cmd
}

func cmdGetVideoLength(inputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputDirectory,
	)
	return cmd
}

func cmdTrimLengthOfVideo(duration string, tempPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", tempPath+"/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		tempPath+"/final.mp4",
	)
	return cmd
}

func cmdCopyFile(oldPath string, newPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", oldPath, "-y", newPath)

	return cmd
}

func cmdScaleImage(imagePath string, height string, width string, imageOutputPath string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-i", imagePath,
		"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
		"-y", imageOutputPath)

	return cmd
}
