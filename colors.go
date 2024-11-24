package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Color int

const (
	Green Color = iota
	Yellow
	Red
	Blue
	Purple
	Black
	White
)

var colorNames = []string{
	"Green ",
	"Yellow",
	" Red  ",
	" Blue ",
	"Purple",
	"Black ",
	"White ",
}

var colorPrinters = []func(format string, a ...interface{}) string{
	color.New(color.BgGreen, color.FgBlack).SprintfFunc(),
	color.New(color.BgYellow, color.FgBlack).SprintfFunc(),
	color.New(color.BgRed, color.FgBlack).SprintfFunc(),
	color.New(color.BgBlue, color.FgBlack).SprintfFunc(),
	color.New(color.BgMagenta, color.FgBlack).SprintfFunc(),
	color.New(color.BgBlack, color.FgWhite).SprintfFunc(),
	color.New(color.BgWhite, color.FgBlack).SprintfFunc(),
}

func ParseColor(s string) (Color, error) {
	switch strings.ToLower(s) {
	case "green", "g":
		return Green, nil
	case "yellow", "y":
		return Yellow, nil
	case "red", "r":
		return Red, nil
	case "blue", "b":
		return Blue, nil
	case "purple", "p":
		return Purple, nil
	case "black", "bk", "k":
		return Black, nil
	case "white", "w":
		return White, nil
	}
	return White, fmt.Errorf("unknown color: %s", s)
}

func (c Color) IsCrazy() bool {
	return c >= Black
}

func (c Color) String() string {
	return colorPrinters[c](colorNames[c])
}
