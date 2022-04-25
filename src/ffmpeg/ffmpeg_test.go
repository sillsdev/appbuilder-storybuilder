package ffmpeg_pkg

import (
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
	CheckCMDError(output, err)

	ffmpeg = strings.TrimSpace(string(output))
}

func Test_CmdGetVersion(t *testing.T) {
	tests := []struct {
		name string
		want *exec.Cmd
	}{
		{
			"get version ffmpeg cmd",
			exec.Command("ffmpeg", "-version"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdGetVersion().String(); got != tt.want.String() {
				t.Errorf("getVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CmdScaleImage(t *testing.T) {
	type args struct {
		imagePath       string
		height          string
		width           string
		imageOutputPath string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			"Scaling image ffmpeg cmd",
			args{imagePath: "../TestInput/Jn01.1-18-title.jpg", height: "852", width: "480", imageOutputPath: "../TestInput/Jn01.1-18-title.jpg"},
			exec.Command("ffmpeg", "-i", "../TestInput/Jn01.1-18-title.jpg",
				"-vf", "scale=852:480,setsar=1:1",
				"-y", "../TestInput/Jn01.1-18-title.jpg"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdScaleImage(tt.args.imagePath, tt.args.height, tt.args.width, tt.args.imageOutputPath).String(); got != tt.want.String() {
				t.Errorf("cmdScaleImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CmdTrimLengthOfVideo(t *testing.T) {
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
			"trim length of video ffmpeg cmd",
			args{duration: "19200ms", tempPath: "temp"},
			exec.Command("ffmpeg",
				"-i", "temp"+"/merged_video.mp4",
				"-c", "copy", "-t", "19200ms",
				"-y",
				"temp"+"/final.mp4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdTrimLengthOfVideo(tt.args.duration, tt.args.tempPath).String(); got != tt.want.String() {
				t.Errorf("cmdTrimLengthOfVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CmdGetVideoLength(t *testing.T) {
	type args struct {
		inputDirectory string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			"video length ffmpeg cmd",
			args{inputDirectory: "../TestInput/sample_video.mp4"},
			exec.Command("ffprobe",
				"-v", "error",
				"-show_entries", "format=duration",
				"-of", "default=noprint_wrappers=1:nokey=1",
				"../TestInput/sample_video.mp4"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdGetVideoLength(tt.args.inputDirectory).String(); got != tt.want.String() {
				t.Errorf("getVideoLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CmdCreateTempVideo(t *testing.T) {
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
				finalOutputDirectory: "temp/temp1-8.mp4"},
			exec.Command(ffmpeg+" -loop", "1", "-i", "../TestInput/VB-John 1v1.jpg", "-t",
				"9400ms", "-filter_complex", "scale=8000:-1,zoompan=z='1/((0.718)-(0.001)*on)':x='0.282*iw+0.000*iw*on':y='0.088*ih+0.000*ih*on':d=235:fps=25,scale=1280:720,setsar=1:1", "-shortest", "-pix_fmt", "yuv420p", "-y", "temp/temp1-8.mp4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdCreateTempVideo(tt.args.ImageDirectory, tt.args.duration, tt.args.zoom_cmd, tt.args.finalOutputDirectory).String(); got != tt.want.String() {
				t.Errorf("cmdCreateTempVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateZoomCommand(t *testing.T) {
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
			if got := CreateZoomCommand(tt.args.Motions, tt.args.Duration); got != tt.want {
				t.Errorf("createZoomCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CmdCopyFile(t *testing.T) {
	type args struct {
		to   string
		from string
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		{
			"copy file ffmpeg cmd",
			args{to: "temp/final.mp4", from: "../final.mp4"},
			exec.Command("ffmpeg", "-i", "temp/final.mp4", "-y", "../final.mp4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdCopyFile(tt.args.to, tt.args.from).String(); got != tt.want.String() {
				t.Errorf("cmdCopyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
