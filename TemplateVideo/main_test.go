package main

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
)

type Test struct {
	d   *fs.DirEntry
	err *error
}

func TestParse(t *testing.T) {
	inputFile := "test.slideshow"
	var expectedOutput string
	data := readData(inputFile)
	for i, slide := range data.Slide {
		if i == 0 {
			// Test background filename
			expectedOutput = "background.mp3"
			backFilename := slide.Audio.Background_Filename.Path
			if backFilename != expectedOutput {
				t.Error(fmt.Sprintf("expected background filename to be %s, but got %s", expectedOutput, backFilename))
			}
		} else {
			expectedOutput = "narration.mp3"
			audio := slide.Audio.Filename.Name
			if audio != expectedOutput {
				t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedOutput, audio))
			}
		}
		expectedOutput = fmt.Sprintf("test-%d.jpg", i)
		image := slide.Image.Name
		if image != expectedOutput {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedOutput, image))
		}
		if slide.Motion.Start != "" {
			expectedOutput = "0.0 0.1 0.2 0.3"
			start := slide.Motion.Start
			if start != expectedOutput {
				t.Error(fmt.Sprintf("expected motion start to be %s, but got %s", expectedOutput, start))
			}
			expectedOutput = "1 2 3 4"
			end := slide.Motion.End
			if end != expectedOutput {
				t.Error(fmt.Sprintf("expected motion end to be %s, but got %s", expectedOutput, end))
			}
		}
		if slide.Transition.Type != "" {
			expectedOutput = "transitionTest"
			transitionType := slide.Transition.Type
			if transitionType != expectedOutput {
				t.Error(fmt.Sprintf("expected transtion type to be %s, but got %s", expectedOutput, transitionType))
			}
			expectedOutput = "1000"
			transitionDuration := slide.Transition.Duration
			if transitionDuration != expectedOutput {
				t.Error(fmt.Sprintf("expected transtion duration to be %s, but got %s", expectedOutput, transitionDuration))
			}
		}
		if slide.Timing.Start != "" {
			expectedOutput = "1234"
			timingStart := slide.Timing.Start
			if timingStart != expectedOutput {
				t.Error(fmt.Sprintf("expected timing start to be %s, but got %s", expectedOutput, timingStart))
			}
			expectedOutput = "5678"
			timingDuration := slide.Timing.Duration
			if timingDuration != expectedOutput {
				t.Error(fmt.Sprintf("expected timing duration to be %s, but got %s", expectedOutput, timingDuration))
			}
		}

	}
}

//expected output should be a png
// expected output should be a png
// func TestScaleImage(t *testing.T) {
// 	inputFile := input_images(Images[i])
// 	input := height
// 	input2 := width
// 	expectedOutput := fmt.Sprintf("test_%d.jpg", i)
// 	if inputFile != expectedOutput {
// 		t.Errorf("expected image here")
// 	}
// }
/*
	Test function to check if we are getting right version by comparing
	if a version is less than 4.3.0//
	if the version is equal to 4.3.0 //
				or
	if the version is higher than 4.3.0
*/
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

/*
	Test function to check if we are getting the template provided
// */

func TestFindTemplate(t *testing.T) {
	name := ".slideshow"
	//want := regexp.MustCompile(`\b` + name + `\b`)
	err := errors.New("Test error")
	if findTemplate(name, nil, err) == err {
		t.Logf("Pass, Expected and empty string")
	} else {
		t.Errorf("Failed, Expected an empty string")
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

// Unit test Ideas

// check the audio and video are the same lenght
// check the component files are correct
// check transitions from one images to another
// check look at what ffmpeg is generating
