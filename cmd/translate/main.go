package main

import (
	"os"
	"strings"

	"github.com/cavanwang/translate"
)

func main() {
	translate.Translate(strings.Join(os.Args[1:], " "))
}
