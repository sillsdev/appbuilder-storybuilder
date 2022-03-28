package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	templateName := "../TestInput/test.slideshow"

	Images, Audios, BackAudioPath, BackAudioVolume, Transitions, TransitionDurations, Timings, Motions := parseSlideshow(templateName)

	expectedImages := []string{"../TestInput/Jn01.1-18-title.jpg", "../TestInput/./VB-John 1v1.jpg", "../TestInput/./VB-John 1v3.jpg", "../TestInput/./VB-John 1v4.jpg", "../TestInput/./VB-John 1v5a.jpg",
		"../TestInput/./VB-John 1v5b.jpg", "../TestInput/./VB-John 1v6.jpg", "../TestInput/Gospel of John-credits.jpg"}
	for i := 0; i < len(expectedImages); i++ {
		if expectedImages[i] != Images[i] {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedImages[i], Images[i]))
		}
	}

	expectedAudios := []string{"../TestInput/./music-intro-Jn.mp3", "../TestInput/narration-j-001.mp3", "../TestInput/narration-j-001.mp3", "../TestInput/narration-j-001.mp3", "../TestInput/narration-j-001.mp3", "../TestInput/narration-j-001.mp3", "../TestInput/narration-j-001.mp3", ""}
	for i := 0; i < len(expectedAudios); i++ {
		if expectedAudios[i] != Audios[i] {
			t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedAudios[i], Audios[i]))
		}
	}

	expectedBackAudioPath := "./music-intro-Jn.mp3"
	if expectedBackAudioPath != BackAudioPath {
		t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedBackAudioPath, BackAudioPath))
	}

	expectedBackAudioVolume := ""
	if expectedBackAudioVolume != BackAudioVolume {
		t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedBackAudioVolume, BackAudioVolume))
	}

	expectedTransitions := []string{"fade", "fade", "circleopen", "fade", "fade", "fade", "wipeleft", "fade"}
	for i := 0; i < len(expectedTransitions); i++ {
		if expectedTransitions[i] != Transitions[i] {
			t.Error(fmt.Sprintf("expected transition to be %s, but got %s", expectedTransitions[i], Transitions[i]))
		}
	}

	expectedTransitionDurations := []string{"1000", "1000", "2000", "1000", "1000", "1000", "3000", "1000"}
	for i := 0; i < len(expectedTransitionDurations); i++ {
		if expectedTransitionDurations[i] != TransitionDurations[i] {
			t.Error(fmt.Sprintf("expected transition duration to be %s, but got %s", expectedTransitionDurations[i], TransitionDurations[i]))
		}
	}

	expectedTimings := []string{"5000", "9400", "5960", "4200", "2280", "2280", "10880", "5000"}
	for i := 0; i < len(expectedTimings); i++ {
		if expectedTimings[i] != Timings[i] {
			t.Error(fmt.Sprintf("expected timing duration to be %s, but got %s", expectedTimings[i], Timings[i]))
		}
	}

	expectedMotions := [][][]float64{{{0, 0, 1, 1}, {0, 0, 1, 1}}, {{0.282, 0.088, 0.718, 0.717}, {0.391, 0.115, 0.475, 0.478}}, {{0.297, 0.204, 0.554, 0.558}, {0.515, 0.381, 0.416, 0.416}},
		{{0.114, 0.071, 0.663, 0.664}, {0.129, 0.159, 0.46, 0.46}}, {{0, 0, 1, 1}, {0, 0, 1, 1}}, {{0.109, 0.097, 0.629, 0.628}, {0.144, 0.142, 0.47, 0.469}},
		{{0.124, 0.071, 0.455, 0.451}, {0.144, 0.053, 0.782, 0.779}}, {{0, 0, 1, 1}, {0, 0, 1, 1}}}

	for i := 0; i < len(expectedMotions); i++ {
		if expectedMotions[i][0][0] != Motions[i][0][0] {
			t.Error(fmt.Sprintf("expected motion[%d][0][0] to be %f, but got %f", i, expectedMotions[i][0][0], Motions[i][0][0]))
		}
		if expectedMotions[i][0][1] != Motions[i][0][1] {
			t.Error(fmt.Sprintf("expected motion[%d][0][1] to be %f, but got %f", i, expectedMotions[i][0][1], Motions[i][0][1]))
		}
		if expectedMotions[i][0][2] != Motions[i][0][2] {
			t.Error(fmt.Sprintf("expected motion[%d][0][2] to be %f, but got %f", i, expectedMotions[i][0][2], Motions[i][0][2]))
		}
		if expectedMotions[i][0][3] != Motions[i][0][3] {
			t.Error(fmt.Sprintf("expected motion[%d][0][3] to be %f, but got %f", i, expectedMotions[i][0][3], Motions[i][0][3]))
		}
		if expectedMotions[i][1][0] != Motions[i][1][0] {
			t.Error(fmt.Sprintf("expected motion[%d][1][0] to be %f, but got %f", i, expectedMotions[i][1][0], Motions[i][1][0]))
		}
		if expectedMotions[i][1][1] != Motions[i][1][1] {
			t.Error(fmt.Sprintf("expected motion[%d][1][1] to be %f, but got %f", i, expectedMotions[i][1][1], Motions[i][1][1]))
		}
		if expectedMotions[i][1][2] != Motions[i][1][2] {
			t.Error(fmt.Sprintf("expected motion[%d][1][2] to be %f, but got %f", i, expectedMotions[i][1][2], Motions[i][1][2]))
		}
		if expectedMotions[i][1][3] != Motions[i][1][3] {
			t.Error(fmt.Sprintf("expected motion[%d][1][3] to be %f, but got %f", i, expectedMotions[i][1][3], Motions[i][1][3]))
		}
	}
}

func TestScaleImages(t *testing.T) {
	imageName := "Jn01.1-18-title.jpg"
	image_path := "../TestInput/" + imageName

	Images := []string{}
	Images = append(Images, image_path)

	scaleImages(Images, "852", "480")

	cmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v:0", "-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0", image_path)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
	output_string := strings.TrimSpace(string(output))

	expectedOutput := "852x480"

	if output_string != expectedOutput {
		t.Error(fmt.Sprintf("expected image %s to have widthxheight = %s, but got %s", imageName, expectedOutput, output_string))
	}
}
