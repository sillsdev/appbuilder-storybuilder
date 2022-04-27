package os

import (
	"errors"
	"fmt"
	"os"
)

func DeleteTemporaryDirectory(tempDirectory string) error {
	fmt.Println("-s not specified, removing temporary videos...")

	var err error
	if tempDirectory == "" {
		err = os.RemoveAll("./temp")
	} else {
		err = os.RemoveAll(tempDirectory)
	}

	return err
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
