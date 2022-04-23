package os

import (
	"errors"
	"fmt"
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

func CreateDirectory(directory string) (string, error) {
	if directory == "" {
		dir, err := os.MkdirTemp("", "storybuilder-*")
		directory = dir
		return directory, err
	} else {
		if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(directory, os.ModePerm)
			return directory, err
		}
	}
	return directory, nil
}
