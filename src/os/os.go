package os

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
)

func DeleteTemporaryVideos(saveTemps *bool) {
	if !*saveTemps {
		fmt.Println("-s not specified, removing temporary videos...")
		err := os.RemoveAll("./temp")
		helper.Check(err)
	}
}

func CreateDirectory(directory string) string {
	if directory == "" {
		dir, err := os.MkdirTemp("", "storybuilder-*")
		helper.Check(err)
		directory = dir
	} else {
		if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(directory, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return directory
}
