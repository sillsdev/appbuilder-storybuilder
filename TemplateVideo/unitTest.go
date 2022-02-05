package main

import (
	"io/ioutil"
	"testing"
)

var slideshow = readData(templateName)

func TestParse(t *testing.T) {

	input := " Parsing works corretcly"
	expectedOutput := " Parsing works correctly"

	data, err := ioutil.ReadFile(input) /// Don't know to to use the parsing input here
	// any help would ve appreciated

	if err != nil {
		t.Error("expected no error, but got %v", err)
	}
	if sring(data != expectedOutput {
		t.Error("expected output to be %s, but got %v", expectedOutput, data)
	}
}

// func TestReadFile(t *testing.T) {
// 	data, err := ioutil.ReadFile("data.slideshow")
// 	if err != nil {

// 	}
// 	if string(readData) != nil {

// 	}
// }
