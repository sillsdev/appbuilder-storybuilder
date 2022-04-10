package helper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

/* Function to split the motion data into 4 pieces and convert them all to floats
 *  Parameters:
 *			stringData (string): The string that contains the four numerical values separated by spaces
 *  Returns:
 *			A float64 array with the four converted values
 */
func ConvertStringToFloat(stringData string) []float64 {
	floatData := []float64{}
	slicedStrings := strings.Split(stringData, " ")
	for _, str := range slicedStrings {
		if str != "" {
			flt, err := strconv.ParseFloat(str, 64)
			Check(err)
			floatData = append(floatData, flt)
		}
	}
	return floatData
}

// Function to check errors from non-CMD output
func Check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}

func RemoveFileNameFromDirectory(slideshowDirectory string) string {
	template_directory_split := strings.Split(slideshowDirectory, "/")
	template_directory := ""

	if len(template_directory_split) == 1 {
		template_directory = "./"
	} else {
		for i := 0; i < len(template_directory_split)-1; i++ {
			template_directory += template_directory_split[i] + "/"
		}
	}
	return template_directory
}
