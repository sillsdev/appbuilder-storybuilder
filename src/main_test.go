package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

var ffmpeg string

func init() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "ffmpeg")
	} else {
		cmd = exec.Command("which", "ffmpeg")
	}
	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)

	ffmpeg = strings.TrimSpace(string(output))
}

func TestParse(t *testing.T) {
	templateName := "../TestInput/test.slideshow"

	finalVideoName, Images, Audios, Transitions, TransitionDurations, Timings, Motions := parseSlideshow(templateName)

	finalVideoName = strings.TrimSuffix(finalVideoName, ".slideshow")

	expectedName := string("test")
	if expectedName != finalVideoName {
		t.Error(fmt.Sprintf("expected final video name to be %s, but got %s", expectedName, finalVideoName))
	}

	expectedImages := []string{"Jn01.1-18-title.jpg", "./VB-John 1v1.jpg", "./VB-John 1v3.jpg", "./VB-John 1v4.jpg", "./VB-John 1v5a.jpg",
		"./VB-John 1v5b.jpg", "./VB-John 1v6.jpg", "Gospel of John-credits.jpg"}
	for i := 0; i < len(expectedImages); i++ {
		if expectedImages[i] != Images[i] {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedImages[i], Images[i]))
		}
	}

	expectedAudios := []string{"./music-intro-Jn.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", "narration-j-001.mp3", ""}
	for i := 0; i < len(expectedAudios); i++ {
		if expectedAudios[i] != Audios[i] {
			t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedAudios[i], Audios[i]))
		}
	}

	expectedTransitions := []string{"fade", "fade", "circleopen", "fade", "fade", "fade", "wipeleft", "fade"}
	for i := 0; i > len(expectedTransitions); i++ {
		if expectedTransitions[i] != Transitions[i] {
			t.Error(fmt.Sprintf("expected transition to be %s, but got %s", expectedTransitions[i], Transitions[i]))
		}
	}

	expectedTransitionDurations := []string{"1000", "1000", "2000", "1000", "1000", "1000", "3000", "1000"}
	for i := 0; i > len(expectedTransitionDurations); i++ {
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

func Test_scaleImages(t *testing.T) {
	type args struct {
		Images []string
		height string
		width  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Scaling images to smaller size",
			args{Images: []string{"../TestInput/Jn01.1-18-title.jpg"}, height: "852", width: "480"},
		},
		{
			"Scaling images to original size",
			args{Images: []string{"../TestInput/Jn01.1-18-title.jpg"}, height: "1280", width: "720"},
			//args{Images: []string{"../TestInput/Jn01.1-18-title.jpg", "../TestInput/./VB-John 1v1.jpg", "../TestInput/./VB-John 1v3.jpg", "../TestInput/./VB-John 1v4.jpg", "../TestInput/./VB-John 1v5a.jpg",
			// "../TestInput/./VB-John 1v5b.jpg", "../TestInput/./VB-John 1v6.jpg", "../TestInput/Gospel of John-credits.jpg"}, height: "1280", width: "720"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaleImages(tt.args.Images, tt.args.height, tt.args.width)
		})
	}

	cmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v:0", "-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0", "../TestInput/Jn01.1-18-title.jpg")

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
	output_string := strings.TrimSpace(string(output))

	expectedOutput := "1280x720"

	//t.Log(output_string)

	if output_string != expectedOutput {
		t.Error(fmt.Sprintf("expected image %s to have widthxheight = %s, but got %s", "Jn01.1-18-title.jpg", expectedOutput, output_string))
	}
}

func TestCheckFFmpegVersion(t *testing.T) {
	got := checkFFmpegVersion("4.2.9")
	want := "F"
	if want != got {
		t.Errorf("Failed, expected " + want + " got " + got + " for 4.2.9")
	} else {
		t.Logf("Pass, expected " + want + " got " + got + " for 4.2.9")
	}

	got = checkFFmpegVersion("5.0")
	want = "X"
	if want != got {
		t.Errorf("Failed, expected " + want + " got " + got + " for 5.0")
	} else {
		t.Logf("Pass, expected " + want + " got " + got + " for 5.0")
	}

	got = checkFFmpegVersion("4.3.0")
	want = "X"
	if want != got {
		t.Errorf("Failed, expected " + want + " got " + got + " for 4.3.0")
	} else {
		t.Logf("Pass, expected " + want + " got " + got + " for 4.3.0")
	}
}

