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
Inputs Template Documentation: [Link](https://docs.google.com/document/d/1J4X6RWUgXnI0aeaLEv4ePBXeZJQJSMgZ-WLQNx7Hcj8/edit?usp=sharing)<br>
Design Document: [Link](https://docs.google.com/document/d/1vjogjaWZ0ww7rJtKz3J4iuVbbFrZF3KASdHBW-zPYfE/edit#)

# Plans/Ideas for Project:

# How-To Documentation
1. Download FFmpeg https://www.ffmpeg.org by selecting the appropriate .zip for your OS (Here's a basic tutorial for [Windows](https://www.wikihow.com/Install-FFmpeg-on-Windows), [Mac](https://manual.audacityteam.org/man/installing_ffmpeg_for_mac.html), and [Linux](https://www.tecmint.com/install-ffmpeg-in-linux/)) 
2. Download and install GO https://golang.org/dl/ (Should include instructions on their page)
3. Edit the base paths for repository and FFmpeg in main.go (these are on line 10)
4. Put any images (.png, .jpg, etc) and audios (.mp3, .wav, etc) into a the same folder as the executable, and also include a .slideshow xml file with parameters for the video (.slideshow documentation is listed above)
5. Run code in a CLI set to main directory of repo with "go run main.go read.go" or just "go run ."
6. There are also several flags you can include at runtime to alter the output or inputs:

    -v : Verbosity, used to modify how much output is reported on the commandline for debugging purposes (less verbose by default)
    
    -s : Save files, used to specify if user wants to preserve the temporary files used in the video production (videos are deleted by default)
    
    -t : Template, used to input a specific template file to use, otherwise the program searches current directory for any .slideshow files and uses the first it finds
    
    -o : Output location, used to specify where to store the finished video, will use current directory by default
    
    -l : Lower quality, used to generate a lower quality video for smaller file size for easier distribution (default videos will be 1280x720)


