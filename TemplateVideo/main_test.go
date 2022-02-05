package main

import (
	"io/ioutil"
	"testing"
)

//var templateName string
//var slideshow = readData("/Users/roddy/Desktop/SeniorProject/SIL-Video/")

func TestParse(t *testing.T) {
	input := " Parsing works corretcly"
	expectedOutput := " Parsing works correctly"

	data, err := ioutil.ReadFile(input) /// Don't know to to use the parsing input here
	// any help would ve appreciated

	if err != nil {
		t.Error("expected no error, but got %", err)
	}
	if string(data) != expectedOutput {
		t.Error("expected output to be string, but got variable", expectedOutput, data)
	}
}

// func TestCheck(t *testing.T) {

// }

// func TestReadFile(t *testing.T) {
// 	data, err := ioutil.ReadFile("data.slideshow")
// 	if err != nil {

// 	}
// 	if string(readData) != nil {

// 	}
//// }//
