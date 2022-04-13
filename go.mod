module github.com/gordon-cs/SIL-Video/Compiler

go 1.17

replace github.com/gordon-cs/SIL-Video/Compiler/slideshow => ./src/slideshow

replace github.com/gordon-cs/SIL-Video/Compiler/helper => ./src/helper

require (
	github.com/gordon-cs/SIL-Video/Compiler/helper v0.0.0-00010101000000-000000000000
	github.com/gordon-cs/SIL-Video/Compiler/slideshow v0.0.0-00010101000000-000000000000 // indirect
)
