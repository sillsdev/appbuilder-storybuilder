package os

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gordon-cs/SIL-Video/Compiler/helper"
)

type opSys struct {
}

func NewOS() opSys {
	return opSys{}
}

func (o opSys) DeleteTemporaryVideos(saveTemps *bool) {
	if !*saveTemps {
		fmt.Println("-s not specified, removing temporary videos...")
		err := os.RemoveAll("./temp")
		helper.Check(err)
	}
}

func (o opSys) CreateTemporaryFolder(tempPath string) string {
	if tempPath == "" {
		dir, err := os.MkdirTemp("", "storybuilder-*")
		helper.Check(err)
		tempPath = dir
	} else {
		o.CreateDirectory(tempPath)
	}
	return tempPath
}
func (o opSys) CreateDirectory(location string) {
	if _, err := os.Stat(location); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(location, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
