package main

import (
	"os"
	"strings"

	"github.com/cavanwang/translate"
)

func main() {
	words := strings.Join(os.Args[1:], " ")
	if v, _ := os.LookupEnv("Pronounce"); v == "English" {
		translate.Pronounce(words, false)
	} else {
		translate.Pronounce(words, true)
	}
}
