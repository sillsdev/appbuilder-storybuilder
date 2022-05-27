# SIL Video Converter

21-22 Gordon College CS Senior Project<br>
Se Hee Hyung, David Gurge, Roddy Ngolomingi, Hyungyu Park<br>

Senior Project Problem Statement [Link](https://docs.google.com/document/d/1Xcbwg4K3Fhv3oUFh-9i_Q81I1Y1p6ym8wsgSIHjBBA0/edit?usp=sharing).<br>
Initial Design Document [Link](https://docs.google.com/document/d/16FA-5HbT2uVkvgAXTeTjRo2QJxEuIR1Bfjdc5Mci7FI/edit?usp=sharing).<br>
End-To-End Proposal [Link](https://docs.google.com/document/d/1h8e6FNbOrI4lRuMVRTbiZil3-PrC2OoKQ6b0vckxl1w/edit?usp=sharing).<br>
1st Lo-Fi Usability Test [Link](https://drive.google.com/file/d/1L9HBFWGztYsH0RSPItrjFPIrZDt0xkz8/view?usp=sharing).<br>
1st Lo-Fi Usability Test Report [Link](https://docs.google.com/document/d/1-MmKXZmo_WDw9Ju-L8kHIel8QrqPs31j3IiaVdt6B-k/edit?usp=sharing)

MVP Document: [Link](https://docs.google.com/document/d/1ZZWAUzAl-bXXmUvLlqPjvj4Cw5By6yFNDDiA70PlY2E/edit?usp=sharing)<br>
Proof of Work Repo (Python): [Link](https://github.com/sillsdev/storybuilder/tree/v2)<br>
Inputs Template Documentation: [Link](slideshow.md)<br>
Design Document: [Link](https://docs.google.com/document/d/1vjogjaWZ0ww7rJtKz3J4iuVbbFrZF3KASdHBW-zPYfE/edit#)<br>
Final Presentation Slides: [Link](https://docs.google.com/presentation/d/1OxTRJqiOaVFwTOPpruyL800moecj-uWdPVHe_ZXdusA/)

# How-To Documentation

1. Download FFmpeg https://www.ffmpeg.org by selecting the appropriate .zip for your OS. Make sure the version number is greater than 4.3.0 to make full use of our code (Here's a basic tutorial for [Windows](https://www.wikihow.com/Install-FFmpeg-on-Windows), [Mac using Homebrew](https://sites.duke.edu/ddmc/2013/12/30/install-ffmpeg-on-a-mac/), and [Linux using a PPA with ffmpeg v4.4.1](https://launchpad.net/~savoury1/+archive/ubuntu/ffmpeg4))
   When installing with Homebrew (`brew install ffmpeg –ANY-OPTIONS-YOU-WANT`), ignore special options. Run `brew install ffmpeg` instead.
2. Download and extract executable for your system from repo's releases
3. Put any images (.png, .jpg, etc) and audios (.mp3, .wav, etc) into a folder, and also include a .slideshow xml file with proper parameters for the video ([.slideshow documentation linked here](https://github.com/gordon-cs/appbuilder-storybuilder/blob/main/slideshow.md))
4. Run code in a CLI set to directory of executable with "./executable_name" or just "executable_name" for Windows
5. There are also several flags you can include at runtime to alter the output or inputs:

   -h : Help, display list of possible flags and their uses

   -t : Template, used to input a specific template file to use, otherwise the program searches executable's current directory for any .slideshow files and uses the first it finds

   -o : Output location, used to specify where to store the finished video, will use executable's current directory by default

   -l : Lower quality, used to generate a lower quality video for smaller file size for easier distribution (default videos will be 1280x720)

   -td : Temporary Directory, used to specify a location to store the temporary files used in video production (default is in your OS' temp directory/storybuilder-\*)

   -v : Verbosity, used to modify how much output is reported on the commandline for debugging purposes (less verbose by default)

   -s : Save files, used to specify if user wants to preserve the temporary files used in the video production (videos are deleted by default)

   -f : Fadetype, include to use the non-xfade default transitions for video

   -ov : Overlay video, used to specify the location of a test video to create an overlay video with the generated video

# Testing Documentation

Our source code contains unit tests per packages, to which we are adding more tests as we progress. This is run as follows:

1. Ensure GoLang is installed properly, from their website [link](https://golang.org/dl/)
2. Navigate to folder with source code with a CLI and run "go test ./..." to execute all the unit tests for all packages and ensure all tests pass.
3. Run "go test test_filename.go" to run specific test files.

# Release Documentation

GitHub Actions are configured to build a release on tagged commits.

In order to generate a release version locally, follow the steps below:

1. Install GoReleaser [link](https://goreleaser.com/install/)
2. In a CLI, navigate to the root directory (contains main.go)
3. Run `goreleaser release --snapshot --rm-dist`

Binaries will be located in the `dist` folder.

If any of these steps cause issues you can reference the [GoReleaser documentation](https://goreleaser.com/)
