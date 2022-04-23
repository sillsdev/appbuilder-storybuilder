package os

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sillsdev/appbuilder-storybuilder/src/helper"
)

func DeleteTemporaryDirectory(saveTemps bool) {
	if !saveTemps {
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
