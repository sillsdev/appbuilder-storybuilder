module github.com/gordon-cs/SIL-Video/Compiler

go 1.17

replace github.com/gordon-cs/SIL-Video/Compiler/slideshow => ./src/slideshow

replace github.com/gordon-cs/SIL-Video/Compiler/helper => ./src/helper

replace github.com/gordon-cs/SIL-Video/Compiler/xml => ./src/xml

replace github.com/gordon-cs/SIL-Video/Compiler/ffmpeg_pkg => ./src/ffmpeg

require (
	github.com/gordon-cs/SIL-Video/Compiler/ffmpeg_pkg v0.0.0-00010101000000-000000000000
	github.com/gordon-cs/SIL-Video/Compiler/helper v0.0.0-00010101000000-000000000000
	github.com/gordon-cs/SIL-Video/Compiler/slideshow v0.0.0-00010101000000-000000000000
)

require github.com/gordon-cs/SIL-Video/Compiler/xml v0.0.0-00010101000000-000000000000 // indirect
