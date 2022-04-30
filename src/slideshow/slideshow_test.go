package slideshow

import (
	"fmt"
	"testing"
)

func TestReadSlideshow(t *testing.T) {
	templateName := "../../TestInput/test.slideshow"

	slideshow := NewSlideshow(templateName, false)

	expectedImages := []string{"../../TestInput/Jn01.1-18-title.jpg", "../../TestInput/./VB-John 1v1.jpg", "../../TestInput/./VB-John 1v3.jpg", "../../TestInput/./VB-John 1v4.jpg", "../../TestInput/./VB-John 1v5a.jpg",
		"../../TestInput/./VB-John 1v5b.jpg", "../../TestInput/./VB-John 1v6.jpg", "../../TestInput/Gospel of John-credits.jpg"}
	for i := 0; i < len(expectedImages); i++ {
		if expectedImages[i] != slideshow.images[i] {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedImages[i], slideshow.images[i]))
		}
	}

	expectedAudios := []string{"../../TestInput/./music-intro-Jn.mp3", "../../TestInput/narration-j-001.mp3", "../../TestInput/narration-j-001.mp3", "../../TestInput/narration-j-001.mp3", "../../TestInput/narration-j-001.mp3", "../../TestInput/narration-j-001.mp3", "../../TestInput/narration-j-001.mp3", ""}
	for i := 0; i < len(expectedAudios); i++ {
		if expectedAudios[i] != slideshow.audios[i] {
			t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedAudios[i], slideshow.audios[i]))
		}
	}

	expectedTransitions := []string{"fade", "fade", "circleopen", "fade", "fade", "wipeleft", "wipeleft"}
	for i := 0; i < len(expectedTransitions); i++ {
		if expectedTransitions[i] != slideshow.transitions[i] {
			t.Error(fmt.Sprintf("expected transition to be %s, but got %s", expectedTransitions[i], slideshow.transitions[i]))
		}
	}

	expectedTransitionDurations := []string{"1000", "1000", "2000", "1000", "1000", "3000", "3000"}
	for i := 0; i < len(expectedTransitionDurations); i++ {
		if expectedTransitionDurations[i] != slideshow.transitionDurations[i] {
			t.Error(fmt.Sprintf("expected transition duration to be %s, but got %s", expectedTransitionDurations[i], slideshow.transitionDurations[i]))
		}
	}

	expectedTimings := []string{"5000", "9400", "5960", "4200", "2280", "2280", "10880", "5000"}
	for i := 0; i < len(expectedTimings); i++ {
		if expectedTimings[i] != slideshow.timings[i] {
			t.Error(fmt.Sprintf("expected timing duration to be %s, but got %s", expectedTimings[i], slideshow.timings[i]))
		}
	}

	expectedMotions := [][][]float64{{{0, 0, 1, 1}, {0, 0, 1, 1}}, {{0.282, 0.088, 0.718, 0.717}, {0.391, 0.115, 0.475, 0.478}}, {{0.297, 0.204, 0.554, 0.558}, {0.515, 0.381, 0.416, 0.416}},
		{{0.114, 0.071, 0.663, 0.664}, {0.129, 0.159, 0.46, 0.46}}, {{0, 0, 1, 1}, {0, 0, 1, 1}}, {{0.109, 0.097, 0.629, 0.628}, {0.144, 0.142, 0.47, 0.469}},
		{{0.124, 0.071, 0.455, 0.451}, {0.144, 0.053, 0.782, 0.779}}, {{0, 0, 1, 1}, {0, 0, 1, 1}}}

	for i := 0; i < len(expectedMotions); i++ {
		if expectedMotions[i][0][0] != slideshow.motions[i][0][0] {
			t.Error(fmt.Sprintf("expected motion[%d][0][0] to be %f, but got %f", i, expectedMotions[i][0][0], slideshow.motions[i][0][0]))
		}
		if expectedMotions[i][0][1] != slideshow.motions[i][0][1] {
			t.Error(fmt.Sprintf("expected motion[%d][0][1] to be %f, but got %f", i, expectedMotions[i][0][1], slideshow.motions[i][0][1]))
		}
		if expectedMotions[i][0][2] != slideshow.motions[i][0][2] {
			t.Error(fmt.Sprintf("expected motion[%d][0][2] to be %f, but got %f", i, expectedMotions[i][0][2], slideshow.motions[i][0][2]))
		}
		if expectedMotions[i][0][3] != slideshow.motions[i][0][3] {
			t.Error(fmt.Sprintf("expected motion[%d][0][3] to be %f, but got %f", i, expectedMotions[i][0][3], slideshow.motions[i][0][3]))
		}
		if expectedMotions[i][1][0] != slideshow.motions[i][1][0] {
			t.Error(fmt.Sprintf("expected motion[%d][1][0] to be %f, but got %f", i, expectedMotions[i][1][0], slideshow.motions[i][1][0]))
		}
		if expectedMotions[i][1][1] != slideshow.motions[i][1][1] {
			t.Error(fmt.Sprintf("expected motion[%d][1][1] to be %f, but got %f", i, expectedMotions[i][1][1], slideshow.motions[i][1][1]))
		}
		if expectedMotions[i][1][2] != slideshow.motions[i][1][2] {
			t.Error(fmt.Sprintf("expected motion[%d][1][2] to be %f, but got %f", i, expectedMotions[i][1][2], slideshow.motions[i][1][2]))
		}
		if expectedMotions[i][1][3] != slideshow.motions[i][1][3] {
			t.Error(fmt.Sprintf("expected motion[%d][1][3] to be %f, but got %f", i, expectedMotions[i][1][3], slideshow.motions[i][1][3]))
		}
	}
}
