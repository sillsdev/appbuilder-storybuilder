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

	expectedImages := []string{"Jn01.1-18-title.jpg", "./VB-John 1v1.jpg", "./VB-John 1v3.jpg", "./VB-John 1v4.jpg", "./VB-John 1v5a.jpg",
		"./VB-John 1v5b.jpg", "./VB-John 1v6.jpg", "Gospel of John-credits.jpg"}
	for i := 0; i < len(expectedImages); i++ {
		if expectedImages[i] != Images[i] {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedImages[i], Images[i]))
		}
	}

	expectedAudios := []string{"../music-intro-Jn.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", ""}
	for i := 0; i < len(expectedAudios); i++ {
		if expectedAudios[i] != Audios[i] {
			t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedAudios[i], Audios[i]))
		}
	}

	expectedBackAudioPath := "../music-intro-Jn.mp3"
	if expectedBackAudioPath != BackAudioPath {
		t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedBackAudioPath, BackAudioPath))
	}

	expectedBackAudioVolume := ""
	if expectedBackAudioVolume != BackAudioVolume {
		t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedBackAudioVolume, BackAudioVolume))
	}

	expectedTransitions := []string{"fade", "fade", "crossfade", "fade", "fade", "wipeleft", "fade"}
	for i := 0; i < len(expectedTransitions); i++ {
		if expectedTransitions[i] != Transitions[i] {
			t.Error(fmt.Sprintf("expected transition to be %s, but got %s", expectedTransitions[i], Transitions[i]))
		}
	}

	expectedTransitionDurations := []string{"1000", "1000", "2000", "1000", "1000", "3000", "1000"}
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

	if Motions[0][0][0] == Motions[1][1][1] {

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

// func TestReadFile(t *testing.T) {
// 	data, err := ioutil.ReadFile("data.slideshow")
// 	if err != nil {

// 	}
// 	if string(readData) != nil {

// 	}
//// }//
