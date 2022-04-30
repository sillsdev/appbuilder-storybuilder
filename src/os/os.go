package os

import (
	"errors"
	"fmt"
	"os"
)

/* Function to remove the temp directory created during the process
 *
 * Parameters:
 *			tempDirectory - the name of the directory to remove, if null, remove default "./temp"
 * Returns:
 *			err - error code in the event of a failure, error is nil if successful
 */
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

/* Function to create a directory at specified location, or in the OS default temp directory by default
 *
 * Parameters:
 *			directory - the name of the directory to create
 *			v - verbose flag to determine what feedback to print
 * Returns:
 *			directory - the path to the created directory
 *			err - error code in the event of a failure, error is nil if successful
 */
func CreateDirectory(directory string, v bool) (string, error) {
	var err error
	var dir string
	if directory == "" {
		dir, err = os.MkdirTemp("", "storybuilder-*")
		directory = dir
	} else {
		if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(directory, os.ModePerm)
		}
	}
	if v {
		println("Created directory: ", dir)
	}
	return directory, err
}
