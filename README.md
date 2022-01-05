# SIL Video Converter
21-22 Gordon College CS Senior Project<br>
Se Hee Hyung, David Gurge, Roddy Ngolomingi, Hyungyu Park<br>

Senior Project Problem Statement [Link](https://docs.google.com/document/d/1Xcbwg4K3Fhv3oUFh-9i_Q81I1Y1p6ym8wsgSIHjBBA0/edit?usp=sharing).<br>
Initial Design Document [Link](https://docs.google.com/document/d/16FA-5HbT2uVkvgAXTeTjRo2QJxEuIR1Bfjdc5Mci7FI/edit?usp=sharing).<br>
End-To-End Proposal [Link](https://docs.google.com/document/d/1h8e6FNbOrI4lRuMVRTbiZil3-PrC2OoKQ6b0vckxl1w/edit?usp=sharing).<br>
1st Lo-Fi Usability Test [Link](https://drive.google.com/file/d/1L9HBFWGztYsH0RSPItrjFPIrZDt0xkz8/view?usp=sharing).<br>
1st Lo-Fi Usability Test Report [Link](https://docs.google.com/document/d/1-MmKXZmo_WDw9Ju-L8kHIel8QrqPs31j3IiaVdt6B-k/edit?usp=sharing)


MVP Document: [Link](https://docs.google.com/document/d/1ZZWAUzAl-bXXmUvLlqPjvj4Cw5By6yFNDDiA70PlY2E/edit?usp=sharing)<br>
Proof of Work Repo (Python): https://github.com/sillsdev/storybuilder/tree/v2<br>
Inputs Template Documentation: [Link](https://docs.google.com/document/d/1J4X6RWUgXnI0aeaLEv4ePBXeZJQJSMgZ-WLQNx7Hcj8/edit?usp=sharing)

# Plans/Ideas for Project:

# How-To Documentation
1. Download FFmpeg https://www.ffmpeg.org by selecting the appropriate .zip for your OS (Here's a basic tutorial for [Windows](https://www.wikihow.com/Install-FFmpeg-on-Windows), [Mac](https://manual.audacityteam.org/man/installing_ffmpeg_for_mac.html), and [Linux](https://www.tecmint.com/install-ffmpeg-in-linux/)) 
2. Download and install GO https://golang.org/dl/ (Should include instructions on their page)
3. Edit the base paths for repository and FFmpeg in main.go (these are on line 10)
4. Put any images (.png, .jpg, etc) and audios (.mp3, .wav, etc) into a folder labeled "input" and include a data.slideshow in the main directory with your main.go file (the .slideshow file documentation is linked above)
5. Also create a folder labeled "output" for the finished videos to placed in
6. Run code in CLI set to main directory of repo with "go run main.go read.go" or just "go run ."

# issue 
1. |Fixed| After running the the command  "go run main.go" or " go run ." it will create a video in your output folder, 
before running the command again we need to delete the created videos in the output folder every time we want to create a new video. 