func TestFindTemplate(t *testing.T) {
	name := ".slideshow"
	err := errors.New("Test error")
	if findTemplate(name, nil, err) != nil {
		t.Logf("Pass, Expected nothing")
	} else {
		t.Errorf("Failed, Expected nothing")
	}

}

func TestErroRFindTemplate(t *testing.T) {
	err := errors.New("Test error")
	if findTemplate(".slideshow", nil, err) == err {
		t.Logf("Pass, Expected an Error")
	} else {
		t.Errorf("Failed, Expedcted an Error")
	}
}

func Test_cmdCreateTempVideo(t *testing.T) {
	type args struct {
		ImageDirectory       string
		duration             string
		zoom_cmd             string
		finalOutputDirectory string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			"Creating temp video ffmpeg command for VB-John 1v1.jpg",
			args{ImageDirectory: "../TestInput/VB-John 1v1.jpg",
				duration:             "9400",
				zoom_cmd:             "scale=8000:-1,zoompan=z='1/((0.718)-(0.001)*on)':x='0.282*iw+0.000*iw*on':y='0.088*ih+0.000*ih*on':d=235:fps=25,scale=1280:720,setsar=1:1",
				finalOutputDirectory: "./temp/temp1-8.mp4"},
			exec.Command(ffmpeg+" -loop", "1", "-i", "./../TestInput/VB-John 1v1.jpg", "-t",
				"9400ms", "-filter_complex", "scale=8000:-1,zoompan=z='1/((0.718)-(0.001)*on)':x='0.282*iw+0.000*iw*on':y='0.088*ih+0.000*ih*on':d=235:fps=25,scale=1280:720,setsar=1:1", "-shortest", "-pix_fmt", "yuv420p", "-y", "./temp/temp1-8.mp4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdCreateTempVideo(tt.args.ImageDirectory, tt.args.duration, tt.args.zoom_cmd, tt.args.finalOutputDirectory).String(); got != tt.want.String() {
				t.Errorf("cmdCreateTempVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createZoomCommand(t *testing.T) {
	type args struct {
		Motions  [][]float64
		Duration []float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Creating zoom command for VB-John 1v1.jpg",
			args{Motions: [][]float64{{0.282, 0.088, 0.718, 0.717}, {0.391, 0.115, 0.475, 0.478}},
				Duration: []float64{9400}},
			"scale=8000:-1,zoompan=z='1/((0.718)-(0.001)*on)':x='0.282*iw+0.000*iw*on':y='0.088*ih+0.000*ih*on':d=235:fps=25,scale=1280:720,setsar=1:1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createZoomCommand(tt.args.Motions, tt.args.Duration); got != tt.want {
				t.Errorf("createZoomCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdGetVideoLength(t *testing.T) {
	type args struct {
		inputDirectory string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			"get correct video duration",
			args{inputDirectory: "../TestInput/sample_video.mp4"},
			// check the command that we are running is the right command.
			exec.Command("ffprobe",
				"-v", "error",
				"-show_entries", "format=duration",
				"-of", "default=noprint_wrappers=1:nokey=1",
				"../TestInput/sample_video.mp4"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdGetVideoLength(tt.args.inputDirectory).String(); got != tt.want.String() {
				t.Errorf("createZoomCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdTrimLengthOfVideo(t *testing.T) {
	type args struct {
		duration string
		tempPath string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			" ffmpeg command for triming video",
			args{duration: "30ms",
				tempPath: "./temp"},
			exec.Command("ffmpeg",
				"-i", "./temp"+"/merged_video.mp4",
				"-c", "copy", "-t", "30ms",
				"-y",
				"./temp"+"/final.mp4"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdTrimLengthOfVideo(tt.args.duration, tt.args.tempPath).String(); got != tt.want.String() {
				t.Errorf("cmdTrimLengthOfVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdCopyFile(t *testing.T) {
	type args struct {
		oldPath string
		newPath string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			" checking the video oldpath and NewpPath ",
			args{oldPath: "./temp/final.pm4",
				newPath: "./output/final.mp4"},

			exec.Command("ffmpeg", "-i", "./temp/final.pm4", "-y", "./output/final.mp4"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdCopyFile(tt.args.oldPath, tt.args.newPath).String(); got != tt.want.String() {
				t.Errorf("cmdAddBackgroundMusic() = %v, want %v", got, tt.want)
			}
		})
	}
}
