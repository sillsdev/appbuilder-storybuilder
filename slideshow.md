# Slideshow Format for StoryBuilder

Author: [Chris Hubbard](mailto:chris_hubbard@sil.org), SIL International

## Introduction

Starting with [version 5.0](https://software.sil.org/scriptureappbuilder/release-notes/#v5.0), [Scripture App Builder](software.sil.org/scriptureappbuilder) (SAB) has had a feature that can use slideshow templates with the audio from a SAB project to generate [PhotoStage](https://www.nchsoftware.com/slideshow/) projects that can produce a slideshow videos.

This has worked for many years, but it is time intensive and error prone and there is a desire to scale the usage of this process. This will require a new slideshow format that includes all the metadata that used to reside in the PhotoStage template project.

The new slideshow format will have two separate stages:
* Template
* Language Specific Instantiation 

## Slideshow Template

The .slideshow format used in the slideshow template has been modified to include the information that was present in the PhotoStage project included in the old template format.  It contains `<timing>` elements for the non-narative slides and `<narration>` elements to specify the Chapter/Verse for the narrative slides. [Here](SampleInput/SAB%20Video%20Production%20v4/templates) is an example of the new template format.

## Slideshow Instantiation

The .slideshow format for a specific language includes the audio needed for the narration along with the specific timing for each narrative slide. [Here](SampleInput/SAB%20Video%20Production%20v4/%5Beng%5D%20World%20English%20Bible) is an example of a slideshow that was generated from a project using the World English Bible.

## Elements Needed For StoryBuilder

The language specific instantiation of the .slideshow file will have extra elements that are not needed for the StoryBuilder process.  These are the ones that are relevant:

```xml
<?xml version="1.0" encoding="utf-8"?>
<slideshow>
  <slide>
    <audio background="play">
      <background-filename volume="50">../music-intro-SA.mp3</background-filename>
    </audio>
    <image>MAT2-title.jpg</image>
    <timing duration="5000"/>
    <transition type="fade" duration="1000"/>
  </slide>
  <slide>
    <audio>
      <filename>narration-001.mp3</filename>
    </audio>
    <image>Mat-02-v01.jpg</image>
    <motion start="0.395 0 0.605 0.468" end="0 0 1 0.774"/>
    <timing duration="10040"/>
    <transition type="wipeleft" duration="1000"/>
  </slide>
  <slide>
    <audio>
      <filename>narration-001.mp3</filename>
    </audio>
    <image>Mat-02-v02.jpg</image>
    <motion start="0 0 1 0.774" end="0.378 0 0.622 0.484"/>
    <timing duration="8080"/>
    <transition type="radial" duration="1000"/>
  </slide>
  <slide>
    <audio>
      <filename>narration-001.mp3</filename>
    </audio>
    <image>Mat-02-v03.jpg</image>
    <motion start="0.163 0.008 0.663 0.516" end="0 0 1 0.774"/>
    <timing duration="5800"/>
  </slide>
  <slide>
    <image>Sweet-credits.jpg</image>
    <timing duration="5000"/>
  </slide>
</slideshow>
```

## &lt;slide>

Container element for each slide of the slideshow.  It must contain `<audio>`, `<image>`, and `<timing>`.

## &lt;audio>

Specifies the filenames for the audio associated with the slide. There are two possible audio tracks for a slide: background and narration. 

The background audio is specified with the `<background-filename>` element and can have a volume attribute which is the percent volume (1-100).  The background attribute on the audio element specifies how background audio should be processed for the slide.  A value of "play" means to use the `<background-filename>` included in this `<audio>` element.  A value of "continue" means to keep playing the audio on this slide.

The narration audio is specified with the `<filename>` element which can be replicated over several `<slide>` elements with different `<timing>` elements for the segment of the narration audio that will be used for that slide.

## &lt;image> with JPEG file

Specifies the image filename (.jpg) for the slide.  There will be additional `<image>` elements that have a lang attribute with a LibreOffice document (.odg) for localization of the slide.  These additional `<image>` elements should be ignored and just use the `<image>` element with an image filename (.jpg).

## &lt;motion>

Specifies the animation to be applied to the image of the slide.  The start and end attributes specify the rectangles for the Ken Burns effect for the slide.  The values in the start and end attributes are string with these properties of the rectangle: left, top, width, height.  The values of these properties are the percentage of the associated width and height of the image.

## &lt;timing>

Specifies the timing of the audio within the slide.  The duration attributes specify the milliseconds within the audio that should be played for the slide. Multiple slides will use the same audio filename and the audio should be played continuously until it is not referenced by a slide.

## &lt;transition>

Specifies the transition that should happen between slides.  The duration attribute specifies the milliseconds of the transition and should be split between the two slides.  The type attribute specifies the name of the transition to be used.  If there is no `<transition>` element, then assume a 1000 millisecond transition using the fade transition. See the [FFmpeg xfade documentation](https://ffmpeg.org/ffmpeg-filters.html#xfade) for the list of transitions.

## Elements Used By Slideshow Generator (StoryBuilder Ignore)

You can see extra data that is used by the slideshow generation process that can be safely ignored by the video generation process.

## &lt;narration>

The narration tag is used by the slideshow generation process to specify the range of verses that should be extracted from the scripture audio.  This element should be ignored by the video generation process.

## &lt;title>

The title tag is used by the slideshow generation process  to generate the title slide for the slideshow. This element should be ignored by the video generation process.

## &lt;image> with ODG files

There are image tags that have .odg files specified.  These are used by the slideshow generation process to localize slides.  The image tags with .odg files should be ignored by the video generation process.