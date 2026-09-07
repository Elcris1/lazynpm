package ui

import (
	"os"
	"golang.org/x/term"
)


func GetLeftColumnWidth() int {
	width, _ := getTermSize()

	return int(width/10) * 4
}

func GetRightColumnWidth() int {
	width, _ := getTermSize()

	return int(width/10) * 6
}

func getTermSize() (int, int) {

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	
	if err != nil {
		width = 100
		height = 100
	}

	return width, height
}