package main

import (
	"os/exec"
)

func cmdCreateTempVideo(ImageDirectory string, duration string, zoom_cmd string, finalOutputDirectory string) *exec.Cmd {
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", "./"+ImageDirectory,
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

func cmdTrimLengthOfVideo(duration string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", "./temp/merged_video.mp4",
		"-c", "copy", "-t", duration,
		"-y",
		"./temp/final.mp4",
	)
	return cmd
}

func cmdAddBackgroundMusic(backgroundAudioPath string, volume string) *exec.Cmd {
	cmd := exec.Command("ffmpeg",
		"-i", "./temp/mergedVideo.mp4",
		"-i", backgroundAudioPath,
		"-filter_complex", "[1:0]volume="+volume+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		"-y", "../finalvideo.mp4",
	)
	return cmd
}
